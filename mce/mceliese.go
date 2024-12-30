// mceliese.go
package mce

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type PublicKey struct {
	RM	RMCode
}

type PrivateKey struct {
	RM			RMCode
	perm    []int
	C_inv   []Bitset
}

func GenerateKeyPair(r, m int) (PublicKey, PrivateKey) {
	privateRM := ReedMuller(r, m)

	privateRM.Print(false)

	perm := RandomPermutaion(privateRM.outBits)
	C, C_inv := MatrixPair(privateRM.inBits)

	publicRM := privateRM.PermuteCols(perm)
	publicRM.M = MatMulMat(C, publicRM.M)

	return PublicKey{RM:publicRM},
	       PrivateKey{RM:privateRM, perm:perm, C_inv:C_inv}
}

func (pubKey PublicKey) Encrypt(str string) Bitset {
	message := PadBlock(ParseText(str), pubKey.RM.inBits)

	return pubKey.RM.Encrypt(message, true)
}

func (privKey PrivateKey) Decrypt(cipherText Bitset) Bitset {
	// undo the permutation
	cipherText = ApplyPermToBitset(cipherText, privKey.perm, true)
	// apply standard rm decryption
	cipherText = privKey.RM.Decrypt(cipherText, true)
	// right multiply by C_inv
	postCipher := Bitset{}
	P := privKey.RM.inBits
	for i := 0; i < len(cipherText); i += P {
		eword := make(Bitset, P)
		copy(eword, cipherText[i:i+P])
		eword = MatMulVecR(eword, privKey.C_inv)

		postCipher = append(postCipher, eword...)
	}

	return postCipher
}

func (pubKey PublicKey) Write(filePath string) {
	fmt.Println("Saving public key...")
	f, _ := os.Create(filePath)

	defer f.Close()

	for _, row := range pubKey.RM.M {
		for _, bit := range row {
			if bit {
				f.WriteString("1")
			} else {
				f.WriteString("0")
			}
		}
		f.WriteString("\n")
	}

	f.Sync()
	return
}

func (privKey PrivateKey) Write(filePath string) {
	fmt.Println("Saving private key...")
	f, _ := os.Create(filePath)

	defer f.Close()

	// Save the coefficients of the Reed-Muller code
	f.WriteString(fmt.Sprintf("RM|%d|%d\n\n", privKey.RM.r, privKey.RM.m))

	for _, row := range privKey.C_inv {
		for _, bit := range row {
			if bit {
				f.WriteString("1")
			} else {
				f.WriteString("0")
			}
		}
		f.WriteString("\n")	
	}
	f.WriteString("\n")

	for _, i := range privKey.perm {
		f.WriteString(fmt.Sprintf("%d ", i))
	}
	f.WriteString("\n")

	f.Sync()
	return
}

func ReadPublic(fname string) (pubKey PublicKey) {
	b,_ := ioutil.ReadFile(fname)

	M := []Bitset{}

	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if len(l) == 0 { continue }

		row := make(Bitset, len(l))
		for j, char := range l {
			if char == rune('1') {
				row[j] = true
			}
		} 
		M = append(M, row)
	}

	pubKey.RM.M = M
	pubKey.RM.inBits = len(M)
	pubKey.RM.outBits = len(M[0])
	pubKey.RM.m = Log2(pubKey.RM.outBits)

	rows := 0
	for pubKey.RM.r = 0; pubKey.RM.r < pubKey.RM.m; pubKey.RM.r++ {
		rows += Choose(pubKey.RM.m, pubKey.RM.r)

		if (rows == pubKey.RM.inBits) {
			break
		}
	}

	return
}

func ReadPrivate(fname string) (privKey PrivateKey) {
	b,_ := ioutil.ReadFile(fname)

	mode := 0

	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if len(l) == 0 { 
			mode++
			continue
		}

		if (mode == 0) {
			parts := strings.Split(l, "|")
			// BlockCiphetype := parts[0] // should always be "RM"
			r, _ := strconv.Atoi(parts[1])
			m, _ := strconv.Atoi(parts[2])

			privKey.RM = ReedMuller(r,m)
		}

		if (mode == 1) {
			row := make(Bitset, len(l))
			for j, char := range l {
				if char == rune('1') {
					row[j] = true
				}
			} 
			privKey.C_inv = append(privKey.C_inv, row)
		}

		if (mode == 2) {
			parts := strings.Split(l, " ")
			for _, p := range parts {
				if p == "" {
					continue
				}

				n, _ := strconv.Atoi(p)
				privKey.perm = append(privKey.perm, n)
			}
		}

	}

	return
}
