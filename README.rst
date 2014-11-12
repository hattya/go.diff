go.diff
=======

A Go_ implementation of "An O(NP) Sequence Comparison Algorithm" [#]_.

.. image:: https://godoc.org/github.com/hattya/go.diff?status.svg
   :target: https://godoc.org/github.com/hattya/go.diff

.. image:: https://drone.io/github.com/hattya/go.diff/status.png
   :target: https://drone.io/github.com/hattya/go.diff/latest

.. _Go: http://golang.org/


Install
-------

.. code:: console

   $ go get -u github.com/hattya/go.diff


Usage
-----

.. code:: go

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


License
-------

go.diff is distributed under the terms of the MIT License.


.. rubric:: Footnotes

.. [#] \S. Wu, U. Manber, G. Myers and W. Miller, "An O(NP) Sequence Comparison Algorithm" August 1989
