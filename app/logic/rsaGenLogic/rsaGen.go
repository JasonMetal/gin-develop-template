package rsaGenLogic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/gin-gonic/gin"
)

type logic struct {
	Ctx  context.Context
	GCtx *gin.Context
}

func NewLogic(ctx *gin.Context) *logic {
	return &logic{Ctx: ctx, GCtx: ctx}
}

// RSA公钥私钥产生
func (l *logic) GenRsaKey() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}
