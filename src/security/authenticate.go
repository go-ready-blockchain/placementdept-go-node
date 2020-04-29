package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/jugalw13/placementdept-go-node/utils"
)

type User struct {
	Name       string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type UserBytes struct {
	Name       string
	PrivateKey []byte
	PublicKey  []byte
}

func GetPrivateandPublicKey(Name string) *User {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
	user := &User{Name, privateKey, &privateKey.PublicKey}
	return user
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		os.Exit(1)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			os.Exit(1)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		os.Exit(1)
	}
	return key
}

func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println("Error Public Key1")
			os.Exit(1)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		fmt.Println("Error Public Key")
		os.Exit(1)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		os.Exit(1)
	}
	return key
}
func GetUserFromDB(name string) *User {
	var ubytes []byte = []byte("{UserBytes}")
	ubytes = utils.GetUserBytesFromDB(name)
	userbytes := &UserBytes{}
	json.Unmarshal(ubytes, &userbytes)
	user := &User{userbytes.Name, BytesToPrivateKey(userbytes.PrivateKey), BytesToPublicKey(userbytes.PublicKey)}
	return user

}
func PutUserIntoDB(user *User) {

	userbytes := &UserBytes{user.Name, PrivateKeyToBytes(user.PrivateKey), PublicKeyToBytes(user.PublicKey)}
	ubytes, _ := json.Marshal(userbytes)

	utils.PutUserBytesIntoDB(ubytes)

}

func GetPublicKeyFromDB(name string) *rsa.PublicKey {
	var publickeybytes []byte = []byte("{PublicKeyBytes}")
	publickeybytes = utils.GetPublicKeyFromDB(name)
	publickey := BytesToPublicKey(publickeybytes)
	return publickey
}
func PutPublicKeyIntoDB(publickey *rsa.PublicKey, name string) {

	publickeybytes := PublicKeyToBytes(publickey)

	utils.PutPublickeyIntoDB(publickeybytes, name)

}

func EncryptMessage(message []byte, receiverPublicKey *rsa.PublicKey) []byte {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		receiverPublicKey,
		message,
		label,
	)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Data Encrypted Successfully!")
	return ciphertext
}

func DecryptMessage(ciphertext []byte, receiverPrivateKey *rsa.PrivateKey) ([]byte, bool) {
	label := []byte("")
	hash := sha256.New()
	message, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		receiverPrivateKey,
		ciphertext,
		label,
	)
	if err != nil {
		fmt.Println(err)
		return []byte(""), false
		os.Exit(1)
	}
	fmt.Println("Data Decrypted Successfully!")
	return message, true
}

func PSSSignature(message []byte, privateKey *rsa.PrivateKey) []byte {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		newhash,
		hashed,
		&opts,
	)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Signature Generated Successfully!")
	return signature
}

func VerifyPSSSignature(publicKey *rsa.PublicKey, signature []byte, plainText []byte) bool {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	PSSmessage := plainText
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	err := rsa.VerifyPSS(
		publicKey,
		newhash,
		hashed,
		signature,
		&opts,
	)
	if err != nil {
		fmt.Println("Signature Verification Failed")
		return false
		os.Exit(1)
	} else {
		fmt.Println("Signature Verified Successfully!")
	}
	return true
}

func GenerateAcademicDeptKeys() {
	var AcademicDeptKeys = GetPrivateandPublicKey("AcademicDept")
	PutUserIntoDB(AcademicDeptKeys)
	PutPublicKeyIntoDB(AcademicDeptKeys.PublicKey, "AcademicDept")
	fmt.Println("Generated Public and Private keys for Academic Dept")

}
func GeneratePlacementDeptKeys() {
	var PlacementDeptKeys = GetPrivateandPublicKey("PlacementDept")
	PutUserIntoDB(PlacementDeptKeys)
	PutPublicKeyIntoDB(PlacementDeptKeys.PublicKey, "PlacementDept")
	fmt.Println("Generated Public and Private keys for Placement Dept")

}

func GenerateStudentKeys(name string) {
	var Student = GetPrivateandPublicKey(name)
	PutUserIntoDB(Student)
	PutPublicKeyIntoDB(Student.PublicKey, name)
	fmt.Println("Generated Public and Private keys for Student: ", name)
}

func GenerateCompanyKeys(companyname string) {
	var Company = GetPrivateandPublicKey(companyname)
	PutUserIntoDB(Company)
	PutPublicKeyIntoDB(Company.PublicKey, companyname)
	fmt.Println("Generated Public and Private keys for Company: ", companyname)
}

// func Test() {
// 	var sender = GetPrivateandPublicKey("alice")
// var receiver = GetPrivateandPublicKey("bob")
// message := []byte("Just A Rather Very Intelligent System")
// ciphertext := EncryptMessage(message, receiver.PublicKey)
// signature := PSSSignature(message, sender.PrivateKey)
// //fmt.Println("Encrypted Message: ", ciphertext)
// decmessage, ok := DecryptMessage(ciphertext, receiver.PrivateKey)
// flag := VerifyPSSSignature(sender.PublicKey, signature, decmessage)
// fmt.Println("Decrypted Message: ", string(decmessage), ok, flag)

// fmt.Println("Sender Public Key: ", sender.PublicKey)
// e := PublicKeyToBytes(sender.PublicKey)
// fmt.Println("Public key Bytes: ", e)
// f := BytesToPublicKey(e)
//fmt.Println("Public key: ", f, sender.PublicKey == f)

// fmt.Println(alice.Name, "\n", PrivateKeyToBytes(alice.PrivateKey), "\n", PublicKeyToBytes(alice.PublicKey))
// fmt.Println("ALICE PRIVATE KEY: ", BytesToPrivateKey(PrivateKeyToBytes(alice.PrivateKey)))

//}

// func main() {
// 	Test()
// }
