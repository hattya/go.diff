//
// go.diff :: diff_test.go
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

package diff_test

import (
	"testing"

	"github.com/hattya/go.diff"
)

var tests = []struct {
	name string
	a, b []rune
	cl   []diff.Change
}{
	{
		name: "same",
		a:    []rune("abc"),
		b:    []rune("abc"),
		cl:   []diff.Change{},
	},
	{
		name: "other",
		a:    []rune("abc"),
		b:    []rune("12345"),
		cl: []diff.Change{
			{A: 0, B: 0, Del: 3, Ins: 5},
		},
	},
	{
		name: "shift",
		a:    []rune("abc"),
		b:    []rune("`abc"),
		cl: []diff.Change{
			{A: 0, B: 0, Del: 0, Ins: 1},
		},
	},
	{
		name: "unshift",
		a:    []rune("`abc"),
		b:    []rune("abc"),
		cl: []diff.Change{
			{A: 0, B: 0, Del: 1, Ins: 0},
		},
	},
	{
		name: "push",
		a:    []rune("abc"),
		b:    []rune("abcd"),
		cl: []diff.Change{
			{A: 3, B: 3, Del: 0, Ins: 1},
		},
	},
	{
		name: "pop",
		a:    []rune("abcd"),
		b:    []rune("abc"),
		cl: []diff.Change{
			{A: 3, B: 3, Del: 1, Ins: 0},
		},
	},
	{
		name: "overlap",
		a:    []rune("b"),
		b:    []rune("abc"),
		cl: []diff.Change{
			{A: 0, B: 0, Del: 0, Ins: 1},
			{A: 1, B: 2, Del: 0, Ins: 1},
		},
	},
	{
		name: "paper",
		a:    []rune("acbdeacbed"),
		b:    []rune("acebdabbabed"),
		cl: []diff.Change{
			{A: 2, B: 2, Del: 0, Ins: 1},
			{A: 4, B: 5, Del: 1, Ins: 0},
			{A: 6, B: 6, Del: 1, Ins: 0},
			{A: 8, B: 7, Del: 0, Ins: 3},
		},
	},
}
var paper = len(tests) - 1

func TestDiff(t *testing.T) {
	for _, tt := range tests {
		cl := diff.Runes(tt.a, tt.b)
		if g, e := len(cl), len(tt.cl); g != e {
			t.Errorf("%v: expected %v, got %v", tt.name, e, g)
		}
		for i, c := range tt.cl {
			if cl[i] != c {
				t.Errorf("%v: expected %#v, got %#v", tt.name, c, cl[i])
			}
		}
	}
}

func TestDiffExchange(t *testing.T) {
	for _, tt := range tests {
		cl := diff.Runes(tt.b, tt.a)
		if g, e := len(cl), len(tt.cl); g != e {
			t.Errorf("%v: expected %v, got %v", tt.name, e, g)
		}
		for i, c := range tt.cl {
			c = diff.Change{
				A:   c.B,
				B:   c.A,
				Del: c.Ins,
				Ins: c.Del,
			}
			if cl[i] != c {
				t.Errorf("%v: expected %#v, got %#v", tt.name, c, cl[i])
			}
		}
	}
}

func TestBytes(t *testing.T) {
	tt := tests[paper]
	cl := diff.Bytes(toByte(tt.a), toByte(tt.b))
	if g, e := len(cl), len(tt.cl); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	for i, c := range cl {
		if c != tt.cl[i] {
			t.Errorf("expected %#v, got %#v", tt.cl[i], c)
		}
	}
}

func TestInts(t *testing.T) {
	tt := tests[paper]
	cl := diff.Ints(toInt(tt.a), toInt(tt.b))
	if g, e := len(cl), len(tt.cl); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	for i, c := range cl {
		if c != tt.cl[i] {
			t.Errorf("expected %#v, got %#v", tt.cl[i], c)
		}
	}
}

func TestRunes(t *testing.T) {
	tt := tests[paper]
	cl := diff.Runes(tt.a, tt.b)
	if g, e := len(cl), len(tt.cl); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	for i, c := range cl {
		if c != tt.cl[i] {
			t.Errorf("expected %#v, got %#v", tt.cl[i], c)
		}
	}
}

func TestStrings(t *testing.T) {
	tt := tests[paper]
	cl := diff.Strings(toString(tt.a), toString(tt.b))
	if g, e := len(cl), len(tt.cl); g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	for i, c := range cl {
		if c != tt.cl[i] {
			t.Errorf("expected %#v, got %#v", tt.cl[i], c)
		}
	}
}

type runes struct {
	A, B []rune
}

func (d *runes) Equal(i, j int) bool { return d.A[i] == d.B[j] }

func BenchmarkDiff(b *testing.B) {
	tt := tests[paper]
	n := len(tt.a)
	m := len(tt.b)
	data := &runes{tt.a, tt.b}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.Diff(n, m, data)
	}
}

func BenchmarkBytes(b *testing.B) {
	tt := tests[paper]
	A := toByte(tt.a)
	B := toByte(tt.b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.Bytes(A, B)
	}
}

func BenchmarkInts(b *testing.B) {
	tt := tests[paper]
	A := toInt(tt.a)
	B := toInt(tt.b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.Ints(A, B)
	}
}

func BenchmarkRunes(b *testing.B) {
	tt := tests[paper]
	A := tt.a
	B := tt.b
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.Runes(A, B)
	}
}

func BenchmarkStrings(b *testing.B) {
	tt := tests[paper]
	A := toString(tt.a)
	B := toString(tt.b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.Strings(A, B)
	}
}

func toByte(a []rune) []byte {
	l := make([]byte, len(a))
	for i, r := range a {
		l[i] = byte(r)
	}
	return l
}

func toInt(a []rune) []int {
	l := make([]int, len(a))
	for i, r := range a {
		l[i] = int(r)
	}
	return l
}

func toString(a []rune) []string {
	l := make([]string, len(a))
	for i, r := range a {
		l[i] = string(r)
	}
	return l
}
