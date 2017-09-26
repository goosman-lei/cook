package util

import (
	"crypto/rc4"
)

type Crypto interface {
	Encrypt(decrypted []byte) (encrypted []byte, err error)
	Decrypt(encrypted []byte) (decrypted []byte, err error)
}

type mrc4 struct {
	// length between 1 and 256 byte
	key []byte
}

func NewCrypto(key []byte) Crypto {
	return &mrc4{
		key: key,
	}
}

func (r *mrc4) Encrypt(decrypted []byte) (encrypted []byte, err error) {
	var cipher *rc4.Cipher
	if cipher, err = rc4.NewCipher(r.key); err != nil {
		return
	}
	encrypted = make([]byte, len(decrypted))
	cipher.XORKeyStream(encrypted, decrypted)
	return
}

func (r *mrc4) Decrypt(encrypted []byte) (decrypted []byte, err error) {
	var cipher *rc4.Cipher
	if cipher, err = rc4.NewCipher(r.key); err != nil {
		return
	}
	decrypted = make([]byte, len(encrypted))
	cipher.XORKeyStream(decrypted, encrypted)
	return
}
