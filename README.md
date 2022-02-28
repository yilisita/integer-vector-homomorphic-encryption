<!--
 * @Author: Wen Jiajun
 * @Date: 2022-02-05 15:18:18
 * @LastEditors: Wen Jiajun
 * @LastEditTime: 2022-02-28 18:52:37
 * @FilePath: \test\intvec\README.md
 * @Description: 
-->

# integer vector homomorphic encryption
This is an implementation of the [integer vector homomorphic encryption](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.471.387&rep=rep1&type=pdf) in Go.

Full tests of the functions have not been conducted (because I don't know how) so there might be some bugs in it.However, most of the functions should be workable.

## To Use
* Get the package:
```
    go get -u "github.com/yilisita/intvec"
```
* Generate the plaintext:
```
    var x1 = []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
    var m1 = intvec.NewPlaintext(x1)
    var x2 = []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
    var m2 = intvec.NewPlaintext(x2)
```
* Generate the Key pairs:
```
    # This is the key generation for one-elememt vector
    # I'm not sure its appropriate, but it works ok
    var sk, pk = intvec.GetKeyPairs(3, 10, keyBound)
```
* Encryption:
```
    var c1 = intvec.Encrypt(pk, m1)
    var c2 = intvec.Encrypt(pk, m2)
``` 
* Homomorphic Calculation:
```
    var cRes = intvec.AddCiphertext(c1, c2)
```
* Decryption:
```
    var decryptedRes = intvec.Decrypt(sk, cRes)
```

# Reference
1. [Efficient homomorphic encryption on integer vectors and its applications](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.471.387&rep=rep1&type=pdf)
2. Matrix implementation is based on [goNum](https://github.com/chfenger/goNum) with several modifications

# PS
Please forgive my poor English, it's not my first language. I'm not fully understanding the paper, so there might be some errors in this package
and some important featurs might be missing.