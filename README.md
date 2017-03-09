# mongopool_txn
library for run transactions on mongodb

# features
1. support the count of max connections config
2. support redial mongo if mongo shutdown and setup again

# how to use

1. init a txnRunner 
```
	txnRunner = mongo_pool_txn.NewTxnRunner("127.0.0.1", "testdb", 2)
	txnRunner.Run()
```

2. do the transaction 
```
	now := time.Now().Format(time.RFC3339)
	ops := []txn.Op{{
		C:      "coltest",
		Id:     now,
		Assert: txn.DocMissing,
		Insert: obj{ID: now, Status: 1},
	}}

	txnRunner.Do(ops)
```