package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type Wallet struct {
	salt              string
	nDeriveIterations uint32
	crypted_key       []byte
	ckey              []byte
	pubkey_hash       []byte
	pubkey            *big.Int
}

var m_curve = secp256k1.S256()

// Calc pubkey hash
func GetPubkeyHash(pubkey string) []byte {
	pubkey_bc, _ := hex.DecodeString(pubkey)
	data := sha256.Sum256(pubkey_bc)
	data = sha256.Sum256(data[:])
	return data[0:16]
}

// New
func NewWallet(options *WalletOptions) *Wallet {
	w := new(Wallet)

	// Supports compressed pubkey only
	pubkey := options.Pubkey
	if len(pubkey) != 66 || !(strings.HasPrefix(pubkey, "02") || strings.HasPrefix(pubkey, "03")) {
		log.Panicf("Invalid pubkey: %s\n", pubkey)
	}

	salt, _ := hex.DecodeString(options.Salt)
	w.salt = string(salt)
	w.nDeriveIterations = options.DeriveIterations
	w.crypted_key, _ = hex.DecodeString(options.Crypted_key)
	w.ckey, _ = hex.DecodeString(options.Ckey)
	w.pubkey_hash = GetPubkeyHash(options.Pubkey)
	w.pubkey, _ = new(big.Int).SetString(options.Pubkey[2:], 16)
	return w
}

func (w *Wallet) CheckPassword(password string) bool {
	secret := w.genPrivateKey(password)
	x, _ := m_curve.ScalarBaseMult(secret)
	return x.Cmp(w.pubkey) == 0
}

// Generate private key
func (w *Wallet) genPrivateKey(password string) []byte {
	buffer := make([]byte, 64)
	data := sha512.Sum512([]byte(password + w.salt))
	for i := uint32(1); i < w.nDeriveIterations; i++ {
		data = sha512.Sum512(data[:])
	}
	chKey := data[0:32]
	chIV := data[32:48]

	block, _ := aes.NewCipher(chKey)
	cbc := cipher.NewCBCDecrypter(block, chIV)
	cbc.CryptBlocks(buffer, w.crypted_key)

	chKey = buffer[0:32]
	chIV = w.pubkey_hash

	block, _ = aes.NewCipher(chKey)
	cbc = cipher.NewCBCDecrypter(block, chIV)
	cbc.CryptBlocks(buffer, w.ckey)

	return buffer[0:32]
}
