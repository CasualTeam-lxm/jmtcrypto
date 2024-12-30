// bitset.go
package mce

type Bitset = []bool

func ReverseBitset(bs Bitset) (bsNew Bitset) {
	bsNew = make(Bitset, len(bs))

	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
	    bsNew[i], bsNew[j] = bs[j], bs[i]
	}

	return
}

func InvertBitset(bs Bitset) (bsNew Bitset) {
	bsNew = make(Bitset, len(bs))

	for i := 0; i < len(bs); i++ {
	    bsNew[i] = !bs[i]
	}

	return	
}

func ParityOfBitset(bs Bitset) (par bool) {
	par = false
	for _, b := range bs {
		if b {
			par = !par
		}
	}

	return
}

func WeightOfBitset(bs Bitset) (wt int) {
	wt = 0
	for _, b := range bs {
		if b {
			wt += 1
		}
	}

	return
}

func BitsetXOR(b1 Bitset, b2 Bitset) (xor Bitset) {
	if len(b1) != len(b2) {
		//ERROR
	}

	xor = make(Bitset, len(b1))

	for i := 0; i < len(b1); i++ {
		xor[i] = b1[i] != b2[i]
	}

	return
}

func BitsetAND(b1 Bitset, b2 Bitset) (and Bitset) {
	and = make(Bitset, len(b1))

	for i := 0; i < len(b1); i++ {
		and[i] = (b1[i] && b2[i])
	}

	return
}

func BitsetOR(b1 Bitset, b2 Bitset) (or Bitset) {
	or = make(Bitset, len(b1))

	for i := 0; i < len(b1); i++ {
		or[i] = (b1[i] || b2[i])
	}

	return
}

func BitsetDOT(b1 Bitset, b2 Bitset) (dot bool) {
	dot = false
	for i := 0; i < len(b1); i++ {
		if (b1[i] && b2[i]) {
			dot = !dot
		}
	}

	return
}

func BitsetAllTrue(bs Bitset) (bool) {
	for _, b := range bs {
		if !b {
			return false
		}
	}

	return true 
}

func BitsetVote(bs Bitset, tie bool) bool {
	if tie {
		return 2*WeightOfBitset(bs) >= len(bs) 
	} else {
		return 2*WeightOfBitset(bs) > len(bs) 
	}
}

func BitsetFlipTopBit(b Bitset) Bitset {
	b[0] = !b[0]

	return b
}

func ApplyPermToBitset(bs Bitset, perm []int, forward bool) (pBs Bitset) {
	bitsPerPerm := len(perm)
	pBs = make(Bitset, len(bs))

	for blockId := 0; blockId < len(bs); blockId += bitsPerPerm {
		for i := 0; i < bitsPerPerm; i++ {
			if forward {
				pBs[blockId + perm[i]] = bs[blockId + i]
			} else {
				pBs[blockId + i] = bs[blockId + perm[i]]
			}
		}
	}

	return
}

func BitsetAllOnes(n int) (bs Bitset) {
	bs = make(Bitset, n)

	for i := 0; i < n; i++ {
		bs[i] = true
	}

	return
}
