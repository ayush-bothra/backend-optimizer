package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	body := map[string]string{
		"client_id":     os.Getenv("CLIENT0_ID"),
		"client_secret": os.Getenv("CLIENT0_SECRET"),
		"audience":      "http://localhost:8080/",
		"grant_type":    "client_credentials",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	payload := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	newbody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(res.Status)
	fmt.Println(string(newbody))
}
