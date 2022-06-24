package hiveenginego

import (
	"encoding/json"
	"strings"
)

type Witness struct {
	Id                 int             `json:"_id"`
	Account            string          `json:"account"`
	ApprovalWeight     json.RawMessage `json:"approvalWeight"`
	SigningKey         string          `json:"signingKey"`
	Ip                 string          `json:"ip"`
	IpVersion          int             `json:"ipVersion"`
	RPCPort            int             `json:"RPCPort"`
	P2PPort            int             `json:"P2PPort"`
	Enabled            bool            `json:"enabled"`
	MissedRounds       int             `json:"missedRounds"`
	MissedRoundsInARow int             `json:"missedRoundsInARow"`
	VerifiedRounds     int             `json:"verifiedRounds"`
	LastRoundVerified  int             `json:"lastRoundVerified"`
	LastBlockVerified  int             `json:"lastBlockVerified"`
}

func (h HiveEngineRpcNode) getAllWitnesses() ([]Witness, error) {
	if len(h.Endpoints.Contracts) == 0 {
		h.Endpoints.Contracts = "/contracts"
	}
	endpoint := h.Endpoints.Contracts

	qParamsIndex := []ContractQueryParamsIndex{{Index: "_id", Descending: false}}
	qParams := ContractQueryParams{Contract: "witnesses", Table: "witnesses", Query: struct{}{}, Limit: 1000, Offset: 0, Index: qParamsIndex}
	query := herpcQuery{method: "find", params: qParams}

	res, err := h.rpcExec(endpoint, query)
	if err != nil {
		return nil, err
	}

	witnesses := []Witness{}

	if err := json.Unmarshal(res, &witnesses); err != nil {
		return nil, err
	}

	for i, witness := range witnesses {
		if strings.Contains(witness.Ip, ":") {
			witness.IpVersion = 6
		} else {
			witness.IpVersion = 4
		}
		witnesses[i] = witness
	}

	return witnesses, nil
}
