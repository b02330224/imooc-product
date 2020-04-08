package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//高级加密标准 (Advanced Encryption Standard)AES

//16,24,32位字符串的话，分别对应AES-128,AES-192,AES-256加密方法
//key不能泄漏
var PwdKey=[]byte("DIS**#KKKDJJSKDI")

//PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//填充的反向操作，删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	} else {
		padding := int(origData[length-1])
		return origData[:(length - padding)], nil
	}
}

//实现解密
func AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))

	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)

	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, nil
}

//实现加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//获取块大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//加密base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), err
}

func DePwdCode(pwd string) ([]byte, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}

	//执行AES解密
	return AesDeCrypt(pwdByte, PwdKey)
}


