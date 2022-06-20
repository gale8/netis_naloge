package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"io/ioutil"
	"log"
	"os"
)

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
	if *file == "" {
		fmt.Println("You need to provide a file name.")
		return
	}

	content, err := getJsonFileContent(*file)
	if err != nil {
		panic(err.Error())
	}

	pemString := os.Getenv("PRIVATE_KEY")
	// DECODE and PARSE pk in structure suitable for signing
	block, _ := pem.Decode([]byte(pemString))
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	// Sign generates a signature for the given payload, and serializes it in compact serialization format.
	jwsObj, err := jws.Sign(content, jwa.RS256, privateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jwsObj))
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
