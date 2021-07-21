package sm4

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"math/bits"
	"strconv"
)

const (
	BlockSize = 16
	KeySize   = 16
)

var sBox = [256]byte{
	0xd6, 0x90, 0xe9, 0xfe, 0xcc, 0xe1, 0x3d, 0xb7,
	0x16, 0xb6, 0x14, 0xc2, 0x28, 0xfb, 0x2c, 0x05,
	0x2b, 0x67, 0x9a, 0x76, 0x2a, 0xbe, 0x04, 0xc3,
	0xaa, 0x44, 0x13, 0x26, 0x49, 0x86, 0x06, 0x99,
	0x9c, 0x42, 0x50, 0xf4, 0x91, 0xef, 0x98, 0x7a,
	0x33, 0x54, 0x0b, 0x43, 0xed, 0xcf, 0xac, 0x62,
	0xe4, 0xb3, 0x1c, 0xa9, 0xc9, 0x08, 0xe8, 0x95,
	0x80, 0xdf, 0x94, 0xfa, 0x75, 0x8f, 0x3f, 0xa6,
	0x47, 0x07, 0xa7, 0xfc, 0xf3, 0x73, 0x17, 0xba,
	0x83, 0x59, 0x3c, 0x19, 0xe6, 0x85, 0x4f, 0xa8,
	0x68, 0x6b, 0x81, 0xb2, 0x71, 0x64, 0xda, 0x8b,
	0xf8, 0xeb, 0x0f, 0x4b, 0x70, 0x56, 0x9d, 0x35,
	0x1e, 0x24, 0x0e, 0x5e, 0x63, 0x58, 0xd1, 0xa2,
	0x25, 0x22, 0x7c, 0x3b, 0x01, 0x21, 0x78, 0x87,
	0xd4, 0x00, 0x46, 0x57, 0x9f, 0xd3, 0x27, 0x52,
	0x4c, 0x36, 0x02, 0xe7, 0xa0, 0xc4, 0xc8, 0x9e,
	0xea, 0xbf, 0x8a, 0xd2, 0x40, 0xc7, 0x38, 0xb5,
	0xa3, 0xf7, 0xf2, 0xce, 0xf9, 0x61, 0x15, 0xa1,
	0xe0, 0xae, 0x5d, 0xa4, 0x9b, 0x34, 0x1a, 0x55,
	0xad, 0x93, 0x32, 0x30, 0xf5, 0x8c, 0xb1, 0xe3,
	0x1d, 0xf6, 0xe2, 0x2e, 0x82, 0x66, 0xca, 0x60,
	0xc0, 0x29, 0x23, 0xab, 0x0d, 0x53, 0x4e, 0x6f,
	0xd5, 0xdb, 0x37, 0x45, 0xde, 0xfd, 0x8e, 0x2f,
	0x03, 0xff, 0x6a, 0x72, 0x6d, 0x6c, 0x5b, 0x51,
	0x8d, 0x1b, 0xaf, 0x92, 0xbb, 0xdd, 0xbc, 0x7f,
	0x11, 0xd9, 0x5c, 0x41, 0x1f, 0x10, 0x5a, 0xd8,
	0x0a, 0xc1, 0x31, 0x88, 0xa5, 0xcd, 0x7b, 0xbd,
	0x2d, 0x74, 0xd0, 0x12, 0xb8, 0xe5, 0xb4, 0xb0,
	0x89, 0x69, 0x97, 0x4a, 0x0c, 0x96, 0x77, 0x7e,
	0x65, 0xb9, 0xf1, 0x09, 0xc5, 0x6e, 0xc6, 0x84,
	0x18, 0xf0, 0x7d, 0xec, 0x3a, 0xdc, 0x4d, 0x20,
	0x79, 0xee, 0x5f, 0x3e, 0xd7, 0xcb, 0x39, 0x48,
}

var cK = [32]uint32{
	0x00070e15, 0x1c232a31, 0x383f464d, 0x545b6269,
	0x70777e85, 0x8c939aa1, 0xa8afb6bd, 0xc4cbd2d9,
	0xe0e7eef5, 0xfc030a11, 0x181f262d, 0x343b4249,
	0x50575e65, 0x6c737a81, 0x888f969d, 0xa4abb2b9,
	0xc0c7ced5, 0xdce3eaf1, 0xf8ff060d, 0x141b2229,
	0x30373e45, 0x4c535a61, 0x686f767d, 0x848b9299,
	0xa0a7aeb5, 0xbcc3cad1, 0xd8dfe6ed, 0xf4fb0209,
	0x10171e25, 0x2c333a41, 0x484f565d, 0x646b7279,
}

var fK = [4]uint32{
	0xa3b1bac6, 0x56aa3350, 0x677d9197, 0xb27022dc,
}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "sm4: invalid key size " + strconv.Itoa(int(k))
}

type sm4Cipher struct {
	enc []uint32
	dec []uint32
}

func NewCipher(key []byte) (cipher.Block, error) {
	n := len(key)
	if n != KeySize {
		return nil, KeySizeError(n)
	}
	c := new(sm4Cipher)
	c.enc = expandKey(key, true)
	c.dec = expandKey(key, false)
	return c, nil
}

func (c *sm4Cipher) BlockSize() int {
	return BlockSize
}

func (c *sm4Cipher) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	processBlock(c.enc, src, dst)
}

func (c *sm4Cipher) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	processBlock(c.dec, src, dst)
}

func expandKey(key []byte, forEnc bool) []uint32 {
	var mK [4]uint32
	mK[0] = binary.BigEndian.Uint32(key[0:4])
	mK[1] = binary.BigEndian.Uint32(key[4:8])
	mK[2] = binary.BigEndian.Uint32(key[8:12])
	mK[3] = binary.BigEndian.Uint32(key[12:16])

	var x [5]uint32
	x[0] = mK[0] ^ fK[0]
	x[1] = mK[1] ^ fK[1]
	x[2] = mK[2] ^ fK[2]
	x[3] = mK[3] ^ fK[3]

	var rk [32]uint32
	if forEnc {
		for i := 0; i < 32; i++ {
			x[(i+4)%5] = encRound(x[i%5], x[(i+1)%5], x[(i+2)%5], x[(i+3)%5], x[(i+4)%5], rk[:], i)
		}
	} else {
		for i := 0; i < 32; i++ {
			x[(i+4)%5] = decRound(x[i%5], x[(i+1)%5], x[(i+2)%5], x[(i+3)%5], x[(i+4)%5], rk[:], i)
		}
	}
	return rk[:]
}

func tau(a uint32) uint32 {
	var aArr [4]byte
	var bArr [4]byte
	binary.BigEndian.PutUint32(aArr[:], a)
	bArr[0] = sBox[aArr[0]]
	bArr[1] = sBox[aArr[1]]
	bArr[2] = sBox[aArr[2]]
	bArr[3] = sBox[aArr[3]]
	return binary.BigEndian.Uint32(bArr[:])
}

func lAp(b uint32) uint32 {
	return b ^ bits.RotateLeft32(b, 13) ^ bits.RotateLeft32(b, 23)
}

func tAp(z uint32) uint32 {
	return lAp(tau(z))
}

func encRound(x0 uint32, x1 uint32, x2 uint32, x3 uint32, x4 uint32, rk []uint32, i int) uint32 {
	x4 = x0 ^ tAp(x1^x2^x3^cK[i])
	rk[i] = x4
	return x4
}

func decRound(x0 uint32, x1 uint32, x2 uint32, x3 uint32, x4 uint32, rk []uint32, i int) uint32 {
	x4 = x0 ^ tAp(x1^x2^x3^cK[i])
	rk[31-i] = x4
	return x4
}

func processBlock(rk []uint32, in []byte, out []byte) {
	var x [BlockSize / 4]uint32
	x[0] = binary.BigEndian.Uint32(in[0:4])
	x[1] = binary.BigEndian.Uint32(in[4:8])
	x[2] = binary.BigEndian.Uint32(in[8:12])
	x[3] = binary.BigEndian.Uint32(in[12:16])

	for i := 0; i < 32; i += 4 {
		x[0] = f0(x[:], rk[i])
		x[1] = f1(x[:], rk[i+1])
		x[2] = f2(x[:], rk[i+2])
		x[3] = f3(x[:], rk[i+3])
	}
	r(x[:])

	binary.BigEndian.PutUint32(out[0:4], x[0])
	binary.BigEndian.PutUint32(out[4:8], x[1])
	binary.BigEndian.PutUint32(out[8:12], x[2])
	binary.BigEndian.PutUint32(out[12:16], x[3])
}

func l(b uint32) uint32 {
	return b ^ bits.RotateLeft32(b, 2) ^ bits.RotateLeft32(b, 10) ^
		bits.RotateLeft32(b, 18) ^ bits.RotateLeft32(b, 24)
}

func t(z uint32) uint32 {
	return l(tau(z))
}

func r(a []uint32) {
	a[0] = a[0] ^ a[3]
	a[3] = a[0] ^ a[3]
	a[0] = a[0] ^ a[3]
	a[1] = a[1] ^ a[2]
	a[2] = a[1] ^ a[2]
	a[1] = a[1] ^ a[2]
}

func f0(x []uint32, rk uint32) uint32 {
	return x[0] ^ t(x[1]^x[2]^x[3]^rk)
}

func f1(x []uint32, rk uint32) uint32 {
	return x[1] ^ t(x[2]^x[3]^x[0]^rk)
}

func f2(x []uint32, rk uint32) uint32 {
	return x[2] ^ t(x[3]^x[0]^x[1]^rk)
}

func f3(x []uint32, rk uint32) uint32 {
	return x[3] ^ t(x[0]^x[1]^x[2]^rk)
}

// 输入的plainText长度必须是BlockSize(16)的整数倍，也就是调用该方法前调用方需先加好padding，
// 可调用util.PKCS5Padding()方法进行加padding操作
func ECBEncrypt(key, plainText []byte) (cipherText []byte, err error) {
	plainTextLen := len(plainText)
	if plainTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}

	c, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipherText = make([]byte, plainTextLen)
	for i := 0; i < plainTextLen; i += BlockSize {
		c.Encrypt(cipherText[i:i+BlockSize], plainText[i:i+BlockSize])
	}
	return cipherText, nil
}

// 输出的plainText是加padding的明文，调用方需要自己去padding，
// 可调用util.PKCS5UnPadding()方法进行去padding操作
func ECBDecrypt(key, cipherText []byte) (plainText []byte, err error) {
	cipherTextLen := len(cipherText)
	if cipherTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}

	c, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	plainText = make([]byte, cipherTextLen)
	for i := 0; i < cipherTextLen; i += BlockSize {
		c.Decrypt(plainText[i:i+BlockSize], cipherText[i:i+BlockSize])
	}
	return plainText, nil
}

// 输入的plainText长度必须是BlockSize(16)的整数倍，也就是调用该方法前调用方需先加好padding，
// 可调用util.PKCS5Padding()方法进行加padding操作
func CBCEncrypt(key, iv, plainText []byte) (cipherText []byte, err error) {
	plainTextLen := len(plainText)
	if plainTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}

	c, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(c, iv)
	cipherText = make([]byte, plainTextLen)
	encrypter.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

// 输出的plainText是加padding的明文，调用方需要自己去padding，
// 可调用util.PKCS5UnPadding()方法进行去padding操作
func CBCDecrypt(key, iv, cipherText []byte) (plainText []byte, err error) {
	cipherTextLen := len(cipherText)
	if cipherTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}

	c, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(c, iv)
	plainText = make([]byte, len(cipherText))
	decrypter.CryptBlocks(plainText, cipherText)
	return plainText, nil
}


func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}