package kvse

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestKvseNew(t *testing.T) {
	want := &DataStore{
		data:            map[string]mapValue{},
		mx:              sync.RWMutex{},
		deletePrecision: time.Millisecond,
	}
	got := New(time.Millisecond)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("New kvse %v, want %v", got, want)
	}

	want = &DataStore{
		data:            map[string]mapValue{},
		mx:              sync.RWMutex{},
		deletePrecision: time.Second,
	}
	got = New(0)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("New kvse %v, want %v", got, want)
	}
}

func TestKvseSet(t *testing.T) {
	kvses := &DataStore{
		data:            map[string]mapValue{},
		mx:              sync.RWMutex{},
		deletePrecision: time.Millisecond,
	}
	for i := 1; i <= 10; i++ {
		kvses.Set(strconv.Itoa(i), i*2, 0)
	}
	for i := 1; i <= 10; i++ {
		if kvses.data[strconv.Itoa(i)].value != 2*i {
			t.Errorf("Set kvse failed. Added at %s: %d, got: %d.", strconv.Itoa(i), 2*i, kvses.data["1"].value)
		}
	}
}

func TestKvseGet(t *testing.T) {
	kvses := &DataStore{
		data:            map[string]mapValue{},
		mx:              sync.RWMutex{},
		deletePrecision: time.Millisecond,
	}
	for i := 1; i <= 10; i++ {
		kvses.data[strconv.Itoa(i)] = mapValue{
			expiration: time.Now(),
			value:      i * 2,
		}
	}
	for i := 1; i <= 10; i++ {
		if val, ok := kvses.Get(strconv.Itoa(i)); !ok || (ok && val != i*2) {
			t.Errorf("Get kvse failed. Added at %s: %d, got: %d.", strconv.Itoa(i), 2*i, kvses.data[strconv.Itoa(i)].value)
		}
	}
}

func TestKvseSetExpiration(t *testing.T) {
	kvses := New(time.Second)
	for i := 1; i <= 10; i++ {
		kvses.Set(strconv.Itoa(i), i*2, 500*time.Millisecond)
	}
	time.Sleep(10 * time.Millisecond)
	for i := 1; i <= 10; i++ {
		if data, ok := kvses.data[strconv.Itoa(i)]; !ok || (ok && data.value != i*2) {
			t.Errorf("Persisted kvse failed. Added at %s: %d, and didn't got it back.", strconv.Itoa(i), 2*i)
		}
	}
	time.Sleep(time.Second)
	for i := 1; i <= 10; i++ {
		if _, ok := kvses.data[strconv.Itoa(i)]; ok {
			t.Errorf("Removing kvse failed. Added at %s: %d, and got it back.", strconv.Itoa(i), 2*i)
		}
	}
}

func TestMultipleKvses(t *testing.T) {
	kvses1 := New(time.Second)
	kvses2 := New(time.Second)
	for i := 1; i <= 10; i++ {
		kvses1.Set(strconv.Itoa(i), i*2, time.Second)
		kvses2.Set(strconv.Itoa(i), i*3, time.Second)
	}
	for i := 1; i <= 10; i++ {
		if val, ok := kvses1.Get(strconv.Itoa(i)); !ok || (ok && val != i*2) {
			t.Errorf("Multiple kvses failed. Added at %s: %d, got: %d.", strconv.Itoa(i), 2*i, val)
		}
		if val, ok := kvses2.Get(strconv.Itoa(i)); !ok || (ok && val != i*3) {
			t.Errorf("Multiple kvses failed. Added at %s: %d, got: %d.", strconv.Itoa(i), 3*i, val)
		}
	}
}
