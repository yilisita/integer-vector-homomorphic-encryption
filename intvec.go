/*
 * @Author: Wen Jiajun
 * @Date: 2022-01-29 15:03:03
 * @LastEditors: Wen Jiajun
 * @LastEditTime: 2022-02-28 18:50:49
 * @FilePath: \test\intvec\intvec.go
 * @Description: an implementation for integer vector encryption schema
 *               see(https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.471.387&rep=rep1&type=pdf)
 */

package intvec

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

const (
	//      = 10
	l      = 100
	aBound = 100
	bBound = 1000
	tBound = 100
)

var (
	w       = math.Pow(2, 45)
	wStr    = strconv.FormatFloat(w, 'f', 0, 64)
	wBig, _ = new(big.Int).SetString(wStr, 10)
)

// PublicKey: the switching matrix
type PublicKey struct {
	Matrix
}

// PrivateKey: the private key matrix
type PrivateKey struct {
	Matrix
}

// a vector
type Plaintext struct {
	Matrix
}

// a vector, usually longer than the plaintext vector
type Ciphertext struct {
	Matrix
}

// NewPlaintext returns a plaintext using the input,
// By default, the row = len(data), the column = 1
func NewPlaintext(data []*big.Int) Plaintext {
	return Plaintext{
		NewMatrix(len(data), 1, data),
	}
}

// NewCiphertext returns a ciphertext using the input,
// By default, the row = len(data), the column = 1
func NewCiphertext(data []*big.Int) Ciphertext {
	return Ciphertext{
		NewMatrix(len(data), 1, data),
	}
}

// GetKeyPairs generate the (privKey, pubKey) pair
// @param: row = the message vector's length
// @param: col = the ciphertext vector's length(self-defined)
// @return: (privKey, pubKey)
func GetKeyPairs(row, col, bound int) (PrivateKey, PublicKey) {
	// Randomly construct the T part of a new private key
	// For production security, use "crypto/rand"
	rand.Seed(time.Now().Unix())
	var T = ZeroMatrix(row, col)
	//var data = []*big.Int
	var data = make([]*big.Int, row*col)
	for i := 0; i < row*col; i++ {
		data[i] = big.NewInt(int64(rand.Intn(bound)))
	}
	T.data = data

	// concatenate an identity matrix I with the T matrix,
	// Thus, we get a whole private key: [I, T]
	I := IdentityE(T.GetRows())
	var (
		ISlice = make([][]*big.Int, 0)
		TSlice = make([][]*big.Int, 0)
	)
	ISlice = Matrix2ToSlices(I)
	TSlice = Matrix2ToSlices(T)
	for i := 0; i < T.GetRows(); i++ {
		ISlice[i] = append(ISlice[i], TSlice[i]...)
	}
	var res = make([]*big.Int, 0)
	for _, s := range ISlice {
		for _, j := range s {
			res = append(res, j)
		}
	}
	PrivKey := NewMatrix(I.GetRows(), I.GetColumns()+T.GetColumns(), res)

	// Now, we construct the public key
	// the public key is the key-switching matrix,
	// which swtich the old private key I
	// to the above new private key [I, T]
	pubKey := KeySwitchMatrix(PrivateKey{I}, T)
	return PrivateKey{PrivKey}, pubKey

}

// TODO: when doing encryption, we should use a public key, in the paper, that is M
// HOW to fix this ?
func Encrypt(pubk PublicKey, x Plaintext) Ciphertext {
	// I := IdentityE(x.GetRows())
	// var xSub = make([]*big.Int, 0)
	// xSub = append(xSub, x.data...)
	// var xS = NewMatrix(x.GetRows(), x.GetColumns(), xSub)

	for i := 0; i < len(x.data); i++ {
		x.data[i] = new(big.Int).Mul(wBig, x.data[i])
	}
	return SwitchCipher(pubk, Ciphertext{x.Matrix})
}

func Decrypt(prvk PrivateKey, c Ciphertext) Plaintext {
	sc := DotPruduct(prvk.Matrix, c.Matrix)
	x := make([]*big.Int, 0)
	var temp float64
	for _, i := range sc.data {
		iFloat, _ := strconv.ParseFloat(i.String(), 64)
		temp = (nearestInteger(iFloat))
		if temp < 0 {
			temp = temp - 1
		}
		tempStr := strconv.FormatFloat(temp, 'f', 0, 64)
		tempRes, _ := new(big.Int).SetString(tempStr, 10)
		x = append(x, tempRes)
	}
	return Plaintext{NewMatrix(sc.GetRows(), 1, x)}
}

// SwitchCipher returns a new ciphertext using the public key M to
// multiply the bit vector of the ciphertext
func SwitchCipher(M PublicKey, c Ciphertext) Ciphertext {
	cstar := getBitVector(c.Matrix)
	return Ciphertext{DotPruduct(M.Matrix, cstar)}
}

// AddCiphertext returns the sum of the 2 input ciphertext
func AddCiphertext(c1, c2 Ciphertext) Ciphertext {
	return Ciphertext{AddMatrix(c1.Matrix, c2.Matrix)}
}

// S: the original private key
// T: the "T" part of the new private key
// return M:
//			M.GetRows() = S.GetRows() + T.GetColumns()
//  		M.GetColumns() = S.GetColumns()
// Therefor, we could only manually adjust T.GetColumns() to reshape new ciphertext's row
// Thus we could reduce the ciphertext's dimensions
func KeySwitchMatrix(S PrivateKey, T Matrix) PublicKey {
	sStar := getBitMatrix(S.Matrix)
	A := getRandomMatrix(T.GetColumns(), sStar.GetColumns(), aBound)
	E := getRandomMatrix(sStar.GetRows(), sStar.GetColumns(), bBound)
	up1 := AddMatrix(sStar, E)
	up2 := SubMatrix(up1, DotPruduct(T, A))
	ASLice := Matrix2ToSlices(A)
	USlice := Matrix2ToSlices(up2)
	for _, j := range ASLice {
		USlice = append(USlice, j)
	}
	var res = make([]*big.Int, 0)
	for _, s := range USlice {
		for _, j := range s {
			res = append(res, j)
		}
	}
	return PublicKey{NewMatrix(E.GetRows()+A.GetRows(), A.GetColumns(), res)}
}

// THIS IS A SERVER SIDE METHOD:
// GetInnerProduct receives a key switching matrix M, two encrypted ciphertexts, i.e. c1 and c2,
// which are equal in length and width, and returns the result of their inner product, which
// should be decrypted with a new private key constructed by the ciphertexts' private keys
// (this means that the input ciphertexts may have different or the same keys).
func GetInnerProduct(c1, c2 Ciphertext, M PublicKey) Ciphertext {
	if c1.GetRows() != c2.GetRows() || c1.GetColumns() != c2.GetColumns() {
		panic("Unmatched shape of c1 and c2: c1.GetRows()/c1.GetColumns() should be equal c2.GetRows()/Columns.")
	}

	if M.GetColumns() != c1.GetRows()*c2.GetRows()*l {
		panic("Cannot use M to reduce the dimension of vec(c1 * c2'): M.GetColumns() and vec(c1 * c2').GetRows() unmatched.")
	}
	// calculate the new ciphertext
	// c1c2T = c1 * c2'
	var c1c2T = DotPruduct(c1.Matrix, c2.Transpose())

	// construct vec(c1c2T)
	var flattenc1c2T []*big.Int
	for i := 0; i < c1c2T.GetRows(); i++ {
		flattenc1c2T = append(flattenc1c2T, c1c2T.RowOfMatrix(i)...)
	}

	// calculate vec(c1c2T) / w
	for k, v := range flattenc1c2T {

		flattenc1c2T[k] = new(big.Int).Div(v, wBig) // original: v / w
	}
	var flattenc1c2TMatrix = NewCiphertext(flattenc1c2T)
	// do the dimension reduction
	return SwitchCipher(M, flattenc1c2TMatrix)
}

// THIS IS A CLIENT SIDE METHOD:
// GetInnerProductKey compute a new temporary private key which is
// very long and cannot be used directly to decrypt the new
// ciphertext returned by GetInnerProduct.
func GetInnerProductLongKey(s1, s2 PrivateKey) PrivateKey {
	if s1.GetRows() != s2.GetRows() || s1.GetColumns() != s2.GetColumns() {
		panic("Unmatched shape of s1 and s2: s1.GetRows()/s1.GetColumns() should be equal s2.GetRows()/s2.GetColumns().")
	}

	var s1Ts2 = DotPruduct(s1.Transpose(), s2.Matrix)
	var s1Ts2Flatten []*big.Int
	for i := 0; i < s1Ts2.GetRows(); i++ {
		s1Ts2Flatten = append(s1Ts2Flatten, s1Ts2.RowOfMatrix(i)...)
	}
	return PrivateKey{NewMatrix(1, s1Ts2.GetRows()*s1Ts2.GetRows(), s1Ts2Flatten)}
}

// THIS IS A CLINET SIDE METHOD:
// Refer to KeySwitchMatrix, we could only adjust the columns of T,
// i.e. the "T" part of the new private key which could be generated
// with method "GetRandomKey" to shape the new ciphertext's GetRows() = S.GetRows() + T.GetColumns().
// Because S.GetRows() = 1 (under the circumstance of Inner Product), so we could reshape ciphertext's dimension to n,
// by adjusting T.GetColumns() = n - 1,
//				i.e. var T = GetRandomMatrix(1, n-1, bound)
//											 |
// ps: recall method "GetSecretKey", we compute a key, s, whose GetRows() equal T.GetRows()
// @Param: s = vec(s1Ts2), that is the key generated by GetInnerProductLongKey
// @Return: M: the key switching matrix to reduce ciphertext's dimension
//			S: the final private key to decrypt the dimension-reduced ciphertext
func GetInnerProductKeyPairs(s PrivateKey) (PrivateKey, PublicKey) {
	var n = int(math.Sqrt(float64(s.GetColumns())))
	var T = getRandomMatrix(1, n-1, tBound)
	var S = getSecretKey(T)
	// the final new key must be returned, otherwise we have no way to obtain it
	return S, KeySwitchMatrix(s, T)
}

// convert a decimal number to a binary one, with the length of 100.
// Empty place will have "0" as a replacement
// e.g.
// convertToBin(big.NewInt(2))
// output: "0000000...0010" (length = 100)
func convertToBin(num *big.Int) string {
	s := ""
	var negative = false
	if num.Sign() < 0 {
		negative = true
		num = num.Abs(num)
	}

	//
	s = fmt.Sprintf("%b", num)
	s = fillStrLengthToL(s, l)
	if negative {
		return "-" + s
	}
	return s
}

// fillStrLengthToL returns a string with the length of L
// "0" will be put to the unoccupied place on the left side
// e.g.
// fillStrLengthToL("my", 5) => "000my"
func fillStrLengthToL(s string, L int) string {
	var res = s
	if len(s) < L {
		var delt = L - len(s)
		for i := 0; i < delt; i++ {
			res = "0" + res
		}
	}
	return res
}

// reverse returns the reverse of the input string.
// e.g.
// reverse("hello") => "olleh"
func reverse(str string) string {
	rs := []rune(str)
	len := len(rs)
	var tt []rune

	tt = make([]rune, 0)
	for i := 0; i < len; i++ {
		tt = append(tt, rs[len-i-1])
	}
	return string(tt[0:])
}

// getRandomMatrix returns a matrix shaped by the parameters (row, col)
// the returned matrix's element is bounded within the bound
// e.g.
// getRandomMatrix(1, 3, 100)
// output might be: [34, 56, 79]
func getRandomMatrix(row, col, bound int) Matrix {
	rand.Seed(time.Now().Unix())
	var T = ZeroMatrix(row, col)
	//var data = []*big.Int
	var data = make([]*big.Int, row*col)
	for i := 0; i < row*col; i++ {
		data[i] = big.NewInt(int64(rand.Intn(bound)))
	}
	T.data = data
	return T
}

// getBitVector returns a bit representation vector, whose element is
// that of the input multiplied l.
// e.g:
// input: [1, 2, 3]
// output:[1, 0, 0...0,
//		   0, 1, 0...0,
//         1, 1, 0...0,]
// length of output: 3 * l = 300
func getBitVector(c Matrix) Matrix {
	var (
		res        = make([]*big.Int, 0)
		sign int64 = 1
		s    string
	)
	for _, i := range c.data {
		s = convertToBin(i)
		if s[0] == '-' {
			sign = -1
			s = "0" + s[1:]
		}
		s = reverse(s)
		for _, j := range s {
			res = append(res, big.NewInt(int64(j-48)*sign))
		}
	}
	A := NewMatrix(l*c.GetRows()*c.GetColumns(), 1, res)
	return A
}

// getBitMatrix return a bit representation of a matrix
// see "getBitVector" for more information
func getBitMatrix(s Matrix) Matrix {
	var powers = make([]string, l)
	for i := 0; i < l; i++ {
		powers[i] = strconv.FormatFloat(math.Pow(2, float64(i)), 'f', 0, 64)
	}
	var res = make([]*big.Int, 0)
	for _, k := range s.data {
		for _, j := range powers {
			temp, _ := new(big.Int).SetString(j, 10)
			res = append(res, new(big.Int).Mul(k, temp))
		}
	}
	var final = NewMatrix(s.GetRows(), s.GetColumns()*l, res)
	return final
}

// getSecretKey returns a privateKey,it simply concatenates an
// identity matrix with the input matrix, thus [I, T], the
// PrivateKey is returned
func getSecretKey(T Matrix) PrivateKey {
	I := IdentityE(T.GetRows())
	var (
		ISlice = make([][]*big.Int, 0)
		TSlice = make([][]*big.Int, 0)
	)
	ISlice = Matrix2ToSlices(I)
	TSlice = Matrix2ToSlices(T)
	for i := 0; i < T.GetRows(); i++ {
		ISlice[i] = append(ISlice[i], TSlice[i]...)
	}
	var res = make([]*big.Int, 0)
	for _, s := range ISlice {
		for _, j := range s {
			res = append(res, j)
		}
	}
	A := NewMatrix(I.GetRows(), I.GetColumns()+T.GetColumns(), res)
	return PrivateKey{A}
}

// nearestInteger returns the nearest interger of the input divided by w
func nearestInteger(x float64) float64 {
	return math.Floor((x + (w+1)/2) / w)
}
