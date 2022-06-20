package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	secretToken := flag.String("secret", "", "a Secret token")
	flag.Parse()
	fmt.Println(*secretToken)

	keySize := 2048
	// reader := strings.NewReader(*secretToken)
	// Generate private key using SECRET FROM INPUT
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		panic(err)
	}

	keyFilePath := "./keys/" + (*secretToken) + "_rsa.pem"
	// CHECK IF FILE WITH THE NAME 'secret'_rsa.pem already EXISTS
	if _, err := os.Stat(keyFilePath); !os.IsNotExist(err) {
		// FILE EXISTS - GET CONTENT
		fileContentInBytes, err := ioutil.ReadFile(keyFilePath)
		if err != nil {
			panic(err)
		}
		// DISPLAY CONTENT
		fmt.Println(string(fileContentInBytes))
	} else {
		// SAVE PRIVATE_KEY IN FILE
		// Encode private key to PKCS#1 ASN.1 PEM.
		pkPem := pem.
			EncodeToMemory(&pem.Block{
				Type: "RSA PRIVATE KEY",
				// CONVERT 'privateKey' into PKCS#1 ASN.1 DER form
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			})
		// SAVE PK TO FILE
		if err := ioutil.WriteFile(keyFilePath, pkPem, 0100); err != nil {
			panic(err)
		}
		// DISPLAY KEY
		fmt.Println(string(pkPem))
	}
}
