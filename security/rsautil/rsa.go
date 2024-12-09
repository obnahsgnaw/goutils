package rsautil

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/hsutil"
)

/*
RSA:
1. 公钥加密， 私钥解密
2. 私钥签名， 公钥验签
3. 公钥长度size bit决定加解密块的长度， 最大 size/8 - 7 byte
4. padding方式 PKCS 1.5,  OAEP
*/

var (
	ErrPublicKeyParseError  = errors.New("rsautil: public key parse failed")
	ErrPublicKeyError       = errors.New("rsautil: public key error")
	ErrPrivateKeyParseError = errors.New("rsautil: private key parse failed")
	ErrBitTooShort          = errors.New("rsautil: bits too short")
)

const (
	pukPKCS1 = 1
	pukPKIX  = 2
	prkPKCS1 = 1
	prkPKCS8 = 2
)

type Rsa struct {
	encoder  coder.Encoder
	signHash crypto.Hash
	pukType  int // 1 = PKCS1 2 = PKIX
	prkType  int // 1 = PKCS1 2 = PKCS8
	disable  bool
}

func New(o ...Option) *Rsa {
	s := &Rsa{
		encoder:  coder.B64StdEncoding,
		signHash: crypto.SHA256,
		disable:  false,
		pukType:  pukPKCS1,
		prkType:  prkPKCS1,
	}
	s.with(o...)
	return s
}

func (s *Rsa) with(o ...Option) {
	for _, opt := range o {
		if opt != nil {
			opt(s)
		}
	}
}

// Generate rsa private key and public key size: 密钥位数bit，加密的message不能比密钥长 (size/8 -11)
func (s *Rsa) Generate(bits int) (privateKey []byte, publicKey []byte, err error) {
	if bits < 512 {
		err = ErrBitTooShort
		return
	}
	var priKey *rsa.PrivateKey
	priKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	var paBuf = bytes.NewBufferString("")
	err = pem.Encode(paBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priKey),
	})
	if err != nil {
		return
	}
	privateKey = paBuf.Bytes()

	var derStream []byte
	if s.pukType == pukPKIX {
		derStream, err = x509.MarshalPKIXPublicKey(&priKey.PublicKey)
		if err != nil {
			return
		}
	} else {
		derStream = x509.MarshalPKCS1PublicKey(&priKey.PublicKey)
	}

	var puBuf = bytes.NewBufferString("")
	err = pem.Encode(puBuf, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derStream,
	})

	if err != nil {
		return
	}
	publicKey = puBuf.Bytes()

	return
}

// Encrypt public key encrypt
func (s *Rsa) Encrypt(data []byte, pubKey []byte, encode bool) (encrypted []byte, err error) {
	if len(data) == 0 {
		return
	}
	if !s.disable {
		var publicKey *rsa.PublicKey
		var chunkData []byte
		if publicKey, err = s.getPublicKey(pubKey); err != nil {
			return
		}
		maxLen := publicKey.N.BitLen()/8 - 11 // EncryptPKCS1v15 11位填充
		if maxLen < 0 {
			err = ErrBitTooShort
			return
		}
		chunks := split(data, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			if chunkData, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk); err != nil {
				return
			}
			buffer.Write(chunkData)
		}
		encrypted = buffer.Bytes()
	} else {
		encrypted = data
	}
	if encode {
		encrypted = []byte(s.encoder.EncodeToString(encrypted))
	}

	return
}

// Decrypt private key decrypt
func (s *Rsa) Decrypt(encrypted []byte, priKey []byte, decode bool) (data []byte, err error) {
	if len(encrypted) == 0 {
		return
	}
	if decode {
		if encrypted, err = s.encoder.DecodeString(string(encrypted)); err != nil {
			return
		}
	}

	if !s.disable {
		var privateKey *rsa.PrivateKey
		var chunkData []byte
		if privateKey, err = s.getPrivateKey(priKey); err != nil {
			return
		}
		maxLen := privateKey.PublicKey.N.BitLen() / 8
		chunks := split(encrypted, maxLen)
		buffer := bytes.NewBufferString("")
		for _, chunk := range chunks {
			if chunkData, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk); err != nil {
				return
			}
			buffer.Write(chunkData)
		}
		data = buffer.Bytes()
	} else {
		data = encrypted
	}

	return
}

// Sign Private key sign
func (s *Rsa) Sign(data []byte, priKey []byte, encode bool) (signature []byte, err error) {
	if len(data) == 0 {
		return
	}

	var privateKey *rsa.PrivateKey
	var hashed []byte
	var chunkData []byte
	if privateKey, err = s.getPrivateKey(priKey); err != nil {
		return
	}
	maxLen := privateKey.PublicKey.N.BitLen()/8 - 11 - s.signHash.Size()
	if maxLen < 0 {
		err = ErrBitTooShort
		return
	}
	chunks := split(data, maxLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		if hashed, err = hsutil.Hash(chunk, s.signHash); err != nil {
			return
		}
		if chunkData, err = rsa.SignPKCS1v15(rand.Reader, privateKey, s.signHash, hashed); err != nil {
			return
		}
		buffer.Write(chunkData)
	}

	signature = buffer.Bytes()
	if encode {
		signature = []byte(s.encoder.EncodeToString(signature))
	}

	return
}

// Verify Public key verify
func (s *Rsa) Verify(data, signature, pubKey []byte, decode bool) (err error) {
	if len(data) == 0 {
		return
	}

	if decode {
		if signature, err = s.encoder.DecodeString(string(signature)); err != nil {
			return
		}
	}

	var publicKey *rsa.PublicKey
	if publicKey, err = s.getPublicKey(pubKey); err != nil {
		return
	}

	var hashed []byte
	maxLen := publicKey.N.BitLen()/8 - 11 - s.signHash.Size()
	if maxLen < 0 {
		err = ErrBitTooShort
		return
	}
	chunks := split(data, maxLen)
	signLen := len(signature) / len(chunks)
	for i, chunk := range chunks {
		if hashed, err = hsutil.Hash(chunk, s.signHash); err != nil {
			return
		}
		chunkSign := signature[i*signLen : i*signLen+signLen]
		if err = rsa.VerifyPKCS1v15(publicKey, s.signHash, hashed, chunkSign); err != nil {
			return err
		}
	}

	return
}

func (s *Rsa) Disable() {
	s.disable = true
}

// get *rsa.PublicKey from byte key
func (s *Rsa) getPublicKey(pk []byte) (publicKey *rsa.PublicKey, err error) {
	var block *pem.Block
	var publicInterface interface{}
	var flag bool

	if block, _ = pem.Decode(pk); block == nil {
		err = ErrPublicKeyParseError
		return
	}

	if s.pukType == pukPKIX {
		if publicInterface, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			return
		}
	} else if s.pukType == pukPKCS1 {
		if publicInterface, err = x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
			return
		}
	} else {
		err = ErrPublicKeyParseError
	}

	if publicKey, flag = publicInterface.(*rsa.PublicKey); !flag {
		err = ErrPublicKeyError
	}

	return
}

// get *rsa.PrivateKey from byte key
func (s *Rsa) getPrivateKey(sk []byte) (privateKey *rsa.PrivateKey, err error) {
	var block *pem.Block

	if block, _ = pem.Decode(sk); block == nil {
		err = ErrPrivateKeyParseError
		return
	}
	if s.prkType == prkPKCS1 {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if s.prkType == prkPKCS8 {
		k, err1 := x509.ParsePKCS8PrivateKey(block.Bytes)

		if err1 != nil {
			err = err1
			return
		}
		privateKey, _ = k.(*rsa.PrivateKey)
	} else {
		err = ErrPrivateKeyParseError
	}

	return
}

// split rsa message block
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
