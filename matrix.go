/*
 * @Author: Wen Jiajun
 * @Date: 2022-02-27 19:46:25
 * @LastEditors: Wen Jiajun
 * @LastEditTime: 2022-02-28 18:16:20
 * @FilePath: \test\intvec\matrix.go
 * @Description:
 */

package intvec

// Matrix
/*
------------------------------------------------------
作者   : Black Ghost
日期   : 2018-11-20
版本   : 0.0.0
         0.0.1 2018-12-11 增加切片与矩阵转换
         0.0.2 2018-12-26 增加错误报告
         0.0.3 2018-12-27 增加追加行/列
------------------------------------------------------
    矩阵的创建及其操作创建及其简单操作/运算
理论：
    参考 OneThin // http://outofmemery.cn/code-snippet
         /16991/go-language-matrix-operation
    进行了主要运算和结构的补充与修改
------------------------------------------------------
注意事项：
    1. r, c 是从零开始算的
------------------------------------------------------
*/

import (
	"fmt"
	"math/big"
)

//数据结构定义----------------------------------------+
// Matrix 定义Matrix数据类型
type Matrix struct {
	rows    int
	columns int
	data    []*big.Int //将矩阵中所有元素作为一维切片
}

func (m *Matrix) GetRows() int {
	return m.rows
}

func (m *Matrix) GetColumns() int {
	return m.columns
}

//矩阵操作-------------------------------------------+
//通过行列号寻找指定矩阵位置在一维切片中的编号
func (m *Matrix) findIndex(r, c int) int {
	//r E [0, n), c E [0, n)
	return r*m.GetColumns() + c
}

// SetMatrix 设置指定行列的值
func (A *Matrix) SetMatrix(r, c int, val *big.Int) {
	if (r >= A.GetRows()) || (c >= A.GetColumns()) {
		panic("Error in goNum.(*Matrix).SetMatrix: Out of range")
	}
	A.data[A.findIndex(r, c)] = val
}

// GetFromMatrix 获取指定行列的值
func (A *Matrix) GetFromMatrix(r, c int) *big.Int {
	if (r >= A.GetRows()) || (c >= A.GetColumns()) {
		panic("Error in goNum.(*Matrix).SetMatrix: Out of range")
	}
	return A.data[A.findIndex(r, c)]
}

// RowOfMatrix 获取指定行的值的切片
func (A *Matrix) RowOfMatrix(i int) []*big.Int {
	if i >= A.GetRows() {
		panic("Error in goNum.(*Matrix).RowOfMatrix: Out of range")
	}
	return A.data[A.findIndex(i, 0):A.findIndex(i, A.GetColumns())]
}

// ColumnOfMatrix 获取指定列的值的切片
func (A *Matrix) ColumnOfMatrix(j int) []*big.Int {
	if j >= A.GetColumns() {
		panic("Error in goNum.(*Matrix).ColumnOfMatrix: Out of range")
	}
	col := make([]*big.Int, A.GetRows())
	for i := 0; i < A.GetRows(); i++ {
		col[i] = A.RowOfMatrix(i)[j]
	}
	return col
}

// Transpose 矩阵转置
func (A *Matrix) Transpose() Matrix {
	B := ZeroMatrix(A.GetColumns(), A.GetRows())
	for i := 0; i < A.GetRows(); i++ {
		for j := 0; j < A.GetColumns(); j++ {
			B.SetMatrix(j, i, A.GetFromMatrix(i, j))
		}
	}
	return B
}

// AppendRow 追加一行,另外一种方法是追加数据到A.data，测试显示其速度表现更差
func (A *Matrix) AppendRow(row []*big.Int) Matrix {
	//判断row长度是否等于A列数
	if len(row) != A.GetColumns() {
		panic("Error in goNum.(*Matrix).AppendRow: Slice length error")
	}
	B := ZeroMatrix(A.GetRows()+1, A.GetColumns())
	n := A.GetRows() * A.GetColumns()
	for i := 0; i < n; i++ {
		B.data[i] = A.data[i]
	}
	for i := 0; i < len(row); i++ {
		B.data[n+i] = row[i]
	}
	return B
}

// AppendColumn 追加一列，对于多次调用，建议组合使用转置和追加行
func (A *Matrix) AppendColumn(col []*big.Int) Matrix {
	//判断row长度是否等于A列数
	if len(col) != A.GetRows() {
		panic("Error in goNum.(*Matrix).AppendColumn: Slice length error")
	}
	B := ZeroMatrix(A.GetRows(), A.GetColumns()+1)
	for i := 0; i < A.GetRows(); i++ {
		for j := 0; j < A.GetColumns(); j++ {
			B.SetMatrix(i, j, A.GetFromMatrix(i, j))
		}
		B.SetMatrix(i, A.GetColumns(), col[i])
	}
	return B
}

//矩阵初始化-----------------------------------------+
// ZeroMatrix r行c列零矩阵
func ZeroMatrix(r, c int) Matrix {
	var data []*big.Int
	for i := 0; i < r*c; i++ {
		data = append(data, big.NewInt(0))
	}
	return Matrix{r, c, data}
}

// IdentityE n阶单位矩阵
func IdentityE(n int) Matrix {
	A := ZeroMatrix(n, n)
	for i := 0; i < len(A.data); i += (n + 1) {
		A.data[i] = big.NewInt(1)
	}
	return A
}

// NewMatrix 以已有数据创建r行c列矩阵
func NewMatrix(r, c int, data []*big.Int) Matrix {
	// make a copy of the original message (deep copy)
	datacp := append([]*big.Int{}, data...)
	if len(data) != r*c {
		panic("goNum.Matrix.New: Length of data does not matched r rows and c columns")
	}
	A := ZeroMatrix(r, c)
	A.data = datacp
	return A
}

// Slices1ToMatrix 一维切片转为矩阵(列向量)
func Slices1ToMatrix(s []*big.Int) Matrix {
	A := ZeroMatrix(len(s), 1)
	for i := 0; i < A.GetRows(); i++ {
		A.data[i] = s[i]
	}
	return A
}

// Slices2ToMatrix 二维切片转为矩阵
func Slices2ToMatrix(s [][]*big.Int) Matrix {
	row := len(s)
	col := len(s[0])
	A := ZeroMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			A.SetMatrix(i, j, s[i][j])
		}
	}
	return A
}

// Matrix1ToSlices 列向量转为一维切片
func Matrix1ToSlices(A Matrix) []*big.Int {
	s := make([]*big.Int, A.GetRows())
	for i := 0; i < A.GetRows(); i++ {
		s[i] = A.data[i]
	}
	return s
}

// ExtendHorizontal extends the matrix horizontally
// e.g: A.ExtendHorizontal(B) => [A, B]
func ExtendHorizontal(m, n Matrix) (Matrix, error) {
	if m.rows != n.rows {
		return Matrix{}, fmt.Errorf("Unequal Rows: Left rows:%v, right rows%v", m.rows, n.rows)
	}

	var resData []*big.Int
	for i := 0; i < m.rows; i++ {
		resData = append(resData, m.RowOfMatrix(i)...)
		resData = append(resData, n.RowOfMatrix(i)...)
	}

	return NewMatrix(m.rows, m.columns+n.columns, resData), nil
}

// Matrix2ToSlices 二维矩阵转为二维切片
func Matrix2ToSlices(A Matrix) [][]*big.Int {
	s := make([][]*big.Int, A.GetRows())
	for i := 0; i < A.GetRows(); i++ {
		s[i] = make([]*big.Int, A.GetColumns())
		for j := 0; j < A.GetColumns(); j++ {
			s[i][j] = A.GetFromMatrix(i, j)
		}
	}
	return s
}

//矩阵运算------------------------------------------+
// AddMatrix 矩阵相加
func AddMatrix(A, B Matrix) Matrix {
	if (A.GetRows() != B.GetRows()) || (A.GetColumns() != B.GetColumns()) {
		panic("goNum.Matrix.Add: A and B does not matched")
	}
	AaddB := ZeroMatrix(A.GetRows(), A.GetColumns())
	for i := 0; i < A.GetRows(); i++ {
		for j := 0; j < A.GetColumns(); j++ {
			AaddB.SetMatrix(i, j, new(big.Int).Add(A.GetFromMatrix(i, j), B.GetFromMatrix(i, j)))
		}
	}
	return AaddB
}

// SubMatrix 矩阵相减
func SubMatrix(A, B Matrix) Matrix {
	if (A.GetRows() != B.GetRows()) || (A.GetColumns() != B.GetColumns()) {
		panic("goNum.Matrix.Sub: A and B does not matched")
	}
	AsubB := ZeroMatrix(A.GetRows(), A.GetColumns())
	for i := 0; i < A.GetRows(); i++ {
		for j := 0; j < A.GetColumns(); j++ {
			AsubB.SetMatrix(i, j, new(big.Int).Sub(A.GetFromMatrix(i, j), B.GetFromMatrix(i, j)))
		}
	}
	return AsubB
}

// NumProductMatrix 矩阵数乘
func NumProductMatrix(A Matrix, c *big.Int) Matrix {
	cA := ZeroMatrix(A.GetRows(), A.GetColumns())
	for i := 0; i < len(cA.data); i++ {
		cA.data[i] = new(big.Int).Mul(c, A.data[i])
	}
	return cA
}

// DotPruduct 矩阵点乘
func DotPruduct(A, B Matrix) Matrix {
	if A.GetColumns() != B.GetRows() {
		panic("goNum.Matrix.DotPruduct: A and B does not matched")
	}
	AdotB := ZeroMatrix(A.GetRows(), B.GetColumns())
	for i := 0; i < A.GetRows(); i++ {
		for j := 0; j < B.GetColumns(); j++ {
			for k := 0; k < A.GetColumns(); k++ {
				temp := new(big.Int).Mul(A.GetFromMatrix(i, k), B.GetFromMatrix(k, j))
				AdotB.data[B.GetColumns()*i+j] = big.NewInt(0).Add(AdotB.data[B.GetColumns()*i+j], temp)
			}
		}
	}
	return AdotB
}

// // CrossVector 向量叉乘，得到垂直于两个向量所在平面的向量
// func CrossVector(a, b []*big.Int) []*big.Int {
// 	if (len(a) != 3) || (len(b) != 3) {
// 		panic("goNum.Matrix.CrossVector: vector a or b length is not 3")
// 	}
// 	acrossb := make([]*big.Int, 3)
// 	acrossb[0] = a[1]*b[2] - a[2]*b[1]
// 	acrossb[1] = a[2]*b[0] - a[0]*b[2]
// 	acrossb[2] = a[0]*b[1] - a[1]*b[0]
// 	return acrossb
// }
