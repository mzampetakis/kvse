// Package kvse is an In memory Key-Value Store with Expiration for each pair
// It provides the necessary structs and methods to maintain a key-value
// map with expiration for each pair.
// Expiration is managed by the lib itself but can be omitted if set to 0.
// Also it supports a minimum precision you want to achieve for deleting expired pairs.
// Providing a 0 precision will cause to use the default one which is time.second.
package kvse

import (
	"fmt"
	"sync"
	"time"
)

// DataStore stores the struct of each list
// Holds a pointer to the list, a mutex lock and the list itself
type DataStore struct {
	data            map[string]mapValue
	mx              sync.RWMutex
	deletePrecision time.Duration
}

type mapValue struct {
	expiration time.Time
	value      interface{}
}

// New func create`s and returns a new kvse DataStore instance
// and initiates a worker to check and delete for expired keys.
// {precision} is the minimum precision you want to achieve for deleting expired pairs.
// Providing the precision as 0 will use a default one: time.second.
func New(precision time.Duration) *DataStore {
	if precision == 0 {
		precision = time.Second
	}
	ds := DataStore{
		data:            map[string]mapValue{},
		mx:              sync.RWMutex{},
		deletePrecision: precision,
	}
	go ds.deleteExpiredKeys()
	return &ds
}

// Has returns a boolean based on whether or not the store contains a value for
// {key}.
func (ds *DataStore) Has(key string) bool {
	ds.mx.RLock()
	defer ds.mx.RUnlock()
	_, ok := ds.data[key]
	return ok
}

// Get retrieves the value associated to the {key} in the store.
func (ds *DataStore) Get(key string) (interface{}, bool) {
	ds.mx.RLock()
	defer ds.mx.RUnlock()
	data, ok := ds.data[key]
	return data.value, ok
}

// Set adds a ne value to a specific key with a {lifespan} duration.
// Setting the {lifespan} to 0 will not let this pair to expire.
func (ds *DataStore) Set(key string, value interface{}, lifespan time.Duration) {
	ds.mx.Lock()
	defer ds.mx.Unlock()
	delete(ds.data, key)
	var expire time.Time
	if lifespan.Nanoseconds() != 0 {
		expire = time.Now().Add(lifespan)
	}
	ds.data[key] = mapValue{
		expiration: expire,
		value:      value,
	}
}

// Remove removes the enrty of the provided {key} if found
func (ds *DataStore) Remove(key string) {
	ds.mx.Lock()
	defer ds.mx.Unlock()
	delete(ds.data, key)
}

func (ds *DataStore) deleteExpiredKeys() {
	for {
		startTime := time.Now()
		ds.checkAndDeleteExpiredKeys()
		if time.Since(startTime) < ds.deletePrecision {
			time.Sleep(ds.deletePrecision - time.Since(startTime))
		}
	}

}

func (ds *DataStore) checkAndDeleteExpiredKeys() {
	ds.mx.Lock()
	defer ds.mx.Unlock()
	for key, data := range ds.data {
		now := time.Now()
		if !data.expiration.IsZero() && data.expiration.Before(now) {
			delete(ds.data, key)
		}
	}
}

func (ds *DataStore) String() string {
	str := "PRINTING KVSE\n"
	ds.mx.RLock()
	defer ds.mx.RUnlock()
	for key, data := range ds.data {
		str = str + fmt.Sprintf(" Key %s \t Val: %d \t Exp: %s \n", key, data.value, data.expiration)
	}
	return str
}
