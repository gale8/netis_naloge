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
	keyName := flag.String("key_name", "", "a Private Key name")
	flag.Parse()
	if *keyName == "" {
		fmt.Println("You must provide key name.")
		return
	}

	keyFilePath := "./keys/" + (*keyName) + "_rsa.pem"
	// CHECK IF FILE WITH THE NAME 'key_name'_rsa.pem already EXISTS
	if _, err := os.Stat(keyFilePath); !os.IsNotExist(err) {
		privateKeyInBytes, err := ioutil.ReadFile(keyFilePath)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(privateKeyInBytes))
	} else {
		// Generate a new private key
		keySize := 2048
		privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			panic(err)
		}
		// Encode private key to PKCS#1 ASN.1 PEM format
		pkPem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
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
