// reedmuller.go
package mce

import (
	"fmt"
)

type RMCode struct {
	r,m      int
	inBits   int
	outBits  int
	M        []Bitset
	diffs    [][]int
	errors   int
}

func ReedMuller(r, m int) RMCode {
	k := 0
	n := 1 << m
	rm := []Bitset{}
	diffs := [][]int{}

	for i := 0; i <= r; i++ {
		k += Choose(m, i)
	}

	counter := 0
	indices := [][]int{}
	ordinals := []int{}

	for i := n; i > 0; i = i / 2 {
		indices = append(indices, []int{counter})
		rm = append(rm, AlternatingBitset(i, n))
		if counter != 0 {
			ordinals = append(ordinals, counter)
		}

		counter++
	}

	for i := 2; i <= r; i++ {
		additional := GetWedges(rm[1:(m+1)], i)
		newIndices := Pool(i, ordinals)
		
		for _, v := range additional {
			rm = append(rm, v)
		}

		for _, v := range newIndices {
			indices = append(indices, v)
		}

	}

	diffs = InvertIndices(m, indices)
	inBits  := k
	outBits := 1 << m
	
	return RMCode{r:r, m:m, M:rm, diffs:diffs, inBits:inBits, outBits:outBits}
}

func (rm RMCode) Encrypt(msg Bitset, addErrors bool) (ctxt Bitset) {
	N := rm.inBits // the number of bit in each plaintext Block
	P := rm.outBits // the number of bytes in each cipher text block

	ctxt = make(Bitset, 0)

	for i := 0; i < len(msg); i += N {
		cipherBlock := make(Bitset, P)
		for b := 0; b < N; b++ {
			if msg[(i + b)] {
				cipherBlock = BitsetXOR(cipherBlock, rm.M[b])
			}
		}
		
		ctxt = append(ctxt, cipherBlock[:]...)
		cipherBlock = make(Bitset, P)
	}

	if addErrors {
		errors := ((1 << (rm.m - rm.r)) - 1) / 2
		bits := (1 << rm.m)
		fmt.Printf("Adding %d errors for every %d bits (%f bits per error).\n\n",
		           errors, bits, float64(bits)/float64(errors))

		ctxt = AddErrors(ctxt, errors, bits)
	}

	return
}

func (rm RMCode) Decrypt(msg Bitset, fixErrors bool) (ptxt Bitset) {
	P := rm.outBits
	// get the characteristic vectors
	charVectors := [][]Bitset{}
	for i := 0; i < len(rm.M); i++ {
		charVectors = append(charVectors, getCharVectors(rm, i))
	}

	for i := 0; i < len(msg); i += P {

		eword := make(Bitset, P)
		copy(eword, msg[i:i+P])
		ewordTemp := make(Bitset, len(eword))
		copy(ewordTemp, eword)

		coeffs := make([]bool, len(rm.M))

		// compare this block to char vectors for each index
		// iterate backwards through charVectors
		for j := len(charVectors) - 1; j >= 0; j-- {
			chrVecs := charVectors[j]
			votesForOne := 0
			for _, cv := range chrVecs {
				if BitsetDOT(cv, eword) {
					votesForOne += 1
				}
			}

			if votesForOne == len(charVectors[j]) - votesForOne {
				fmt.Println("DANGER!")
			}

			if fixErrors {
				if votesForOne > len(charVectors[j]) - votesForOne {
					ewordTemp = BitsetXOR(rm.M[j], ewordTemp)
					coeffs[j] = true
				}
			} else {
				if votesForOne == len(charVectors[j]) {
					ewordTemp = BitsetXOR(rm.M[j], ewordTemp)
					coeffs[j] = true
				} 				
			}

			if (j == 0) || (len(rm.diffs[j]) != len(rm.diffs[j-1])) {
				copy(eword, ewordTemp)
			}
		}

		flag := BitsetVote(eword, true)

		if flag {
			coeffs = BitsetFlipTopBit(coeffs)
		}

		ptxt = append(ptxt, coeffs...)
	}

	return
}

func getCharVectors(rm RMCode, row int) (chars []Bitset) {
	initial := BitsetAllOnes(rm.outBits)
	chars = []Bitset{ initial }

	for _, index := range rm.diffs[row] {
		fold := rm.M[index] // grab the ith row of the r-m matrix
		notFold := InvertBitset(fold) //

		temp := make([]Bitset, 2*len(chars))
		for i, v := range chars {
			temp[2*i]     = BitsetAND(v, fold)
			temp[2*i + 1] = BitsetAND(v, notFold)
		}

		chars = temp
	}

	return
}

// 'n' errors per 'k' bytes
func AddErrors(ctext Bitset, n, k int) Bitset {
	ctextErr := make(Bitset, len(ctext))
	copy(ctextErr, ctext)

	for blockId := 0; blockId < len(ctext); blockId += k {
		errors := RandomPermutaion(k)
		for errIndex := 0; errIndex < n; errIndex++ {
			ctextErr[blockId + errors[errIndex]] = !ctextErr[blockId + errors[errIndex]]
		}
	}
	return ctextErr 
}

func (rm RMCode) PermuteCols(perm []int) (RMCode) {
	mNew := make([]Bitset, len(rm.M))

	for i := 0; i < len(mNew); i++ {
		mNew[i] = ApplyPermToBitset(rm.M[i], perm, false)
	}

	return RMCode{r:rm.r, m:rm.m, M:mNew, diffs:rm.diffs, inBits:rm.inBits,
	              outBits:rm.outBits}	
}

func (rm RMCode) Print(showMatrix bool) {
	fmt.Printf("(%d, %d) | In bits = %d | Out bits = %d | ",
	           rm.r, rm.m, rm.inBits, rm.outBits)
	fmt.Printf("Expansion ratio = %f\n\n", float64(rm.outBits)/float64(rm.inBits))

	if showMatrix {
		for _, r := range rm.M {
			PrintBin(r, true)
		}
		fmt.Printf("\n")
	}

	return
}
