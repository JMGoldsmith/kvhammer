package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

const (
	baseURL = "http://127.0.0.1:8200/v1/kv"
)

type Write struct {
	Data struct {
		Foo string `json:"foo"`
		Zip string `json:"zip"`
	} `json:"data"`
}

type MetaData struct {
	CustomMetadata struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
		Baz string `json:"baz"`
	} `json:"custom_metadata"`
}

func main() {
	// var wg sync.WaitGroup
	var vaultToken string
	var totalRuns int
	var totalRequests int
	flag.StringVar(&vaultToken, "token", "", "token to use for the Vault requests.")
	flag.IntVar(&totalRuns, "runs", 1, "total number of runs for the test.")
	flag.IntVar(&totalRequests, "requests", 10, "Total number of requests to send after writing a KV")
	flag.Parse()
	fmt.Printf("I am a token" + vaultToken)
	runCount := totalRuns
	requestCount := totalRequests
	// grab flags here and then run go funcs according to the flags
	// maybe randomize by default?
	// Launch requests concurrently
	for i := 0; i < runCount; i++ {
		//	wg.Add(1)
		//	go func() {
		//	defer wg.Done()
		secret := makeWriteRequest(vaultToken)
		fmt.Printf("I am a token" + vaultToken)
		makeMetaDataRequest(secret, vaultToken)
		for i := 0; i < requestCount; i++ {
			makeReadRequest(secret, vaultToken)
		}
	}
	//	}

	//wg.Wait()

}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomLength(min, max int) int {
	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	return seededRand.Intn(max-min+1) + min
}

func createKVPayLoad() (string, error) {
	var write Write
	// Generate random lengths between 500 and 1024 bytes
	fooLength := randomLength(500, 1024)
	zipLength := randomLength(500, 1024)

	// Assign random strings of the generated lengths to Foo and Zip
	write.Data.Foo = generateRandomString(fooLength)
	write.Data.Zip = generateRandomString(zipLength)

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(write)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func createMetaPayLoad() (string, error) {
	var metaData MetaData
	// Generate random lengths between 500 and 1024 bytes
	fooLength := randomLength(500, 1024)
	barLength := randomLength(500, 1024)
	bazLength := randomLength(500, 1024)

	// Assign random strings of the generated lengths to Foo and Zip
	metaData.CustomMetadata.Foo = generateRandomString(fooLength)
	metaData.CustomMetadata.Bar = generateRandomString(barLength)
	metaData.CustomMetadata.Baz = generateRandomString(bazLength)

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(metaData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func makeWriteRequest(vaultToken string) string {
	client := &http.Client{}
	secret := generateRandomString(20)
	payload, _ := createKVPayLoad()

	req, err := http.NewRequest("POST", baseURL+"/data/"+secret, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Printf("Failed to create request for %s: %v\n", secret, err)
	}
	req.Header.Add("X-Vault-Token", vaultToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request to %s failed: %v\n", secret, err)
	}
	defer resp.Body.Close()
	return secret
}

// need to get the value of the secret from the call to create
func makeMetaDataRequest(mount string, vaultToken string) {
	client := &http.Client{}
	payload, _ := createMetaPayLoad()

	req, err := http.NewRequest("POST", baseURL+"/metadata/"+mount, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Printf("Failed to create request for %s: %v\n", mount, err)
	}
	req.Header.Add("X-Vault-Token", vaultToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request to %s failed: %v\n", mount, err)
	}
	defer resp.Body.Close()
}

func makeReadRequest(secret string, vaultToken string) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"/data/"+secret, nil)
	if err != nil {
		fmt.Printf("Failed to create request for %s: %v\n", secret, err)
		return
	}
	req.Header.Add("X-Vault-Token", vaultToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request to %s failed: %v\n", secret, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request to %s succeeded with status: %s\n", secret, resp.Status)
}
