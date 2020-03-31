//
// go.diff :: diff.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// Package diff implements the difference algorithm, which is based upon
// S. Wu, U. Manber, G. Myers, and W. Miller,
// "An O(NP) Sequence Comparison Algorithm" August 1989.
package diff

type Interface interface {
	// Equal returns whether the elements at i and j are equal.
	Equal(i, j int) bool
}

// Bytes returns the differences between byte slices.
func Bytes(a, b []byte) []Change {
	return Diff(len(a), len(b), &bytes{a, b})
}

type bytes struct {
	A, B []byte
}

func (p *bytes) Equal(i, j int) bool { return p.A[i] == p.B[j] }

// Ints returns the differences between int slices.
func Ints(a, b []int) []Change {
	return Diff(len(a), len(b), &ints{a, b})
}

type ints struct {
	A, B []int
}

func (p *ints) Equal(i, j int) bool { return p.A[i] == p.B[j] }

// Runes returns the differences between rune slices.
func Runes(a, b []rune) []Change {
	return Diff(len(a), len(b), &runes{a, b})
}

type runes struct {
	A, B []rune
}

func (p *runes) Equal(i, j int) bool { return p.A[i] == p.B[j] }

// Strings returns the differences between string slices.
func Strings(a, b []string) []Change {
	return Diff(len(a), len(b), &strings{a, b})
}

type strings struct {
	A, B []string
}

func (p *strings) Equal(i, j int) bool { return p.A[i] == p.B[j] }

// Diff returns the differences between data.
// It makes O(NP) (the worst case) calls to data.Equal.
func Diff(m, n int, data Interface) []Change {
	c := &context{data: data}
	if n >= m {
		c.M = m
		c.N = n
	} else {
		c.M = n
		c.N = m
		c.xchg = true
	}
	c.Δ = c.N - c.M
	return c.compare()
}

type Change struct {
	A, B int // position in a and b
	Del  int // number of elements that deleted from a
	Ins  int // number of elements that inserted into b
}

type context struct {
	data Interface
	M, N int
	Δ    int
	fp   []point
	xchg bool
}

func (c *context) compare() []Change {
	c.fp = make([]point, (c.M+1)+(c.N+1)+1)
	for i := range c.fp {
		c.fp[i].y = -1
	}

	Δ := c.Δ + (c.M + 1)
	for p := 0; c.fp[Δ].y != c.N; p++ {
		for k := -p; k < c.Δ; k++ {
			c.snake(k)
		}
		for k := c.Δ + p; k > c.Δ; k-- {
			c.snake(k)
		}
		c.snake(c.Δ)
	}

	lcs, n := c.reverse(c.fp[Δ].lcs)
	cl := make([]Change, 0, n+1)
	x, y := 0, 0
	for ; lcs != nil; lcs = lcs.next {
		if x < lcs.x || y < lcs.y {
			if !c.xchg {
				cl = append(cl, Change{x, y, lcs.x - x, lcs.y - y})
			} else {
				cl = append(cl, Change{y, x, lcs.y - y, lcs.x - x})
			}
		}
		x = lcs.x + lcs.n
		y = lcs.y + lcs.n
	}
	if x < c.M || y < c.N {
		if !c.xchg {
			cl = append(cl, Change{x, y, c.M - x, c.N - y})
		} else {
			cl = append(cl, Change{y, x, c.N - y, c.M - x})
		}
	}
	return cl
}

func (c *context) snake(k int) {
	var y int
	var prev *lcs
	kk := k + (c.M + 1)

	h := &c.fp[kk-1]
	v := &c.fp[kk+1]
	if h.y+1 >= v.y {
		y = h.y + 1
		prev = h.lcs
	} else {
		y = v.y
		prev = v.lcs
	}

	x := y - k
	n := 0
	for x < c.M && y < c.N {
		var eq bool
		if !c.xchg {
			eq = c.data.Equal(x, y)
		} else {
			eq = c.data.Equal(y, x)
		}
		if !eq {
			break
		}
		x++
		y++
		n++
	}

	p := &c.fp[kk]
	p.y = y
	if n == 0 {
		p.lcs = prev
	} else {
		p.lcs = &lcs{
			x:    x - n,
			y:    y - n,
			n:    n,
			next: prev,
		}
	}
}

func (c *context) reverse(curr *lcs) (*lcs, int) {
	n := 0
	var prev, next *lcs
	for ; curr != nil; n++ {
		prev = curr.next
		curr.next = next
		next = curr
		curr = prev
	}
	return next, n
}

type point struct {
	y   int
	lcs *lcs
}

type lcs struct {
	x, y int
	n    int
	next *lcs
}
