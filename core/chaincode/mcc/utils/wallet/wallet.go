package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	
	enc "github.com/btcsuite/btcutil/base58"
)

type Wallet struct {
	Privatekey string
	Publickey  string
}
type SignData struct {
	R string
	S string
}

type PubKeyData struct {
	X string
	Y string
}

// Generate new wallet
func GenerateWallet() (*Wallet, error) {
	w, err := generateKey()
	return &w, err
}

func generateKey() (Wallet, error) {
	var wallet Wallet

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	encPriv, encPub := encode(privateKey, publicKey)

	wallet.Privatekey = encPriv
	wallet.Publickey = encPub

	fmt.Println("public  Key ", encPub)
	fmt.Println("private Key ", encPriv)

	return wallet, nil
}

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {

	pubK := elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)
	encodePub := b58encode(pubK)

	encodePriv := privateToHexString(privateKey)

	return encodePriv, encodePub
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

func Verify(pubKey string, msg string, R, S string) bool {

	digest := sha1.Sum([]byte(msg))

	rbigInt := toBigInt(R)
	sbigInt := toBigInt(S)

	println()

	pub := publickKeyFromHexString(pubKey)

	return ecdsa.Verify(&pub, digest[:], rbigInt, sbigInt)

}

func toBigInt(val string) *big.Int {
	valBigInt := new(big.Int)
	valBigInt, ok := valBigInt.SetString(val, 10)
	if !ok {
		fmt.Println("SetString: error val")
		return nil
	}
	return valBigInt
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

/*
func GetHash(s string) string {
	/// analog java sha-256

	sum224 := sha256.Sum256([]byte(s))
	sh := make([]byte, sha256.Size)

	for i := range sum224 {
		sh[i] = byte(sum224[i])
	}
	fmt.Printf("%d", sh)
	fmt.Println()

	hashTostr := b58encode(sh)
	return hashTostr
}
*/

func GetHash(s string) string {
	/// analog java sha-256

	sum224 := sha256.Sum256([]byte(s))
	signed := make([]int8, sha256.Size)
	for i := range sum224 {
		signed[i] = int8(sum224[i])
	}
	//	fmt.Printf("%d", signed)
	//	fmt.Println()

	unsigned := make([]byte, 0, len(signed))
	for _, b := range signed {
		unsigned = append(unsigned, byte(b))
	}

	// fmt.Println("unsigned", unsigned)

	hashTostr := enc.Encode(unsigned)

	return hashTostr
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

func privateToHexString(priv *ecdsa.PrivateKey) string {
	return hex.EncodeToString(priv.D.Bytes())
}

// NewGreeter makes a new Greeter.
func getPubKeyData(pubK string) *PubKeyData {
	pkd := getXandYfromPubKey(pubK)
	return &pkd
}

func getXandYfromPubKey(pubK string) PubKeyData {

	var pubKeyData PubKeyData

	pub := publickKeyFromHexString(pubK)

	fmt.Println("X for public", pub.X)
	fmt.Println("Y for public", pub.Y)

	pubKeyData.X = fmt.Sprint(pub.X)
	pubKeyData.Y = fmt.Sprint(pub.Y)

	return pubKeyData
}

//// additional

func Sign(privateKey string, hash string) (*SignData, error) {
	priv := privateKeyFromHexString(privateKey)
	r1, s1, err1 := ecdsa.Sign(rand.Reader, priv, []byte(hash))
	var signData SignData
	if err1 != nil {
		return nil, err1
	}
	signData.R = fmt.Sprint(r1)
	signData.S = fmt.Sprint(s1)

	return &signData, nil
}

func privateKeyFromHexString(key string) *ecdsa.PrivateKey {
	b := fromHex(key)

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = new(big.Int).SetBytes(b)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(b)
	return priv
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

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

//

func GetPubKeyFromXandY(x string, y string) string {

	xPoint := fmt.Sprintf(strings.Replace(x, " ", "", -1))
	yPoint := fmt.Sprintf(strings.Replace(y, " ", "", -1))

	X := new(big.Int)
	X, ok := X.SetString(xPoint, 16)
	if !ok {
		fmt.Println("SetString: error val")
		return ""
	}

	Y := new(big.Int)
	Y, okk := Y.SetString(yPoint, 16)
	if !okk {
		fmt.Println("SetString: error val")
		return ""
	}

	pubKey := elliptic.Marshal(elliptic.P256(), X, Y)

	pubKeyStr := b58encode(pubKey)

	return pubKeyStr
}

// return pub key
func GetPubKeyHash(pubKey string) string {

	//pubKey := GetPubKeyFromXandY(x, y)

	fmt.Println("public key received", pubKey)

	pKeyHash := GetHash(pubKey)
	finKey := pKeyHash[len(pKeyHash)-40 : len(pKeyHash)]

	return finKey
}