package main

import (
	"sync"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/data"
)

var once sync.Once
var defaultStore *MuStore

type MuStore struct {
	store map[data.Data]data.Data
}

func GetStore() *MuStore {
	once.Do(func() {
		if defaultStore == nil {
			defaultStore = &MuStore{store: make(map[data.Data]data.Data)}
		}
	})

	return defaultStore
}

func (store *MuStore) Set(key data.Data, value data.Data) {
	store.store[key] = value
}

func (store *MuStore) Get(key data.Data) data.Data {
	if value, ok := store.store[key]; !ok {
		return data.NewNullMuString()
	} else {
		return value
	}
}
