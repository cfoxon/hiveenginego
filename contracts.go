package hiveenginego

type ContractQueryParams struct {
    Contract string                     `json:"contract"`
    Table    string                     `json:"table"`
    Query    interface{}                `json:"query"`
    Limit    int                        `json:"limit"`
    Offset   int                        `json:"offset"`
    Index    []ContractQueryParamsIndex `json:"indexes"`
}

type QueryIDRange struct {
    Id QueryIntRange `json:"_id"`
}

type QueryIntRange struct {
    GreaterThanEqual int `json:"$gte,omitempty"`
    LessThanEqual int `json:"$lte,omitempty"`
}

type BroadcastTx struct {

}

type ContractQueryParamsQuery struct {
    Account string `json:"account,omitempty"`
    NftId   string `json:"_id,omitempty"`
}

type ContractQueryParamsIndex struct {
    Index      string `json:"index,omitempty"`
    Descending bool   `json:"descending"`
}

type ContractTx struct {
    ContractName string `json:"contractName"`
    ContractAction string `json:"contractAction"`
    ContractPayload interface{} `json:"contractPayload"`
}

func (h HiveEngineRpcNode) QueryContract(qParams ContractQueryParams) ([]byte, error){
    if h.Endpoints.Contracts == "" {
        h.Endpoints.Contracts = "/contract"
    }
    query := herpcQuery{method: "find", params: qParams}
    qRes, err := h.rpcExec(h.Endpoints.Contracts, query)
    if err != nil {
        return nil, err
    }
    return qRes, nil
}

func (h HiveEngineRpcNode) QueryContractBatch(qParams []ContractQueryParams) ([][]byte, error){
    if h.Endpoints.Contracts == "" {
        h.Endpoints.Contracts = "/contract"
    }

    var queries []herpcQuery
    for _, qParam := range qParams{
        query := herpcQuery{method: "find", params: qParam}
        queries = append(queries, query)
    }
    qRes, err := h.rpcExecBatch(h.Endpoints.Contracts, queries)
    if err != nil{
        return nil, err
    }

    return qRes, nil
}