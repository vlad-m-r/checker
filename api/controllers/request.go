package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vlad-m-r/checker/api/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

	switch request.Method {
	case http.MethodPost:
		response, httpError = http.Post(r.URL, "application/json", ioReader)
	case http.MethodPut:
		req, err := http.NewRequest("PUT", r.URL, ioReader)

		if err != nil {
			r.recordError("The HTTP request failed with error: " + err.Error())
		}

		if req != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		response, httpError = http.DefaultClient.Do(req)
	case http.MethodGet:
		response, httpError = http.Get(r.URL)
	default:
		httpError = errors.New("unknown method: " + request.Method)
	}

	if httpError != nil {
		r.recordError("The HTTP request failed with error: " + httpError.Error())
	}

	if response != nil {
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
			r.recordError("Failed to unmarshal output to interface:" + response.Status)
		}

		r.runAsserts(request, data)
	} else {
		r.recordError("got empty response from server")
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
