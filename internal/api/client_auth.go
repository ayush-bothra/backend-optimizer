package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func RegisterUser(email, password string) (map[string]interface{}, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/dbconnections/signup"
	
	body := map[string]string {
		"client_id": os.Getenv("CLIENT1_ID"),
		"email": email,
		"password": password,
		"connection": "Username-Password-Authentication", // auth_DB
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}


func LoginUser(email, password string) (map[string]interface{}, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	body := map[string]string{
		"grant_type": "password",
        "username":   email,
        "password":   password,
		"client_id": os.Getenv("CLIENT1_ID"),
		"client_secret": os.Getenv(("CLIENT1_SECRET")),
		"audience": os.Getenv("AUTH0_AUDIENCE"),
		"realm": "Username-Password-Authentication",
	}

	// Marshal returns the JSON encoding of v.
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	// Set sets the header entries associated with key to the single element value. 
	req.Header.Set("Content-type", "application/json")

	// Do sends an HTTP request and returns an HTTP response, 
	// following policy (such as redirects, cookies, auth) as configured on the client.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err 
	}
	// close the response after completion
	defer res.Body.Close()

	var tokenResp map[string]interface{}
	
	// Decode reads the next JSON-encoded value from its
	// input and stores it in the value pointed to by v.
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		return nil, err 
	}

	return tokenResp, nil
}