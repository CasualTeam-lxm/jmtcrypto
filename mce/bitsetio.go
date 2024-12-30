// bitsetio.go
package mce

import (
	"fmt"
	// "errors"
	// "strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
//
// Convert to and from Bitsets
//
func CharToBitset(c rune) (bs Bitset) {
	bs = make(Bitset, 8)

	for b := 0; b < 8; b++ {
		bs[b] = ((c >> (7 - b)) & 1) == 1
	}

	return
}

func ParseText(s string) Bitset {
	msg := make(Bitset, 0)
	for _, char := range s {
		binChar := CharToBitset(char)
		msg = append(msg, binChar...)
	}

	return msg
}

func DeparseMessage(bs Bitset) string {
	var sb strings.Builder
	byte := uint8(0)
	for i, b := range bs {
		byte <<= 1
		if b {
			byte += 1
		}

		if i % 8 == 7 {
			sb.WriteString(string(rune(byte)))
			byte = 0
		}
	}
	return sb.String()
}

func PrintHex(b Bitset, newLine bool) {
	plain := DeparseMessage(b)
	for _, char := range plain {
		fmt.Printf("%02x ", char)
	}

	if newLine {
		fmt.Println("")
	}
}

func PadBlock(b Bitset, n int) Bitset {
	need := (len(b) % n)

	if need == 0 {
		return b
	}

	for i := 0; i < n - need; i++ {
		b = append(b, true)
	}

	return b
}

func PrintBin(bs Bitset, newLine bool) {
	for i, b := range bs {
		if b {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		if i % 8 == 7 {
			fmt.Print(" ")
		}
	}
	if newLine {
		fmt.Println("")
	}
}

func PrintASCII(b Bitset, newLine bool) {
	fmt.Print(DeparseMessage(b))
	if newLine {
		fmt.Println("")
	}
}
