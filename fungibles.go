package hiveenginego

import (
	"encoding/json"
)

/*
type tokenArray struct {
	tokArray []HiveEngineFungibleToken
}
*/

type HiveEngineFungibleToken struct {
	FungibleId           int             `json:"_id"`
	Issuer               string          `json:"issuer"`
	Symbol               string          `json:"symbol"`
	Name                 string          `json:"name"`
	MetaData             json.RawMessage `json:"metadata"`
	DelegationEnabled    bool            `json:"delegationEnabled"`
	Precision            int             `json:"precision"`
	StakingEnabled       bool            `json:"stakingEnabled"`
	UndelegationCooldown int             `json:"undelegationCooldown"`
	UnstakingCooldown    int             `json:"unstakingCooldown"`
}

type FungibleBalance struct {
	Id                   int    `json:"_id"`
	Account              string `json:"account"`
	Symbol               string `json:"symbol"`
	Balance              string `json:"balance"`
	Stake                string `json:"stake"`
	PendingUnstake       string `json:"pendingUnstake"`
	DelegationsIn        string `json:"delegationsIn"`
	DelegationsOut       string `json:"delegationsOut"`
	PendingUndelegations string `json:"pendingUndelegations"`
}

type FungibleTokenTransfer struct {
	Symbol   string `json:"symbol"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

func CreateFungibleTokenTransfer(symbol string, to string, quantity string, memo string) ContractTx {
	tokenTransfer := FungibleTokenTransfer{
		Symbol:   symbol,
		To:       to,
		Quantity: quantity,
		Memo:     memo,
	}
	contractTrx := ContractTx{
		ContractName:    "tokens",
		ContractAction:  "transfer",
		ContractPayload: tokenTransfer,
	}
	return contractTrx
}

func (h HiveEngineRpcNode) getFungibleTokenCount() (int, error) {
	if len(h.Endpoints.Contracts) == 0 {
		h.Endpoints.Contracts = "/contracts"
	}
	endpoint := h.Endpoints.Contracts
	qParamsIndex := []ContractQueryParamsIndex{}
	qParamsIndex = append(qParamsIndex, ContractQueryParamsIndex{Index: "_id", Descending: true})
	qParams := ContractQueryParams{Contract: "tokens", Table: "tokens", Query: "", Limit: 1, Offset: 0, Index: qParamsIndex}
	query := herpcQuery{method: "find", params: qParams}

	res, err := h.rpcExec(endpoint, query)
	if err != nil {
		return 0, err
	}

	token := []HiveEngineFungibleToken{}
	if err := json.Unmarshal(res, &token); err != nil {
		return 0, err
	}

	countSymbolNFT := token[0].FungibleId

	return countSymbolNFT, nil
}

func (h HiveEngineRpcNode) GetAllFungibleTokens() ([]HiveEngineFungibleToken, error) {
	totalTokens, err := h.getFungibleTokenCount()
	if err != nil {
		return nil, err
	}

	offsetsNeeded := totalTokens / 1000

	qParamsIndex := []ContractQueryParamsIndex{}

	var queries []herpcQuery

	for i := 0; i <= offsetsNeeded; i++ {
		offset := i * 1000
		qParams := ContractQueryParams{Contract: "tokens", Table: "tokens", Limit: 1000, Offset: offset, Index: qParamsIndex}
		query := herpcQuery{method: "find", params: qParams}
		queries = append(queries, query)
	}

	tokens := []HiveEngineFungibleToken{}

	if len(queries) > 0 {
		if len(h.Endpoints.Contracts) == 0 {
			h.Endpoints.Contracts = "/contracts"
		}
		endpoint := h.Endpoints.Contracts
		ress, err := h.rpcExecBatch(endpoint, queries)
		if err != nil {
			return nil, err
		}
		var batchResult []HiveEngineFungibleToken
		for _, res := range ress {
			thisresult := []HiveEngineFungibleToken{}
			if err := json.Unmarshal(res, &thisresult); err != nil { // Parse []byte to the go struct pointer
				return nil, err
			}
			batchResult = append(tokens, thisresult...)
		}
		tokens = append(tokens, batchResult...)
	}
	return tokens, nil
}
