package main

import (
	"sync"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/data"
)

var once sync.Once
var defaultStore *muStore

type muStore struct {
	store map[data.Data]data.Data
}

func getStore() *muStore {
	once.Do(func() {
		if defaultStore == nil {
			defaultStore = &muStore{store: make(map[data.Data]data.Data)}
		}
	})

	return defaultStore
}

func (store *muStore) Set(key data.Data, value data.Data) {
	store.store[key] = value
}

func (store *muStore) Get(key data.Data) data.Data {
	if value, ok := store.store[key]; !ok {
		return data.NewNullMuString()
	} else {
		return value
	}
}
