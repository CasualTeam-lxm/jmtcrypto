// combinatorics.go
package mce

import (
	"math/rand"
)

func Choose(n, k int) int {
	if k > n {
		panic("Choose: k > n")
	}
	if k < 0 {
		panic("Choose: k < 0")
	}
	if n <= 1 || k == 0 || n == k {
		return 1
	}
	if newK := n - k; newK < k {
		k = newK
	}
	if k == 1 {
		return n
	}
	// Our return value, and this allows us to skip the first iteration.
	ret := n - k + 1
	for i, j := ret+1, 2; j <= k; i, j = i+1, j+1 {
		ret = ret * i / j
	}
	return ret
}

// Set Difference: A - B
func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func InvertIndices(m int, is [][]int) ([][]int) {
	all := []int{}
	for i := 1; i <= m; i++ {
		all = append(all, i)
	}	

	for i, v := range is {
		is[i] = Difference(all, v)
	}

	return is
}

func rPool(p int, n []int, c []int, cc [][]int) [][]int {
	if len(n) == 0 || p <= 0 {
		return cc
	}
	p--
	for i := range n {
		r := make([]int, len(c)+1)
		copy(r, c)
		r[len(r)-1] = n[i]
		if p == 0 {
			cc = append(cc, r)
		}
		cc = rPool(p, n[i+1:], r, cc)
	}
	return cc
}

func Pool(p int, n []int) [][]int {
	return rPool(p, n, nil, nil)
}

func RandomPermutaion(n int) []int {
	return rand.Perm(n)
}

func Log2(n int) (l int) {
	l = -1
	for n != 0 {
		n >>= 1
		l++
	}
	return
}

func GetWedges(set []Bitset, depth int) []Bitset {
		ones := make(Bitset, len(set[0]))
		for i := range ones {
				ones[i] = true
		}

		return GetWedgesHelper(set, depth, 0, ones, []Bitset{})
}

func GetWedgesHelper(set []Bitset, depth int, start int,
                      product Bitset, accum []Bitset) []Bitset {
	if depth == 0 {
		return append(accum, product)
	} else {
		for i := start; i <= len(set) - depth; i++ {
			accum = GetWedgesHelper(set, depth - 1, i + 1,
			BitsetAND(product, set[i]), accum)
		}
		return accum
	}
}

func AlternatingBitset(run, n int) (v Bitset) {
	v = make(Bitset, n)
	runcount := 0
	runflag := true

	for i := 0; i < n; i++ {
		if runflag {
			v[i] = true
		}

		runcount += 1
		if runcount == run {
			runflag = !runflag
			runcount = 0
		}
	}

	return
}
