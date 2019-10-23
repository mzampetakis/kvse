# kvse
An in-memory Key Value Store with Expiration

## Installing
```
$ go get -u github.com/mzampetakis/kvse
```

## Usage
### `kvse.New`
- `New() DataStore`
- `New` returns a new instance of a `DataStore`.

### `kvse.Set`
- `Set(key string, value interface, lifespan time.Duration)`
- `Set` adds the `value` to the data store associated to `key` and will be deleted after `lifespan` duration.

### `kvse.Get`
- `Get(key string) interface{}, bool`
- `Get` retrieves the value associated to `key`, and a boolean variable if found.

### `kvse.Has`
- `Has(key string) bool`
- `Has` returns a `bool` based on whether or not `key` exists in the data store. 

## Running the tests
```
go test
```

## Running the benchmarks
```
go test -bench=.    
```
