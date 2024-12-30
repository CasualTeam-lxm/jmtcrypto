// matrix.go
package mce

import (
	"fmt"
	"math/rand"
	"time"
)

func Identity(n int) (Id []Bitset) {
	for i := 0; i < n; i++ {
		b := make(Bitset, n)
		b[i] = true

		Id = append(Id, b)
	}

	return
}

func MatrixPair(n int) ([]Bitset, []Bitset) {
	Id    := Identity(n)

	C     := make([]Bitset, len(Id))
	C_inv := make([]Bitset, len(Id))

	perm := RandomPermutaion(n)

	// doing this inplace in more trouble than it is worth!
	for i := 0; i < len(perm); i++ {
		C[perm[i]] = Id[i]
		C_inv[i]   = Id[i]
	}
	
	seed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(seed)

	swaps := n*n
	is, js:= []int{}, []int{}
	for len(is) < swaps {
		i := rng.Intn(n)
		j := rng.Intn(n)
		if i != j {
			is = append(is, i)
			js = append(js, j)
		}
	}

	// apply row operations
	for p := 0; p < len(is); p++ {
		i := is[p]
		j := js[p]
		C[i] = BitsetXOR(C[i], C[j])

		i = is[swaps - p - 1]
		j = js[swaps - p - 1]
		C_inv[i] = BitsetXOR(C_inv[i], C_inv[j])
	}

	copy(Id, C_inv)

	// doing this inplace in more trouble than it is worth!
	for i := 0; i < len(perm); i++ {
		C_inv[i] = Id[perm[i]]
	}

	return C, C_inv
	
}

func Column(M []Bitset, c int) (col Bitset) {
	col = make(Bitset, len(M))

	for i := 0; i < len(M); i++ {
		col[i] = M[i][c]
	}

	return
}

func MatMulMat(M []Bitset, C []Bitset) (C_new []Bitset) {
	C_new = make([]Bitset, len(M))

	for j, row := range M {
		rowNew := make(Bitset, len(C[0]))
		for i := 0; i < len(C[0]); i++ {
			col := Column(C, i)
			rowNew[i] = BitsetDOT(row, col)
		}	
		C_new[j] = rowNew
	}

	return
}

func MatMulVecR(V Bitset, M []Bitset) (Bitset) {
	V_new := make(Bitset, len(V))

	for i := 0; i < len(M[0]); i++ {
		col := Column(M, i)
		V_new[i] = BitsetDOT(V, col)
	}	

	return V_new
}

func PrintMatrix(M []Bitset) {
	for _, r := range M {
		PrintBin(r, true)
	}
	fmt.Printf("\n")
}
