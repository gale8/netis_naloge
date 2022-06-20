package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	// HANDLE ENDPOINTS
	router.HandleFunc("/sign", Sign).Methods("POST")
	router.HandleFunc("/public", Public).Methods("POST")
	router.HandleFunc("/validate", nil).Methods("POST")
	// RUN SERVER
	err := http.ListenAndServe(":8080", router)
	fmt.Println("Server is listening on  port 8080.")
	if err != nil {
		panic(err)
	}
}

type Item struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

func getRequestData(w http.ResponseWriter, r *http.Request, isBodyRequired bool) (string, []byte, error) {
	keyName := r.URL.Query().Get("keyName")
	if keyName == "" {
		return "", []byte{}, errors.New("insufficient data")
	}
	if !isBodyRequired {
		return keyName, nil, nil
	}

	// DECODE BODY
	var bodyData Item
	decoder := json.NewDecoder(r.Body)
	decodingErr := decoder.Decode(&bodyData)
	if decodingErr != nil {
		return "", []byte{}, decodingErr
	}
	marshal, err := json.Marshal(bodyData)
	if err != nil {
		return "", []byte{}, err
	}

	return keyName, marshal, nil
}

func getKeyFileData(keyName string, w http.ResponseWriter) ([]byte, error) {
	keysFile := "../Naloga2/keys/" + keyName + "_rsa.pem"
	// GET KEY FROM OUR KEYS FOLDER
	if _, err := os.Stat(keysFile); !os.IsNotExist(err) {
		fileContentInBytes, err := ioutil.ReadFile(keysFile)
		if err != nil {
			return []byte{}, err
		} else {
			return fileContentInBytes, nil
		}
	} else {
		return []byte{}, err
	}
}

func Sign(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /sign")
	keyName, bodyDataBytes, err := getRequestData(w, r, true)
	if err != nil {
		sendErrorResponse(w, "Error decoding body data.")
		return
	}
	fileContentInBytes, err := getKeyFileData(keyName, w)
	if err != nil {
		sendErrorResponse(w, "Error fetching Private key data")
		return
	}
	// decode and get private key string
	pemFormattedBlock, _ := pem.Decode(fileContentInBytes)
	privateKey, _ := x509.ParsePKCS1PrivateKey(pemFormattedBlock.Bytes)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, nil)
	if err != nil {
		sendErrorResponse(w, "Signer error.")
		panic(err)
	}
	jwsObj, err := signer.Sign(bodyDataBytes)
	if err != nil {
		panic(err)
	}
	err = sendOkResponse(w, jwsObj.Signatures[0].Signature)
	if err != nil {
		sendErrorResponse(w, "Error converting struct data to bytes.")
		panic(err)
	}
}

func Public(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /public")
	keyName, _, err := getRequestData(w, r, false)
	if err != nil {
		sendErrorResponse(w, "Error decoding body data.")
		return
	}
	fileContentInBytes, err := getKeyFileData(keyName, w)
	if err != nil {
		sendErrorResponse(w, "Error fetching Private key data")
		return
	}
	// decode and get private key string
	pemFormattedBlock, _ := pem.Decode(fileContentInBytes)
	privateKey, _ := x509.ParsePKCS1PrivateKey(pemFormattedBlock.Bytes)
	publicKey := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}
	publicKeyPem := string(pem.EncodeToMemory(&pubKeyBlock))
	err = sendOkResponse(w, publicKeyPem)
	if err != nil {
		sendErrorResponse(w, "Error converting struct data to bytes.")
		panic(err)
	}
}

func sendOkResponse(w http.ResponseWriter, respData interface{}) error {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(respData)
	return err
}

func sendErrorResponse(w http.ResponseWriter, errMsg string) {
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(Response{
		Message: errMsg,
	})
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Message string
}
