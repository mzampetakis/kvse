package kvse

import (
	"strconv"
	"testing"
	"time"
)

var exists bool

type benchCase struct {
	name      string
	times     int
	precision time.Duration
	lifetime  time.Duration
}

var benchCases = []benchCase{
	{"10 times|1000 Sec precision|1 Sec lifetime", 10, 1000 * time.Second, time.Second},
	{"10 times|1000 Sec precision|1 ms lifetime", 10, 1000 * time.Second, time.Millisecond},
	{"10 times|1000 Sec precision|1 μs lifetime", 10, 1000 * time.Second, time.Microsecond},

	{"100 times|1000 Sec precision|1 Sec lifetime", 100, 1000 * time.Second, time.Second},
	{"100 times|1000 Sec precision|1 ms lifetime", 100, 1000 * time.Second, time.Millisecond},
	{"100 times|1000 Sec precision|1 μs lifetime", 100, 1000 * time.Second, time.Microsecond},

	{"1k times|1000 Sec precision|1 Sec lifetime", 1000, 1000 * time.Second, time.Second},
	{"1k times|1000 Sec precision|1 ms lifetime", 1000, 1000 * time.Second, time.Millisecond},
	{"1k times|1000 Sec precision|1 μs lifetime", 1000, 1000 * time.Second, time.Microsecond},

	{"10 times|1 Sec precision|1 Sec lifetime", 10, time.Second, time.Second},
	{"10 times|1 Sec precision|1 ms lifetime", 10, time.Second, time.Millisecond},
	{"10 times|1 Sec precision|1 μs lifetime", 10, time.Second, time.Microsecond},

	{"100 times|1 Sec precision|1 Sec lifetime", 100, time.Second, time.Second},
	{"100 times|1 Sec precision|1 ms lifetime", 100, time.Second, time.Millisecond},
	{"100 times|1 Sec precision|1 μs lifetime", 100, time.Second, time.Microsecond},

	{"1k times|1 Sec precision|1 Sec lifetime", 1000, time.Second, time.Second},
	{"1k times|1 Sec precision|1 ms lifetime", 1000, time.Second, time.Millisecond},
	{"1k times|1 Sec precision|1 μs lifetime", 1000, time.Second, time.Microsecond},

	{"10 times|1 ms precision|1 Sec lifetime", 10, time.Millisecond, time.Second},
	{"10 times|1 ms precision|1 ms lifetime", 10, time.Millisecond, time.Millisecond},
	{"10 times|1 ms precision|1 μs lifetime", 10, time.Millisecond, time.Microsecond},

	{"100 times|1 ms precision|1 Sec lifetime", 100, time.Millisecond, time.Second},
	{"100 times|1 ms precision|1 ms lifetime", 100, time.Millisecond, time.Millisecond},
	{"100 times|1 ms precision|1 μs lifetime", 100, time.Millisecond, time.Microsecond},

	{"1k times|1 ms precision|1 Sec lifetime", 1000, time.Millisecond, time.Second},
	{"1k times|1 ms precision|1 ms lifetime", 1000, time.Millisecond, time.Millisecond},
	{"1k times|1 ms precision|1 μs lifetime", 1000, time.Millisecond, time.Microsecond},

	{"10 times|1 μs precision|1 Sec lifetime", 10, time.Microsecond, time.Second},
	{"10 times|1 μs precision|1 ms lifetime", 10, time.Microsecond, time.Millisecond},
	{"10 times|1 μs precision|1 μs lifetime", 10, time.Microsecond, time.Microsecond},

	{"100 times|1 μs precision|1 Sec lifetime", 100, time.Microsecond, time.Second},
	{"100 times|1 μs precision|1 ms lifetime", 100, time.Microsecond, time.Millisecond},
	{"100 times|1 μs precision|1 μs lifetime", 100, time.Microsecond, time.Microsecond},

	{"1k times|1 μs precision|1 Sec lifetime", 1000, time.Microsecond, time.Second},
	{"1k times|1 μs precision|1 ms lifetime", 1000, time.Microsecond, time.Millisecond},
	{"1k times|1 μs precision|1 μs lifetime", 1000, time.Microsecond, time.Microsecond},
}

func BenchmarkKvseSet(b *testing.B) {
	for _, benchCase := range benchCases {
		kvses := New(benchCase.precision)
		b.Run(benchCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < benchCase.times; i++ {
					kvses.Set(strconv.Itoa(i), i*2, time.Hour)
				}
			}
		})
	}
}

func BenchmarkKvseGet(b *testing.B) {
	for _, benchCase := range benchCases {
		kvses := New(benchCase.precision)
		for n := 0; n < b.N; n++ {
			for i := 0; i < benchCase.times; i++ {
				kvses.Set(strconv.Itoa(i), i*2, time.Hour)
			}
		}
		b.Run(benchCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < benchCase.times; i++ {
					_, exists = kvses.Get(strconv.Itoa(i))
				}
			}
		})
	}
}

func BenchmarkKvseHas(b *testing.B) {
	for _, benchCase := range benchCases {
		kvses := New(benchCase.precision)
		for n := 0; n < b.N; n++ {
			for i := 0; i < benchCase.times; i++ {
				kvses.Set(strconv.Itoa(i), i*2, time.Hour)
			}
		}
		b.Run(benchCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < benchCase.times; i++ {
					exists = kvses.Has(strconv.Itoa(i))
				}
			}
		})
	}
}

func BenchmarkKvseRemove(b *testing.B) {
	for _, benchCase := range benchCases {
		kvses := New(benchCase.precision)
		for n := 0; n < b.N; n++ {
			for i := 0; i < benchCase.times; i++ {
				kvses.Set(strconv.Itoa(i), i*2, time.Hour)
			}
		}
		b.Run(benchCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < benchCase.times; i++ {
					kvses.Remove(strconv.Itoa(i))
				}
			}
		})
	}
}
