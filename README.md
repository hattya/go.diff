# go.diff

A [Go](https://go.dev/) implementation of "An O(NP) Sequence Comparison Algorithm".

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/hattya/go.diff)](https://pkg.go.dev/github.com/hattya/go.diff)
[![GitHub Actions](https://github.com/hattya/go.diff/actions/workflows/ci.yml/badge.svg)](https://github.com/hattya/go.diff/actions/workflows/ci.yml)
[![Appveyor](https://ci.appveyor.com/api/projects/status/ryyeqn70w488ac8f/branch/master?svg=true)](https://ci.appveyor.com/project/hattya/go-diff)
[![Codecov](https://codecov.io/gh/hattya/go.diff/branch/master/graph/badge.svg)](https://codecov.io/gh/hattya/go.diff)


## Installation

```console
$ go get -u github.com/hattya/go.diff
```


## Usage

```go
package main

import (
	"fmt"

	"github.com/hattya/go.diff"
)

func main() {
	a := []rune("acbdeacbed")
	b := []rune("acebdabbabed")
	cl := diff.Runes(a, b)
	i := 0
	for _, c := range cl {
		for ; i < c.A; i++ {
			fmt.Printf("  %c\n", a[i])
		}
		for j := c.A; j < c.A+c.Del; j++ {
			fmt.Printf("- %c\n", a[j])
		}
		for j := c.B; j < c.B+c.Ins; j++ {
			fmt.Printf("+ %c\n", b[j])
		}
		i += c.Del
	}
	for ; i < len(a); i++ {
		fmt.Printf("  %c\n", a[i])
	}
}
```


## License

go.diff is distributed under the terms of the MIT License.


## References

S. Wu, U. Manber, G. Myers, and W. Miller, "An O(NP) Sequence Comparison Algorithm" August 1989
