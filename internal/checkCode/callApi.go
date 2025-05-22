package checkCode

import (
	"beta-turnstiles/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type ErrorResponse struct {
	StatusText string `json:"status"`
	AppCode    int    `json:"code"`
	//ErrorText  string `json:"error"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"` // sometimes set instead of error for turnstile funcs for compatibility with mairs
}

func callAPIBase(endpoint string, method string, body interface{}, headers map[string]string, resBody interface{}) error {
	if config.GlobalConfig == nil {
		return errors.New("no configuration provided")
	}
	var b io.Reader
	if body != nil {
		jsonData, serr := json.Marshal(body)
		if serr != nil {
			return serr
		}
		b = bytes.NewBuffer(jsonData)
	}

	// send request
	client := http.Client{}
	req, rerr := http.NewRequest(method, config.GetApiEndpoint()+endpoint, b)
	if rerr != nil {
		return rerr
	}
	for k, v := range headers {
		// we do not use Add here, we use Set to avoid duplicates
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO if config.apiKey is set, add to header here
	if config.GlobalConfig.ApiKey != nil {
		req.Header.Set("Authorization", *config.GlobalConfig.ApiKey)
	}

	log.Println("request: ", *req)

	// do the request
	resp, herr := client.Do(req)
	if herr != nil {
		log.Println("got error making request", herr)
		return herr
	}
	log.Println("got response ", resp)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// we have an error
		if resp.StatusCode == http.StatusUnauthorized {
			return errors.New("unauthorized")
		}
		var errBody ErrorResponse
		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)
		decoder := json.NewDecoder(tee)
		decodeerr := decoder.Decode(&errBody)
		if decodeerr != nil {
			return decodeerr
		}

		return errors.New(errBody.Error.Message)
	}

	// get response body
	if resBody != nil {
		decoder := json.NewDecoder(resp.Body)
		decodeerr := decoder.Decode(resBody)
		if decodeerr != nil {
			return decodeerr
		}
	}
	return nil
}
