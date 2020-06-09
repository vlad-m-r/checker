package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vlad-m-r/checker/api/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type RequestController struct {
	URL           string
	Name          string
	Requests      []models.Request
	errors        []error
	ErrorOccurred bool
	sChan         chan struct{}
	rChan         chan *CheckResult
}

func RequestControllerFactory(check models.Check, sChan chan struct{}, rChan chan *CheckResult) {
	r := RequestController{
		URL:      check.URL,
		Name:     check.Name,
		Requests: check.Requests,
		sChan:    sChan,
		rChan:    rChan,
	}

	log.Println("Processing checks for:", r.Name)

	for _, request := range r.Requests {
		go r.runCheck(request)
	}
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string // Add the request string

	//url := fmt.Sprintf("Method: %v, URL: %v, Proto: %v", r.Method, r.URL, r.Proto)
	request = append(request, "Method: "+r.Method)
	request = append(request, "URL: "+r.URL.String())
	request = append(request, "Proto: "+r.Proto)
	request = append(request, fmt.Sprintf("Host: %v", r.Host)) // Loop through headers
	for name, headers := range r.Header {
		for _, h := range headers {
			request = append(request, fmt.Sprintf("Header - %v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		_ = r.ParseForm()
		request = append(request, "Form: "+r.Form.Encode())
	} // Return the request as a string

	return strings.Join(request, "\n")
}

func (r *RequestController) runCheck(request models.Request) {
	start := time.Now()

	r.sChan <- struct{}{}

	var ioReader io.Reader

	if len(request.Payload) > 0 {
		ioReader = bytes.NewBuffer([]byte(request.Payload))
	} else {
		ioReader = nil
	}

	var response *http.Response
	var httpError error
	var req *http.Request
	var reqError error

	// Create request
	switch request.Method {
	case http.MethodPost:
		req, reqError = http.NewRequest(http.MethodPost, r.URL, ioReader)
		if req != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	case http.MethodPut:
		req, reqError = http.NewRequest(http.MethodPut, r.URL, ioReader)
		if req != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	case http.MethodGet:
		req, reqError = http.NewRequest(http.MethodGet, r.URL, ioReader)
		//response, httpError = http.Get(r.URL)

	default:
		httpError = errors.New("unknown method: " + request.Method)
	}

	if reqError != nil {
		httpError = errors.New("reqError: " + reqError.Error())
	}

	if req == nil {
		httpError = errors.New("req is nil ")
	} else {
		// Add headers
		for _, header := range request.Headers {
			req.Header.Set(header.Name, header.Value)
		}

		log.Println("Request prepared: " + formatRequest(req))

		// Fire request
		response, httpError = http.DefaultClient.Do(req)

		if httpError != nil {
			r.recordError("The HTTP request failed with error: " + httpError.Error())
		}

		if response != nil {
			bodyBytes, err := ioutil.ReadAll(response.Body)

			if err == nil {
				bodyString := string(bodyBytes)
				log.Printf("%s (method %s): %s", r.URL, request.Method, bodyString)
			}

			defer response.Body.Close()

			statusOK := response.StatusCode >= 200 && response.StatusCode < 300

			if !statusOK {
				r.recordError("bad response code:" + response.Status)
			}

			responseBody, readError := ioutil.ReadAll(response.Body)
			if readError != nil {
				r.recordError("Failed to read response body:" + response.Status)
			}

			var data map[string]interface{}
			if unmarshalError := json.Unmarshal(responseBody, &data); unmarshalError != nil {
				r.recordError("Failed to unmarshal output to interface: " + response.Status)
			}

			r.runAsserts(request, data)
		} else {
			r.recordError("got empty response from server")
		}
	}

	seconds := time.Since(start).Seconds()

	// Dump data into result channel
	r.rChan <- &CheckResult{r.Name, r.URL, r.errors, seconds, r.ErrorOccurred}

	// Read from semaphore channel
	<-r.sChan
}

func (r *RequestController) runAsserts(request models.Request, response map[string]interface{}) {
	for _, assert := range request.Asserts {
		switch assert.Type {
		case "keyExists":
			log.Println("Running keyExists assert")
			if _, keyExists := response[assert.Key]; !keyExists {
				r.recordError("failed to find " + assert.Key + " in response: key does not exist")
			}
			log.Println("Running keyExists assert: done")
		}
	}
}

func (r *RequestController) recordError(errorMessage string) {
	r.errors = append(r.errors, errors.New(errorMessage))
	r.ErrorOccurred = true
}
