// A simple, speedy, and better pseudo-random number generator.
//
// xer was an algorithm designed with the explicit purpose of being fast,
// simple, and reversible. xer2 is intended to be faster, simpler, and still
// reversible. Having done a fairly large amount of study on PRNGs, I've
// learned that xer is essentially an n-tap lagged Fibonacci generator. This
// has allowed me to take advantage of some of the properties of LFGs to
// improve the algorithm. For one thing, having n taps is probably a bad idea.
// It means that a larger state size does not necessarily mean a longer
// period. So xer2 has only one tap. xer also used xor as its binary op. I
// liked xor because it's completely symmetrical and localized. I don't like
// xor because it's completely symmetrical and localized. xer2 uses addition
// because addition causes the elements of each vector to interfere with each
// other - the carry operation. But I don't like addition because the carry
// operation means that the probability of each element in a random vector
// changing when added to another random vector increases with that element's
// order; higher bits are more likely to be carried into than lower ones, with
// bit 0 having P=0. So I fixed this problem. All you need to do is change
// which bits receive the carry results, so permute their order. The easiest
// bit permutation to carry out in a processor is rotation. xer used rotation
// right by popcount, but unless you're taking only very few samples, there's
// no difference between that and just rotating once. So that is what xer2
// does.
//
// So basically, xer2 is an LFG that cyclically rotates each step. It's
// probably slightly slower than a basic {j, k} additive LFG, but it should
// recover from bad states more quickly.
//
// I mentioned that one of the goals of xer was reversibility. If you have a
// given xer state, you can travel backwards and recover every value it's
// produced simply by changing the direction of rotation and decrementing the
// feed counter instead of incrementing it. xer2 is reversed similarly, but
// has the added step of changing the operation from addition to subtraction.
package xer2

// xer2.Source satisfies rand.Source but also provides the reversal operation
// and the option to use all bits generated.
type Source struct {
	feed, tap int
	state     []uint64
}

// Create a new xer2 source with n elements and d elements between feed and
// tap. The seed is extended from the given 64-bit value to the n-element IV
// via a linear congruential generator.
func New(n, d int, seed int64) (x *Source) {
	x = &Source{d, 0, make([]uint64, n)}
	x.Seed(seed)
	return x
}

// Create a new xer2 source with d elements between feed and tap using the
// given state iv. The number of elements will equal the length of the IV.
func NewIV(d int, iv []uint64) (x *Source) {
	x = &Source{d, 0, iv}
	return x
}

// Produce the full 64 bits of one iteration.
func (x *Source) Uint64() uint64 {
	x.feed++; x.tap++
	if x.feed >= len(x.state) {
		x.feed = 0
	} else if x.tap >= len(x.state) {
		x.tap = 0
	}
	sum := x.state[x.feed] + x.state[x.tap]
	sum = (sum >> 1) | (sum << 63)
	x.state[x.feed] = sum
	return sum
}

// Produce the lower 63 bits of one iteration.
func (x *Source) Int64() int64 {
	return int64(x.Uint64() & 0x7fffffffffffffff)
}

// Seed a xer2 source using Knuth's MMIX LCG to fill the state.
func (x *Source) Seed(seed int64) {
	const a uint64 = 6364136223846793005
	const c uint64 = 1442695040888963407
	const fuse = 20
	s := uint64(seed)
	for i := 0; i < fuse; i++ {
		s = a*s + c
		s = a*s + c
	}
	for i := range x.state {
		s = a*s + c
		x.state[i] = s >> 32
		s = a*s + c
		x.state[i] |= s & 0xffffffff00000000
	}
	x.feed = (x.feed - x.tap + len(x.state)) % len(x.state)
	x.tap = 0
	for _ = range x.state {
		x.Uint64()
	}
}

// Set the state directly. The feed and tap will be reset. Panics if the
// lengths do not match.
func (x *Source) SetState(state []uint64) {
	if len(state) != len(x.state) {
		panic("xer2 length mismatch when setting state")
	}
	x.state = state
	x.feed = (x.feed - x.tap + len(x.state)) % len(x.state)
	x.tap = 0
}

// Retrieve a copy of the state. It will be rotated such that the current tap
// is at element 0. (NewIV(d, x.SaveState()) will produce a generator
// equivalent to x.)
func (x *Source) SaveState() (state []uint64) {
	state = make([]uint64, len(x.state))
	copy(state, x.state[x.tap:])
	copy(state[len(x.state)-x.tap:], x.state[:x.tap])
	return state
}

/* why doesn't this work
// Take a step backward through the generator. This will produce the previous
// result of Uint64(), and calling Uint64() again after this will produce
// the same value again.
func (x *Source) Reverse() uint64 {
	sum := x.state[x.feed]
	sum = (sum << 1) | (sum >> 63)
	x.state[x.feed] = sum - x.state[x.tap]
	x.tap--; x.feed--
	if x.tap < 0 {
		x.tap = len(x.state) - 1
	} else if x.feed < 0 {
		x.feed = len(x.state) - 1
	}
	return sum
}
*/
