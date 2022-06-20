package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	// HANDLE ENDPOINTS
	router.HandleFunc("/sign", Sign).Methods("POST")
	router.HandleFunc("/public", Public).Methods("POST")
	router.HandleFunc("/validate", Validate).Methods("POST")
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

type JWS struct {
	JwsObject string `json:"jws_object"`
}

func getRequestData(r *http.Request, isBodyRequired bool) (string, []byte, error) {
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
	keyName, bodyDataBytes, err := getRequestData(r, true)
	if err != nil {
		sendErrorResponse(w, "Error decoding input data.")
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
	// sign json content from our request
	jwsObj, err := jws.Sign(bodyDataBytes, jwa.RS256, privateKey)
	if err != nil {
		panic(err)
	}
	err = sendOkResponse(w, string(jwsObj))
	if err != nil {
		sendErrorResponse(w, "Error when sending response.")
		panic(err)
	}
}

func Public(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /public")
	keyName, _, err := getRequestData(r, false)
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

func Validate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /sign")
	keyName, _, err := getRequestData(r, false)
	if err != nil {
		sendErrorResponse(w, "Error decoding input data.")
		return
	}
	// DECODE BODY
	var bodyData JWS
	decoder := json.NewDecoder(r.Body)
	decodingErr := decoder.Decode(&bodyData)
	if decodingErr != nil {
		sendErrorResponse(w, "Error decoding input data.")
		return
	}
	decodedJws := []byte(bodyData.JwsObject)
	fileContentInBytes, err := getKeyFileData(keyName, w)
	if err != nil {
		sendErrorResponse(w, "Error fetching Private key data")
		return
	}
	// decode and get private key string
	pemFormattedBlock, _ := pem.Decode(fileContentInBytes)
	privateKey, _ := x509.ParsePKCS1PrivateKey(pemFormattedBlock.Bytes)
	// VERIFY TOKEN:
	_, err = jws.Verify(decodedJws, jwa.RS256, privateKey)
	if err != nil {
		sendErrorResponse(w, "Error: wrong private key")
		return
	}
	err = sendOkResponse(w, true)
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
