package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"log"
	"os"
)

// NALOGA 1
func main() {
	// LOAD ENV FILE
	errEnv := godotenv.Load("environment.env")
	if errEnv != nil {
		log.Fatalf("Error loading .env file")
	}

	// GET USER INPUT: -file
	file := flag.String("file", "", "a .json file")
	// Once all flags are declared, call flag.Parse() to execute the command-line parsing.
	flag.Parse()

	// GET PRIVATE_KEY from .ENV and DECODE/PARSE IT
	pemString := os.Getenv("Private_KEY")
	block, _ := pem.Decode([]byte(pemString))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	// SIGNER INIT with OUR PRIVATE KEY AND RS256 Signature algorithm
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	if err != nil {
		panic(err)
	}
	// GET PAYLOAD FROM OUR FILE:
	content, err := getJsonFileContent(*file)
	if err != nil {
		panic(err)
	}
	// SIGN OUR PAYLOAD DATA WITH OUR PRIVATE KEY
	jwsObj, err := signer.Sign(content)
	if err != nil {
		panic(err)
	}
	// 'jwsObj' is a PROTECTED JWS OBJECT
	fmt.Println(jwsObj.Signatures)
}

func getJsonFileContent(fileName string) ([]byte, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return []byte{}, err
	}
	// get byte array of our content
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}
