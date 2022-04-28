/*
 * @Author: Wen Jiajun
 * @Date: 2022-01-29 15:03:03
 * @LastEditors: Wen Jiajun
 * @LastEditTime: 2022-04-28 14:50:51
 * @FilePath: \integer-vector-homomorphic-encryption\intvec_test.go
 * @Description: an implementation for integer vector encryption schema
 *               see(https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.471.387&rep=rep1&type=pdf)
 */

package intvec

import (
	"math/big"
	"reflect"
	"testing"
)

func TestPublicKey_Convert2Byte(t *testing.T) {
	tests := []struct {
		name string
		pk   PublicKey
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "1",
			pk: PublicKey{
				NewMatrix(1, 1, []*big.Int{
					big.NewInt(1989),
				})},
			want: []byte(`{"rows":1,"cols":1,"datastr":["1989"]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Convert2Byte(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKey.Convert2Byte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrivateKeyFromByte(t *testing.T) {
	type args struct {
		skByte []byte
	}
	tests := []struct {
		name    string
		args    args
		want    PrivateKey
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				skByte: []byte(`{"rows":1,"cols":1,"datastr":["1989"]}`),
			},
			want: PrivateKey{
				NewMatrix(1, 1, []*big.Int{big.NewInt(1989)}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrivateKeyFromByte(tt.args.skByte)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrivateKeyFromByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrivateKeyFromByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_GetT(t *testing.T) {
	sk, _ := GetKeyPairs(10, 10, 10000)
	tests := []struct {
		name string
		sk   PrivateKey
		want *Matrix
	}{
		// TODO: Add test cases.
		{
			"1",
			*sk,
			&Matrix{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sk.GetT(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKey.GetT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInnerProduct(t *testing.T) {
	var d1 []*big.Int
	var d2 []*big.Int
	var length = 10
	var sum int = 0
	for i := 0; i < length; i++ {
		sum += i
		d1 = append(d1, big.NewInt(int64(i)))
		d2 = append(d2, big.NewInt(2))
	}

	var sk, pk = GetKeyPairs(length, length, 1000)

	c1 := Encrypt(pk, NewPlaintext(d1))
	c2 := Encrypt(pk, NewPlaintext(d2))

	var sknew, _ = GetKeyPairs(length, 1, 1000)
	var t_part = sknew.GetT()
	var m = KeySwitchMatrix(sk, t_part)

	res_cipher := GetInnerProduct(c1, c2, m)
	sknewlong := GetInnerProductLongKey(sknew, sknew)
	dec_res_cipher := Decrypt(sknewlong, res_cipher)

	res_data := dec_res_cipher.GetData()
	res := res_data[0].Int64()

	if 2*int64(sum) != res {
		t.Errorf("Expected %v, Got %v instead.", sum, res)
	}

}
