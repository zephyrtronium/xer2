package xer2

import "testing"

func TestItWorks(t *testing.T) {
	x := New(17, 10, 0)
	for i := 0; i < 100; i++ {
		x.Uint64()
	}
}

func TestSeedingIsConsistent(t *testing.T) {
	x, y := New(17, 10, 0), New(17, 10, 0)
	t.Log("starting")
	for i := 0; i < 1000; i++ {
		if x.Uint64() != y.Uint64() {
			t.Error("seeding same constants was dumb at", i)
		}
	}
	t.Log("did same constant")
	x.SetState(y.SaveState())
	for i := 0; i < 1000; i++ {
		if x.Uint64() != y.Uint64() {
			t.Error("state restoration failed at", i)
		}
	}
	t.Log("did restoration")
	x.Seed(0); y.Seed(0)
	for i := 0; i < 1000; i++ {
		if x.Uint64() != y.Uint64() {
			t.Error("constant seeding after use failed at", i)
		}
	}
}

func Benchmark607_334(b *testing.B) {
	x := New(607, 334, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Uint64()
	}
}

func Benchmark17_10(b *testing.B) {
	x := New(17, 10, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Uint64()
	}
}

/*
func TestReversibility(t *testing.T) {
	x := New(17, 10, 0)
	for i := 0; i < 100; i++ {
		//t.Log(x.feed, x.state)
		a, b := x.Uint64(), x.Reverse()
		if a != b {
			t.Errorf("%v %016x %016x", i, a, b)
		}
		x.Uint64()
	}
}
*/