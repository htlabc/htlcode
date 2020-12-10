package impl

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"htl/myraft.com/entry"
	"strconv"
	"sync"
)

var LAST_INDEX_KEY []byte = []byte("LAST_INDEX_KEY")
var DEFAULTLOGMODULE_TABLE = "TB_DEFAULT_LOGMODULE"

type DefaultLogModule struct {
	dbDir        string
	logModuleDir string
	logModulerDb *bolt.DB
	lock         sync.Mutex
}

func NewDefaultLogModule() *DefaultLogModule {
	var lock *sync.Mutex = &sync.Mutex{}
	d := new(DefaultLogModule)
	d.dbDir = DbDir
	d.logModuleDir = DbDir + "/logmodule/" + DbFileName
	lock.Lock()
	defer lock.Unlock()
	db, err := bolt.Open(d.logModuleDir, 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	d.logModulerDb = db
	return d
}

func (l *DefaultLogModule) Write(logEntry *entry.LogEntry) {
	success := false
	l.lock.Lock()
	defer l.lock.Unlock()
	logEntry.SetIndex(l.GetLastIndex() + 1)
	err := l.logModulerDb.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucket([]byte(DEFAULTLOGMODULE_TABLE))
		data, _ := json.Marshal(logEntry)
		key := logEntry.GetCommand().GetKey()
		b.Put([]byte(key), data)
		return nil
	})
	if err != nil {
		fmt.Println("exec logmoduler.write failed: because of ", err)
		return
	}
	success = true
	fmt.Printf("DefaultLogModule write rocksDB success, logEntry info : %s", logEntry)
	if success {
		l.updateLastIndex(logEntry.GetIndex())
	}
}

func (l *DefaultLogModule) updateLastIndex(index int64) {
	value := strconv.FormatInt(index, 10)
	l.logModulerDb.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucket([]byte(DEFAULTLOGMODULE_TABLE))
		b.Put([]byte(LAST_INDEX_KEY), []byte(value))
		return nil
	})
}

func (l *DefaultLogModule) Read(index int64) *entry.LogEntry {
	key := strconv.FormatInt(index, 10)
	var entry entry.LogEntry
	data := make([]byte, 0)
	l.logModulerDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULTLOGMODULE_TABLE))
		if b != nil {
			data = b.Get([]byte(key))
		}
		return nil
	})
	json.Unmarshal(data, entry)
	return &entry
}

func (l *DefaultLogModule) RemoveOnstartIndex(startIndex int64) {
	lastIndex := l.GetLast().GetIndex()
	var count int64 = 0
	l.logModulerDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULTLOGMODULE_TABLE))
		for i := startIndex; i < lastIndex; i++ {
			if b != nil {
				key := strconv.FormatInt(i, 10)
				b.Delete([]byte(key))
				count++
			}
		}
		return nil
	})
	l.updateLastIndex(lastIndex - count)
}

func (l *DefaultLogModule) GetLast() *entry.LogEntry {
	var entry entry.LogEntry
	data := make([]byte, 0)
	l.logModulerDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULTLOGMODULE_TABLE))
		if b != nil {
			data = b.Get([]byte(strconv.FormatInt(l.GetLastIndex(), 10)))
		}
		return nil
	})
	json.Unmarshal(data, entry)
	return &entry
}

func (l *DefaultLogModule) GetLastIndex() int64 {
	var lastIndex int64 = 0
	data := ""
	l.logModulerDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULTLOGMODULE_TABLE))
		if b != nil {
			data = string(b.Get([]byte(LAST_INDEX_KEY)))
		}
		return nil
	})
	lastIndex, _ = strconv.ParseInt(data, 10, 64)
	return lastIndex
}
