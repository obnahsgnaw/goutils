package esutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
	"github.com/obnahsgnaw/goutils/security/coder"
	"sync"
)

/*
DES（Data Encryption Standard） 数据块大小为8个字节 密钥长度是64位（其中8位用于校验）(8byte) 3DES（即Triple DES）是DES向AES过渡的加密算法
AES (Advanced Encryption Standard，高级加密标准) AES的数据块大小为16个字节 密钥长度是128位（AES算法比DES算法更安全） 最终生成的加密密钥长度有128位、192位和256位这三种
AES主要有五种工作模式(其实还有很多模式) ：
	ECB (Electroniccodebook，电子密码本)、  不需要初始化向量（IV）相同明文得到相同的密文
	CBC (Cipher-block chaining，密码分组链接)、 第一个明文块与一个叫初始化向量的数据块进行逻辑异或运算。这样就有效的解决了ECB模式所暴露出来的问题，即使两个明文块相同，加密后得到的密文块也不相同。但是缺点也相当明显，如加密过程复杂，效率低等
	CFB (Cipher feedback，密文反馈)、 CFB模式能够将密文转化成为流密文 不需要填充
	OFB (Output feedback，输出反馈)、 不再直接加密明文块，其加密过程是先使用块加密器生成密钥流，然后再将密钥流和明文流进行逻辑异或运算得到密文流
	PCBC (Propagating cipher-block chaining，增强型密码分组链接)
*/

func esErr(msg string, err error) error {
	if err == nil {
		return errors.New("esutil:" + msg)
	}

	return fmt.Errorf("esutil:"+msg+": %w", err)
}

var (
	ErrIvLengthError  = esErr(fmt.Sprintf("iv size error, aes=%d, des=%d", aes.BlockSize, des.BlockSize), nil)
	ErrModeNotSupport = esErr("mode not support now", nil)
)

// ADes Aes Des
type ADes struct {
	t         EsType
	m         EsMode
	aesBlocks sync.Map
	desBlocks sync.Map
	disable   bool
	encoder   coder.Encoder
}

func New(esType EsType, mode EsMode, o ...Option) *ADes {
	s := &ADes{
		t:         esType,
		m:         mode,
		aesBlocks: sync.Map{},
		desBlocks: sync.Map{},
		encoder:   coder.B64StdEncoding,
	}
	s.with(o...)
	return s
}

func (e *ADes) with(o ...Option) {
	for _, opt := range o {
		if opt != nil {
			opt(e)
		}
	}
}

func (e *ADes) Type() EsType {
	return e.t
}

func (e *ADes) Mode() EsMode {
	return e.m
}

func (e *ADes) Disable() {
	e.disable = true
}

func (e *ADes) Encrypt(data, key []byte, encode bool) (encrypted, iv []byte, err error) {
	iv = e.t.RandIv()
	encrypted, err = e.EncryptWithIv(data, key, encode, iv)
	return
}

func (e *ADes) EncryptWithIv(data, key []byte, encode bool, iv []byte) (encrypted []byte, err error) {
	if len(data) == 0 {
		return
	}
	if !e.disable {
		var block cipher.Block
		var esBlock *Block
		if esBlock, err = e.getEsBlock(key); err != nil {
			return
		}
		if block, err = e.getModeBlock(esBlock, iv, true); err != nil {
			return
		}
		padData := pkcs7Padding(data, block.BlockSize())

		encrypted = make([]byte, len(padData))
		block.Encrypt(encrypted, padData)
	} else {
		encrypted = data
	}
	if encode {
		encrypted = []byte(e.encoder.EncodeToString(encrypted))
	}
	return
}

func (e *ADes) Decrypt(encrypted, key, iv []byte, decode bool) (data []byte, err error) {
	if len(encrypted) == 0 {
		return
	}
	var block cipher.Block
	var esBlock *Block
	if decode {
		if encrypted, err = e.encoder.DecodeString(string(encrypted)); err != nil {
			return
		}
	}
	if !e.disable {
		if esBlock, err = e.getEsBlock(key); err != nil {
			return
		}
		if block, err = e.getModeBlock(esBlock, iv, false); err != nil {
			return
		}

		decryptData := make([]byte, len(encrypted))
		block.Decrypt(decryptData, encrypted)

		data = pkcs7UnPadding(decryptData)
	} else {
		data = encrypted
	}
	return
}

func (e *ADes) getEsBlock(key []byte) (block *Block, err error) {
	kName := string(key)
	if e.t == Des {
		if b, ok := e.desBlocks.Load(kName); ok {
			return b.(*Block), nil
		} else {
			block, err = desBlock(key)
			if err != nil {
				return nil, err
			}
			e.desBlocks.Store(kName, block)
			return
		}
	} else {
		if b, ok := e.aesBlocks.Load(kName); ok {
			return b.(*Block), nil
		} else {
			block, err = aesBlock(key)
			if err != nil {
				return nil, err
			}
			e.desBlocks.Store(kName, block)
			return
		}
	}
}

func (e *ADes) getModeBlock(block *Block, iv []byte, enc bool) (b cipher.Block, err error) {
	if e.m == CbcMode {
		b, err = block.CbcBlock(iv, enc)
	} else {
		err = ErrModeNotSupport
	}
	return
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	num := blockSize - len(data)%blockSize
	padData := bytes.Repeat([]byte{byte(num)}, num)
	data = append(data, padData...)
	return data
}

func pkcs7UnPadding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	num := int(data[len(data)-1])

	return data[:len(data)-num]
}
