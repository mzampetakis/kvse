// Package kvse is an In memory Key-Value Store with Expiration for each pair
// It provides the necessary structs and methods to maintain a key-value
// map with expiration for each pair.
// Expiration is managed by the lib itself but can be omitted if set to 0.
// Also it supports a minimum precision you want to achieve for deleting expired pairs.
// Providing a 0 precision will cause to use the default one which is time.second.
package kvse

import (
	"sync"
	"time"
)

// DataStore stores the struct of each list
// Holds a pointer to the list, a mutex lock and the list itself
type DataStore struct {
	data            map[string]mapValue
	mx              *sync.RWMutex
	deletePrecision time.Duration
}

type mapValue struct {
	expiration int64
	value      interface{}
}

// New func creates and returns a new kvse Datastore
// and initiates a worker to check and delete for expired keys.
// {precision} is the minimum precision you want to achieve for deleting expired pairs.
// Providing the precision as 0 it will use a default one: time.second.
func New(precision time.Duration) DataStore {
	if precision == 0 {
		precision = time.Second
	}
	ds := DataStore{
		data:            map[string]mapValue{},
		mx:              &sync.RWMutex{},
		deletePrecision: precision,
	}
	go ds.checkAndDeleteExpiredKeys()
	return ds
}

// Has returns a boolean based on whether or not the store contains a value for
// {key}.
func (ds DataStore) Has(key string) bool {
	ds.mx.RLock()
	defer ds.mx.RUnlock()
	_, ok := ds.data[key]
	return ok
}

// Get retrieves the value associated to the {key} in the store.
func (ds DataStore) Get(key string) (interface{}, bool) {
	ds.mx.RLock()
	defer ds.mx.RUnlock()
	data, ok := ds.data[key]
	return data.value, ok
}

// Set adds a ne value to a specific key with a {lifespan} duration.
// Setting the {lifespan} to 0 will not let this pair to expire.
func (ds DataStore) Set(key string, value interface{}, lifespan int64) interface{} {
	ds.mx.Lock()
	defer ds.mx.Unlock()
	_, ok := ds.data[key]
	if ok {
		delete(ds.data, key)
	}
	if lifespan > 0 {
		lifespan = time.Now().Unix() + lifespan
	}
	ds.data[key] = mapValue{
		expiration: lifespan,
		value:      value,
	}
	return ds
}

// Remove removes the enrty of the provided {key} if found
func (ds DataStore) Remove(key string) interface{} {
	ds.mx.Lock()
	defer ds.mx.Unlock()
	_, ok := ds.data[key]
	if ok {
		delete(ds.data, key)
	}
	return ds
}

func (ds DataStore) checkAndDeleteExpiredKeys() {
	for true {
		ds.mx.Lock()
		startTime := time.Now()
		for key, data := range ds.data {
			now := time.Now().Unix()
			if data.expiration != 0 && data.expiration <= now {
				delete(ds.data, key)
			}
		}
		ds.mx.Unlock()
		if time.Since(startTime) < ds.deletePrecision {
			time.Sleep(ds.deletePrecision - time.Since(startTime))
		}
	}

}
