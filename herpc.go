package hiveenginego

import (
	"errors"
	"github.com/cfoxon/jrc"
	"strconv"
)

type HiveEngineRpcNode struct {
	Endpoints engineApiEndpoints
	RpcNode   rpcServer
}

type rpcServer struct {
	address  string
	MaxConn  int
	MaxBatch int
	UseFast  bool
}

type herpcQuery struct {
	method string
	params interface{}
}

type engineApiEndpoints struct {
	Blockchain string
	Contracts  string
}

func NewHiveEngineRpc(addr string) *HiveEngineRpcNode {
	return NewHiveEngineRpcWithOpts(addr, "/blockchain", "/contracts", 1, 4)
}

func NewHiveEngineRpcWithOpts(addr string, blockchain string, contracts string, maxConn int, maxBatch int) *HiveEngineRpcNode {
	return &HiveEngineRpcNode{
		Endpoints: engineApiEndpoints{
			Blockchain: blockchain,
			Contracts:  contracts,
		},
		RpcNode: rpcServer{
			address:  addr,
			MaxConn:  maxConn,
			MaxBatch: maxBatch},
	}
}

func (h *HiveEngineRpcNode) rpcExec(endpoint string, query herpcQuery) ([]byte, error) {
	rpcClient, err := jrc.NewServer(h.RpcNode.address+endpoint, jrc.MaxCon(h.RpcNode.MaxConn), jrc.MaxBatch(h.RpcNode.MaxBatch))
	if err != nil {
		return nil, err
	}
	jr2query := jrc.RpcRequest{Method: query.method, JsonRpc: "2.0", Id: 1, Params: query.params}
	resp, err := rpcClient.Exec(jr2query)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, errors.New(strconv.Itoa(resp.Error.Code) + "    " + resp.Error.Message)
	}

	return resp.Result, nil
}

func (h *HiveEngineRpcNode) rpcExecBatch(endpoint string, queries []herpcQuery) ([][]byte, error) {
	rpcClient, err := jrc.NewServer(h.RpcNode.address+endpoint, jrc.MaxCon(h.RpcNode.MaxConn), jrc.MaxBatch(h.RpcNode.MaxBatch))
	var jr2queries jrc.RPCRequests
	for i, query := range queries {
		jr2query := &jrc.RpcRequest{Method: query.method, JsonRpc: "2.0", Id: i, Params: query.params}
		jr2queries = append(jr2queries, jr2query)
	}

	resps, err := rpcClient.ExecBatch(jr2queries)
	if err != nil {
		return nil, err
	}

	var batchResult [][]byte
	for _, resp := range resps {
		if resp.Error != nil {
			return nil, errors.New(strconv.Itoa(resp.Error.Code) + "    " + resp.Error.Message)
		}
		batchResult = append(batchResult, resp.Result)
	}

	return batchResult, nil
}

func (h *HiveEngineRpcNode) rpcExecBatchFast(endpoint string, queries []herpcQuery) ([][]byte, error) {
	rpcClient, err := jrc.NewServer(h.RpcNode.address+endpoint, jrc.MaxCon(h.RpcNode.MaxConn), jrc.MaxBatch(h.RpcNode.MaxBatch))
	var jr2queries jrc.RPCRequests
	for i, query := range queries {
		jr2query := &jrc.RpcRequest{Method: query.method, JsonRpc: "2.0", Id: i + 1, Params: query.params}
		jr2queries = append(jr2queries, jr2query)
	}

	resps, err := rpcClient.ExecBatchFast(jr2queries)
	if err != nil {
		return nil, err
	}

	var batchResult [][]byte
	for _, resp := range resps {
		batchResult = append(batchResult, resp)
	}

	return batchResult, nil
}
