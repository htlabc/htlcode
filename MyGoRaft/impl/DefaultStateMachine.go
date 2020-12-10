package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"htl/myraft.com/entry"
	inter "htl/myraft.com/inter"
	"log"
	"sync"
)

var pre_db_table_text string = "raft_blotdb_"

type DefaultStateMachine struct {
	filename        string
	dbDir           string
	stateMachineDir string
	machineDb       *bolt.DB
	inter.StateMachine
}

var DbDir = ""
var DbFileName = ""

func NewDefaultStateMachine() *DefaultStateMachine {
	var lock *sync.Mutex = &sync.Mutex{}
	d := new(DefaultStateMachine)
	d.dbDir = DbDir
	d.stateMachineDir = DbDir + "/stateMachine/" + DbFileName
	lock.Lock()
	defer lock.Unlock()
	db, err := bolt.Open(d.stateMachineDir, 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	d.machineDb = db
	return d
}

func (d *DefaultStateMachine) Apply(log *entry.LogEntry) {
	var command *entry.Command = log.GetCommand()
	if command != nil {
		panic(errors.New("command can not be null, logEntry :"))
	}

	var key string = command.GetKey()
	tablename := pre_db_table_text + key
	d.machineDb.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucket([]byte(tablename))
		data, _ := json.Marshal(log)
		b.Put([]byte(key), data)
		return nil
	})
}

func (d *DefaultStateMachine) Get(key string) *entry.LogEntry {
	var data string
	var entry entry.LogEntry
	tablename := pre_db_table_text + key

	d.machineDb.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tablename))
		if b != nil {
			data = string(b.Get([]byte(key)))
		}
		json.Unmarshal([]byte(data), entry)
		return nil
	})
	return &entry
}

func (d *DefaultStateMachine) GetString(key string) string {
	var data string
	tablename := pre_db_table_text + key
	d.machineDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tablename))
		if b != nil {
			data = string(b.Get([]byte(key)))
		}
		return nil
	})
	return data
}

func (d *DefaultStateMachine) SetString(key string, value string) {
	tablename := pre_db_table_text + key
	err := d.machineDb.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucket([]byte(tablename))
		if b != nil {
			err := b.Put([]byte(key), []byte(value))
			if err != nil {
				log.Panic("数据存储失败")
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (d *DefaultStateMachine) DelString(key ...string) {
	d.machineDb.Update(func(tx *bolt.Tx) error {
		for _, value := range key {
			tx.DeleteBucket([]byte(value))
		}
		return nil
	})
}
