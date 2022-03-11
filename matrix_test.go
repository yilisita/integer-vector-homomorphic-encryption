/*
 * @Author: Wen Jiajun
 * @Date: 2022-02-27 19:46:25
 * @LastEditors: Wen Jiajun
 * @LastEditTime: 2022-03-11 17:55:42
 * @FilePath: \integer-vector-homomorphic-encryption\matrix_test.go
 * @Description:
 */

package intvec

import (
	"math/big"
	"reflect"
	"testing"
)

func TestMatrix_Marshal(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "1",
			m: Matrix{
				rows: 1, columns: 1, data: []*big.Int{big.NewInt(1989)},
			},
			want: []byte{49, 49, 49, 57, 56, 57},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Marshal(); !reflect.DeepEqual(got, tt.want) {
				t.Logf("%v", string(got))
				t.Logf("%v", string(tt.want))
				t.Errorf("Matrix.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Unmarshal(t *testing.T) {
	type args struct {
		matrixByte []byte
	}
	tests := []struct {
		name    string
		m       *Matrix
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			m:    &Matrix{},
			args: args{
				matrixByte: []byte(`{"rows":1,"cols":1,"datastr":["1989"]}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Unmarshal(tt.args.matrixByte); (err != nil) != tt.wantErr {
				t.Errorf("Matrix.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
}
