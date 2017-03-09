package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2/txn"

	"github.com/ssor/mongopool_txn"
)

var (
	txnRunner *mongo_pool_txn.TxnRuner
)

func main() {
	InitMongo()
	loop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("[OK] Quit")
}

func loop() {

	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			conn_to_db := func() {
				err := SaveUserLoginInfoToDB(txnRunner)
				if err != nil {
					fmt.Println("*** save user err: ", err)
				} else {
					fmt.Println("[OK] save user to do success")
				}
			}
			conn_to_db()
		}
	}()
}

func InitMongo() {
	txnRunner = mongo_pool_txn.NewTxnRunner("127.0.0.1", "testdb", 2)
	txnRunner.Run()
}

type obj struct {
	ID     string `bson:"_id"`
	Status int    `bson:"status"`
}

func SaveUserLoginInfoToDB(txnRunner *mongo_pool_txn.TxnRuner) error {
	now := time.Now().Format(time.RFC3339)
	ops := []txn.Op{{
		C:      "coltest",
		Id:     now,
		Assert: txn.DocMissing,
		Insert: obj{ID: now, Status: 1},
	}}

	return txnRunner.Do(ops)
}
