package mongo_pool_txn

import (
	"github.com/ssor/mongopool"
	"gopkg.in/mgo.v2/txn"
)

var (
	collectionForTxn = "tc"
)

type OpArray []txn.Op

func NewTxnRunnerWithMongoPool(db string, mongoPool *mongo_pool.MongoSessionPool) *TxnRuner {
	tr := &TxnRuner{}
	tr.dbName = db
	tr.mongoPool = mongoPool
	return tr

}

func NewTxnRunner(hosts, db string, max_session_count int) *TxnRuner {
	tr := &TxnRuner{}
	tr.dbName = db
	tr.mongoHosts = hosts
	tr.maxSessionCount = max_session_count
	return tr
}

type TxnRuner struct {
	mongoHosts      string
	maxSessionCount int
	dbName          string
	mongoPool       *mongo_pool.MongoSessionPool
}

func (tr *TxnRuner) Run() {
	tr.mongoPool = mongo_pool.NewMongoSessionPool(tr.mongoHosts, tr.maxSessionCount)
	tr.mongoPool.Run()
}

func (tr *TxnRuner) Do(ops OpArray) error {
	if len(ops) == 0 {
		return nil
	}
	var err error
	session, err := tr.mongoPool.GetSession()
	defer func() {
		tr.mongoPool.ReturnSession(session, err)
	}()
	if err != nil {
		return err
	}
	txnRunner := txn.NewRunner(session.DB(tr.dbName).C(collectionForTxn))
	err = txnRunner.Run(ops, "", nil)
	if err != nil {
		//EOF:数据库连接中断
		return err
	}

	return nil
}
