package sm4

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)



func SM4Encrypt(src []byte, key []byte) []byte {

	block, e := sm4.NewCipher(key)
	if e != nil {
		fmt.Println("newCrypther faild !")
	}
	a := block.BlockSize() - len(src)%block.BlockSize()
	repeat := bytes.Repeat([]byte{byte(a)},a)
	newsrc := append(src, repeat...)

	dst := make([]byte, len(newsrc))
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(dst,newsrc)
	return dst

}

func SM4Decrypt(dst,key []byte)[]byte {

	block, e := sm4.NewCipher(key)
	if e != nil {
		fmt.Println("newcipher faild! ")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])

	src := make([]byte, len(dst))
	blockMode.CryptBlocks(src, dst)

	num := int(src[len(src)-1])
	newsrc := src[:len(src)-num]
	return newsrc
}

/*
	a := []byte("333333")
	key := []byte("9A8159B49AB10B4BA55BC477B1C0E79E")
	decrypto := SM4Encrypt(a,key)
	fmt.Println("sm4加密后：",hex.EncodeToString(decrypto))
	i := SM4Decrypt(decrypto, key)
	fmt.Println("sm4解密后：",string(i))
*/