package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"log"
	"github.com/joho/godotenv"
	"strings"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	payload := strings.NewReader(`{
		"client_id":"KDdKFRqwVgBSJNwJtzy9r1dKflKCeRdV",
		"client_secret":"CU6yPm_BUcFhF-HleBPYjCjIpT9IOXwpBTtX_QLarB7ERqFjKw6OS7t58t1B53cX",
		"audience":"http://localhost/8080/",
		"grant_type":"client_credentials"
	}`)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(res.Status)
	fmt.Println(string(body))
}
