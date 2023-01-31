package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type Wallet struct {
	Privatekey string
	Publickey  string
}
type SignData struct {
	R string
	S string
}

// NewGreeter makes a new Greeter.
func GenerateWallet() (*Wallet, error) {
	w, err := generateKey()
	return &w, err
}

// NewGreeter makes a new Greeter.
func Sign(privateKey string, hash string) (*SignData, error) {
	priv := privateKeyFromHexString(privateKey)
	r1, s1, err1 := ecdsa.Sign(rand.Reader, priv, []byte(hash))
	var signData SignData
	if err1 != nil {
		return nil, err1
	}
	signData.R = fmt.Sprint(r1)
	signData.S = fmt.Sprint(s1)
	// fmt.Printf("\n")
	// fmt.Printf("signData.R :" + signData.R)
	// fmt.Printf("\n")
	// fmt.Printf("signData.S :" + signData.S)
	// fmt.Printf("\n")
	// fmt.Printf("\n")
	// fmt.Printf("r1.String() :" + r1.String())
	// fmt.Printf("\n")
	// fmt.Printf("s1.String() :" + s1.String())
	// fmt.Printf("\n")
	return &signData, nil
}

func Verify(publickey string, hash string, r string, s string) bool {
	pubkey := publickKeyFromHexString(publickey)
	//fmt.Println("publickey : " + publickey)
	/*
		nR := new(big.Int)
		_, err := fmt.Sscan(r, nR)
		if err != nil {
			log.Println("error scanning value:", err)
		}

		nS := new(big.Int)
		_, err2 := fmt.Sscan(s, nS)
		if err2 != nil {
			log.Println("error scanning value:", err2)
		}
	*/
	//fmt.Printf("hash         :" + hash)
	nR := new(big.Int)
	nR, rOk := nR.SetString(r, 10)
	if !rOk {
		//fmt.Println("SetString: error")
		return false
	}

	nS := new(big.Int)
	nS, sOK := nS.SetString(s, 10)
	if !sOK {
		//fmt.Println("SetString: error")
		return false
	}
	/*
	fmt.Printf("\n")
	fmt.Printf("nR         : %d", nR)
	fmt.Printf("\n")
	fmt.Printf("nS         : %d", nS)
	fmt.Printf("\n")
	*/
	if ecdsa.Verify(&pubkey, []byte(hash), nR, nS) == false {
		//log.Fatal("verify failed")
		return false
	}
	return true
}

func main() {
	c := elliptic.P256()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	privatekeystring := privateToHexString(priv)
	priv2 := privateKeyFromHexString(privatekeystring)
	privatekeystring2 := privateToHexString(priv2)
	publickeystring := publickToHexString(priv2.PublicKey)
	publickeystring2 := publickToHexString(priv.PublicKey)
	pub2 := publickKeyFromHexString(publickeystring)

	fmt.Printf("\n")
	fmt.Println("privatekeystring 1 : " + privatekeystring)
	fmt.Printf("\n")
	fmt.Println("privatekeystring 2 : " + privatekeystring2)
	fmt.Printf("\n")
	fmt.Println("publickeystring 1 : " + publickeystring)
	fmt.Printf("\n")
	fmt.Println("publickeystring 2 : " + publickeystring2)
	fmt.Printf("\n")

	hash := []byte("Hello, please use ecdsa sign me")

	r, s, err := sign(priv, hash)
	if err != nil {
		log.Fatal(err)
	}

	if !verify(&pub2, hash, r, s) {
		log.Fatal("verify failed")
	} else {
		fmt.Printf("\nverify OK\n\n")
	}
}

// GenerateKey new wallet
func generateKey() (Wallet, error) {
	var wallet Wallet
	c := elliptic.P256()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return wallet, err
	}
	privatekeystring := privateToHexString(priv)
	publickeystring := publickToHexString(priv.PublicKey)

	if len(privatekeystring) == 0 {
		return wallet, errors.New("Can't make privateKey")
	}

	if len(publickeystring) == 0 {
		return wallet, errors.New("Can't make Publickey")
	}
	wallet.Privatekey = privatekeystring
	wallet.Publickey = publickeystring
	return wallet, nil
}
func publickToHexString(pub ecdsa.PublicKey) string {
	pkh := elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
	// result := b58encode(pkh)
	// fmt.Printf("\n")
	// fmt.Println("strToHex b58encode : " + result)
	// fmt.Printf("\n")
	return b58encode(pkh)
}
func privateToHexString(priv *ecdsa.PrivateKey) string {
	return hex.EncodeToString(priv.D.Bytes())
}
func sign(priv *ecdsa.PrivateKey, hash []byte) (r, s *big.Int, err error) {
	r1, s1, err1 := ecdsa.Sign(rand.Reader, priv, hash)
	if err1 != nil {
		return nil, nil, err1
	}
	return r1, s1, nil
}
func verify(pub *ecdsa.PublicKey, hash []byte, r, s *big.Int) bool {
	if !ecdsa.Verify(pub, hash, r, s) {
		log.Fatal("verify failed")
		return false
	}
	return true
}
func publickKeyFromHexString(key string) ecdsa.PublicKey {
	b, _ := b58decode(key)
	// b := fromHexToByte(bKey)
	x, y := elliptic.Unmarshal(elliptic.P256(), b)
	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return pubKey
}
func privateKeyFromHexString(key string) *ecdsa.PrivateKey {
	b := fromHex(key)

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = new(big.Int).SetBytes(b)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(b)
	return priv
}

// Bytes2Hex returns the hexadecimal encoding of d.
func bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}
func toHex(b []byte) string {
	hex := bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}
func strToHex(str string) string {
	hex := bytes2Hex([]byte(str))
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}
func fromHexToByte(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return []byte(s)
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func fromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex2Bytes(s)
}

// b58encode encodes a byte slice b into a base-58 encoded string
func b58encode(b []byte) (s string) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	const base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	/* Convert big endian bytes to big int */
	x := new(big.Int).SetBytes(b)

	/* Initialize */
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x / 58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(base58Table[r.Int64()]) + s
	}

	return s
}

// b58decode decodes a base-58 encoded string into a byte slice b.
func b58decode(s string) (b []byte, err error) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	const base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	/* Initialize */
	x := big.NewInt(0)
	m := big.NewInt(58)

	/* Convert string to big int */
	for i := 0; i < len(s); i++ {
		b58index := strings.IndexByte(base58Table, s[i])
		if b58index == -1 {
			return nil, errors.New("Invalid base-58 character encountered")
		}
		b58value := big.NewInt(int64(b58index))
		x.Mul(x, m)
		x.Add(x, b58value)
	}

	/* Convert big int to big endian bytes */
	b = x.Bytes()

	return b, nil
}
