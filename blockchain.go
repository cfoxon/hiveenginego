package hiveenginego

import (
	"encoding/json"
)

type EngineStatus struct {
	LastBlockNumber             int    `json:"lastBlockNumber"`
	LastBlockRefHiveBlockNumber int    `json:"lastBlockRefHiveBlockNumber"`
	LastParsedHiveBlockNumber   int    `json:"lastParsedHiveBlockNumber"`
	LastVerifiedBlockNumber     int    `json:"lastVerifiedBlockNumber"`
	SSCnodeVersion              string `json:"SSCnodeVersion"`
	ChainId                     string `json:"chainId"`
	Domain			    string `json:"domain"`
	Lightnode		    bool   `json:"lightnode"`
	LastHash		    string `json:"lastHash"`
}

type EngineBlock struct {
	BlockNumber         int             `json:"blockNumber"`
	RefHiveBlockNumber  int             `json:"refHiveBlockNumber"`
	Timestamp           string          `json:"timestamp"`
	Transactions        json.RawMessage `json:"transactions"`
	VirtualTransactions json.RawMessage `json:"virtualTransactions,omitempty"`
	Hash                string          `json:"hash"`
	DatabaseHash        string          `json:"databasehash"`
	MerkleRoot          string          `json:"merkleRoot"`
	Round               int             `json:"round"`
	RoundHash           string          `json:"roundHash"`
	Witness             string          `json:"witness"`
	SigningKey          string          `json:"signingKey"`
	RoundSignature      string          `json:"roundSignature"`
}

type blockChainQueryParams struct {
	BlockNumber int `json:"blockNumber"`
}

func (h HiveEngineRpcNode) GetStatus() (*EngineStatus, error) {
	if len(h.Endpoints.Blockchain) == 0 {
		h.Endpoints.Blockchain = "/blockchain"
	}

	endpoint := h.Endpoints.Blockchain
	query := herpcQuery{method: "getStatus"}

	res, err := h.rpcExec(endpoint, query)
	if err != nil {
		return nil, err
	}

	status := &EngineStatus{}

	if err := json.Unmarshal(res, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (h HiveEngineRpcNode) GetLatestBlockInfo() (*EngineBlock, error) {
	if len(h.Endpoints.Blockchain) == 0 {
		h.Endpoints.Blockchain = "/blockchain"
	}

	endpoint := h.Endpoints.Blockchain
	query := herpcQuery{method: "getLatestBlockInfo"}

	res, err := h.rpcExec(endpoint, query)
	if err != nil {
		return nil, err
	}

	block := &EngineBlock{}
	if err := json.Unmarshal(res, &block); err != nil {
		return nil, err
	}

	return block, nil
}

func (h HiveEngineRpcNode) GetBlockRange(startBlock int, endBlock int) ([][]byte, error) {
	if len(h.Endpoints.Blockchain) == 0 {
		h.Endpoints.Blockchain = "/blockchain"
	}

	if h.RpcNode.MaxConn == 0 {
		h.RpcNode.MaxConn = 1
	}

	if h.RpcNode.MaxBatch == 0 {
		h.RpcNode.MaxBatch = 2
	}

	var queries []herpcQuery

	for i := startBlock; i <= endBlock; i++ {
		params := blockChainQueryParams{BlockNumber: i}
		query := herpcQuery{method: "getBlockInfo", params: params}
		queries = append(queries, query)
	}

	endpoint := h.Endpoints.Blockchain

	res, err := h.rpcExecBatch(endpoint, queries)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h HiveEngineRpcNode) GetBlockRangeFast(startBlock int, endBlock int) ([][]byte, error) {
	if len(h.Endpoints.Blockchain) == 0 {
		h.Endpoints.Blockchain = "/blockchain"
	}

	if h.RpcNode.MaxConn == 0 {
		h.RpcNode.MaxConn = 1
	}

	if h.RpcNode.MaxBatch == 0 {
		h.RpcNode.MaxBatch = 2
	}

	var queries []herpcQuery

	for i := startBlock; i <= endBlock; i++ {
		params := blockChainQueryParams{BlockNumber: i}
		query := herpcQuery{method: "getBlockInfo", params: params}
		queries = append(queries, query)
	}

	endpoint := h.Endpoints.Blockchain

	res, err := h.rpcExecBatchFast(endpoint, queries)
	if err != nil {
		return nil, err
	}

	return res, nil
}
