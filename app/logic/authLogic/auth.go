package authLogic

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"develop-template/app/entity/resp/authRespEntity"
	config2 "develop-template/helper/config"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	time2 "time"
)

type logic struct {
	Ctx  context.Context
	GCtx *gin.Context
}

func NewLogic(ctx *gin.Context) *logic {
	return &logic{Ctx: ctx, GCtx: ctx}
}

// GetCheckTokenData
//
//	@Description: http://127.0.0.1:8080/Auth/2b33386cb6d14ab7b5f6738f7fc1704cgo/X-Dweb-Host/file.sys.dweb:443
//	@receiver l
//	@return res
func (l *logic) GetCheckTokenData() (res *authRespEntity.RespCheckData) {
	defer func() {
		r := recover()
		fmt.Println("============panic GetCheckTokenData============", r)
	}()
	token := l.GCtx.Param("token") // 可以获取路径中的 name 参数
	//action := l.GCtx.Param("action") // 可以获取 *action 之后的所有路径
	authToken := config2.GetKeyValue("domain", "authToken")
	fmt.Println("authToken Token: ", authToken)
	var data = authToken

	if token == authToken {
		prvKey := config2.GetContentFromPem("privateKey")
		pubKey := config2.GetContentFromPem("publicKey")
		fmt.Println(prvKey)
		fmt.Println(pubKey)
		fmt.Println("-------------------------------进行签名与验证操作-----------------------------------------")
		fmt.Println("对消息进行签名操作...")
		signData := RsaSignWithSha256([]byte(data), []byte(prvKey))
		fmt.Println("消息的签名信息： ", hex.EncodeToString(signData))
		fmt.Println("\n对签名信息进行验证...")
		if RsaVerySignWithSha256([]byte(data), signData, []byte(pubKey)) {
			fmt.Println("签名信息验证成功，确定是正确私钥签名！！")
		}

		fmt.Println("-------------------------------进行加密解密操作-----------------------------------------")
		ciphertext := RsaEncrypt([]byte(data), []byte(pubKey))
		fmt.Println("公钥加密后的数据：", hex.EncodeToString(ciphertext))
		sourceData := RsaDecrypt(ciphertext, []byte(prvKey))
		fmt.Println("私钥解密后的数据：", string(sourceData))
		//私钥解密后的数据
		res = new(authRespEntity.RespCheckData)
		res.Id = 0
		res.Result = false
		res.NowTime = int(time2.Now().UnixMilli())
		if string(sourceData) == data {
			res.Result = true
			res.Description = "恭喜Token验证成功！"
			return res
		} else {
			res.Description = "抱歉，解密失败，无法通过验证！！！"
			return res
		}
	}
	res.Description = "验证失败！"
	return res
}

// 签名
func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

// 验证
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// 公钥加密
func RsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// 私钥解密
func RsaDecrypt(ciphertext, keyBytes []byte) []byte {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error!"))
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}
