# HiveEngineGo - A client for the Hive Engine side chain on the Hive blockchain

At this time, there are only a few functions from the client. More will be added.

### Example usage:
create a client:
```
herpc := hiveenginego.NewHiveEngineRpc("http://MyHiveEngineApi")
```

Query latest block info:
```
latestBlockInfo, err := herpc.GetLatestBlockInfo()
//Returns a struct
latestBlockNum := latestBlockInfo.BlockNumber
```

Get All NFT of a given symbol (return rpc resonse as raw bytes):
```
rawNftBytes, err := herpc.GetSymbolAllNftFast("STAR")
```

Get block range as the raw response from the rpc (in bytes):
```
rpcResponsesBytes, err := herpc.GetBlockRangeFast(start, end)
```
WARNING: It is not recommended to stream blocks from public APIs. They are provided as a service to users and saturating them with block requests may (rightfully) result in your IP getting banned
