package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

const port = "8080"
const subscriptionKeyEnvVar = "COGNITIVE_TRANSLATOR_TEXT_SUBSCRIPTION_KEY"
const endpointEnvVar = "COGNITIVE_TRANSLATOR_API_BASE_ENDPOINT"
const endpointSuffix = "/translate?api-version=3.0"
const subscriptionKeyHeader = "Ocp-Apim-Subscription-Key"
const translateRoute = "/translate-text"

var subscriptionKey string
var baseEndpoint string

func init() {
	log.Println("Checking env vars")
	subscriptionKey = os.Getenv(subscriptionKeyEnvVar)

	if "" == subscriptionKey {
		log.Fatalf("Please set the %s environment variable", subscriptionKeyEnvVar)
	}
	baseEndpoint = os.Getenv(endpointEnvVar)

	if "" == baseEndpoint {
		log.Fatalf("Please set the %s environment variable", endpointEnvVar)
	}
}

func main() {
	//front-end
	uiHandler := http.FileServer(http.Dir("ui"))
	http.Handle("/", uiHandler)

	//backend API
	http.HandleFunc(translateRoute, func(w http.ResponseWriter, r *http.Request) {

		//JSON request from front-end
		type request struct {
			Text   string `json:"text"`
			ToLang string `json:"to"`
		}

		var userReq request
		err := json.NewDecoder(r.Body).Decode(&userReq)

		if err != nil {
			log.Println("Unable to decode translate request from user", err)
		}

		log.Printf("Translate text '%s' to language '%s'\n", userReq.Text, userReq.ToLang)

		response, err := translateText(userReq.Text, userReq.ToLang)
		if err != nil {
			log.Println("Failed to translate", err)
		}
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Println("Unable to return translation response back to user", err)
		}
	})
	log.Println("Starting translator service...")

	http.ListenAndServe(":"+port, nil)
}

func translateText(text, toLang string) (translateAPIResponse, error) {

	//append the target language
	cognitiveServiceEndpoint := baseEndpoint + endpointSuffix + "&to=" + toLang

	//keep it simple
	reqBody := `[{'Text':'` + text + `'}]`
	log.Println("reqBody", reqBody)

	//POST request
	req, err := http.NewRequest(http.MethodPost, cognitiveServiceEndpoint, strings.NewReader(reqBody))

	if err != nil {
		log.Println("Error creating HTTP POST request", err)
		return nil, err
	}
	//add headers
	req.Header.Add(subscriptionKeyHeader, subscriptionKey)
	req.Header.Add("Content-Type", "application/json")

	//Invoke API
	log.Println("Invoking endpoint", cognitiveServiceEndpoint)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error invoking Cognitive API", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Println("Cognitive API returned HTTP", res.StatusCode)
		return nil, errors.New("Cognitive API returned HTTP " + res.Status)
	}

	//decode response to a struct
	var result translateAPIResponse
	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		log.Println("Error decoding Cognitive API response", err)
		return nil, err

	}
	return result, nil
}

type translateAPIResponse []struct {
	DetectedLanguage struct {
		Language string  `json:"language"`
		Score    float64 `json:"score"`
	} `json:"detectedLanguage"`
	Translations []struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
}
