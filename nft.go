package hiveenginego

import (
	"encoding/json"
	"strconv"
	"strings"
)

type EngineNft struct {
	Id              int             `json:"_id"`
	Account         string          `json:"account"`
	OwnedBy         string          `json:"ownedBy"`
	LockedTokens    json.RawMessage `json:"lockedTokens"`
	Properties      json.RawMessage `json:"properties"`
	PreviousAccount string          `json:"previousAccount"`
	PreviousOwnedBy string          `json:"previousOwnedBy"`
}

type NftTransferPayload struct {
	Nfts []NftsForNftTransfer `json:"nfts"`
	To   string               `json:"to"`
	Memo string               `json:"memo,omitempty"`
}

type NftsForNftTransfer struct {
	Symbol string   `json:"symbol"`
	Ids    []string `json:"ids"`
}

func CreateNftTransfer(symbol string, nftIds []int, to string, memo string) []ContractTx {
	var nftIdsStr []string
	for _, nftId := range nftIds {
		nftIdsStr = append(nftIdsStr, strconv.Itoa(nftId))
	}
	var preppedTx []ContractTx
	if len(nftIdsStr) > 50 {
		for i := 0; i <= len(nftIdsStr); {
			var l int
			if i+50 > len(nftIdsStr) {
				l = len(nftIdsStr)
			} else {
				l = i + 50
			}
			strid := nftIdsStr[i:l]
			nftTransfer := NftsForNftTransfer{Symbol: symbol, Ids: strid}
			payload := NftTransferPayload{Nfts: []NftsForNftTransfer{nftTransfer}, To: to, Memo: memo}
			contrx := ContractTx{ContractName: "nft", ContractAction: "transfer", ContractPayload: payload}
			preppedTx = append(preppedTx, contrx)
			i += 50
		}
	} else {
		nftTransfer := NftsForNftTransfer{Symbol: symbol, Ids: nftIdsStr}
		payload := NftTransferPayload{Nfts: []NftsForNftTransfer{nftTransfer}, To: to, Memo: memo}
		contrx := ContractTx{ContractName: "nft", ContractAction: "transfer", ContractPayload: payload}
		preppedTx = append(preppedTx, contrx)
	}
	return preppedTx
}

func (h HiveEngineRpcNode) getSymbolNFTCount(nftSymbol string) (int, error) {
	nftSymbolUpper := strings.ToUpper(nftSymbol)
	if len(h.Endpoints.Contracts) == 0 {
		h.Endpoints.Contracts = "/contracts"
	}
	endpoint := h.Endpoints.Contracts

	qParamsIndex := []ContractQueryParamsIndex{{Index: "_id", Descending: true}}
	qParams := ContractQueryParams{Contract: "nft", Table: nftSymbolUpper + "instances", Query: struct{}{}, Limit: 1, Offset: 0, Index: qParamsIndex}
	query := herpcQuery{method: "find", params: qParams}

	res, err := h.rpcExec(endpoint, query)
	if err != nil {
		return 0, err
	}

	nft := []EngineNft{}

	if err := json.Unmarshal(res, &nft); err != nil {
		return 0, nil
	}

	countSymbolNFT := nft[0].Id

	return countSymbolNFT, nil
}

func (h HiveEngineRpcNode) GetSymbolAllNft(nftSymbol string) ([]EngineNft, error) {
	if len(h.Endpoints.Contracts) == 0 {
		h.Endpoints.Contracts = "/contracts"
	}
	if h.RpcNode.MaxConn == 0 {
		h.RpcNode.MaxConn = 1
	}
	if h.RpcNode.MaxBatch == 0 {
		h.RpcNode.MaxBatch = 2
	}

	nftSymbolUpper := strings.ToUpper(nftSymbol)
	collectionSize, err := h.getSymbolNFTCount(nftSymbolUpper)
	if err != nil {
		return nil, err
	}
	offsetsNeeded := collectionSize / 1000

	qParamsIndex := []ContractQueryParamsIndex{}

	var queries []herpcQuery
	for i := 0; i <= offsetsNeeded; i++ {
		offset := i * 1000
		queryFilter := QueryIDRange{QueryIntRange{offset, offset + 999}}
		qParams := ContractQueryParams{Contract: "nft", Table: nftSymbolUpper + "instances", Query: queryFilter, Limit: 1000, Offset: 0, Index: qParamsIndex}
		query := herpcQuery{method: "find", params: qParams}
		queries = append(queries, query)
	}

	endpoint := h.Endpoints.Contracts

	ress, err := h.rpcExecBatch(endpoint, queries)
	if err != nil {
		return nil, err
	}

	var nfts []EngineNft
	for _, res := range ress {
		thisresult := []EngineNft{}
		if err := json.Unmarshal(res, &thisresult); err != nil { // Parse []byte to the go struct pointer
			return nil, err
		}
		nfts = append(nfts, thisresult...)
	}

	return nfts, nil
}

func (h HiveEngineRpcNode) GetSymbolAllNftFast(nftSymbol string) ([][]byte, error) {
	if len(h.Endpoints.Contracts) == 0 {
		h.Endpoints.Contracts = "/contracts"
	}
	if h.RpcNode.MaxConn == 0 {
		h.RpcNode.MaxConn = 1
	}
	if h.RpcNode.MaxBatch == 0 {
		h.RpcNode.MaxBatch = 2
	}

	nftSymbolUpper := strings.ToUpper(nftSymbol)
	collectionSize, err := h.getSymbolNFTCount(nftSymbolUpper)
	if err != nil {
		return nil, err
	}
	offsetsNeeded := collectionSize / 1000

	qParamsIndex := []ContractQueryParamsIndex{}

	var queries []herpcQuery
	for i := 0; i <= offsetsNeeded; i++ {
		offset := i * 1000
		queryFilter := QueryIDRange{QueryIntRange{offset, offset + 999}}
		qParams := ContractQueryParams{Contract: "nft", Table: nftSymbolUpper + "instances", Query: queryFilter, Limit: 1000, Offset: 0, Index: qParamsIndex}
		query := herpcQuery{method: "find", params: qParams}
		queries = append(queries, query)
	}

	endpoint := h.Endpoints.Contracts

	ress, err := h.rpcExecBatchFast(endpoint, queries)
	if err != nil {
		return nil, err
	}

	return ress, nil
}
