package checkCode

import (
	"beta-turnstiles/internal/config"
	"errors"
	"log"
	"strings"
)

type initTurnstileResponseData struct {
	GateName string `json:"gateName"`
	SiteName string `json:"siteName"`
	//SwipeInterval int `json:"swipeInterval"` // in s
	CtrlDir int `json:"ctrlDir"` // 0 only in, 1 only out, 2 bidirectional
	//WorkMode int `json:"workMode"` // 0 regular. 1 pass ct. 2 always on. 3 handheld, identiy sampling
	NetworkMode int `json:"networkMode"` // 0 online 1 offline, 2 both
	//AccessType int `json:"accessType"` // 0 regular, 1 scene channel, 2 regular+sessions
	//SystemDate time.Time `json:"systemDate"`

}

type turnstileReponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type initTunrstileResponse struct {
	Data  initTurnstileResponseData `json:"data"`
	Error turnstileReponseError     `json:"error"`
}

// check in with BETA to see if this IP address is ok
func CheckIn() (string, error) {
	var res initTunrstileResponse
	err := callAPIBase("checkIn", "GET", nil, nil, &res)
	if err != nil {
		return "", err
	}
	return res.Data.SiteName, nil
}

type checkAvailableRequestBody struct {
	UniqueCode string `json:"uniqueCode"`
	AccessDir  int    `json:"accessDir"` // 0 in, 1 out
}

type checkAvailableResponseBodyData struct {
	GUID         string `json:"guid"`
	AccessDir    int    `json:"accessDir"`    //0 in, 1/2 out
	SerialNumber string `json:"serialNumber"` // tiket id
}
type checkAvailableResponseBody struct {
	Data  checkAvailableResponseBodyData `json:"data"`
	Error turnstileReponseError          `json:"error"`
}

func CheckCode(code string) error {
	body := checkAvailableRequestBody{
		UniqueCode: code,
		AccessDir:  0,
	}
	if config.GlobalConfig.DirectionOut {
		body.AccessDir = 1
	}

	var res checkAvailableResponseBody
	err := callAPIBase("checkAvailable", "POST", body, nil, &res)
	if err != nil {
		return err
	}
	if res.Error.Code != 0 {
		return errors.New(res.Error.Message)
	}
	return nil
}

type doConsumeRequestBody struct {
	GUID string `json:"guid"` // needs to include the code and the access direction
}

type doConsumeResponseBodyData struct {
	Result        bool   `json:"result"`
	GUID          string `json:"guid"`
	VerboseResult string `json:"verboseResult"`
}
type doConsumeResponseBody struct {
	Data  doConsumeResponseBodyData `json:"data"`
	Error turnstileReponseError     `json:"error"`
}

func ConsumeCode(code string) (string, error) {
	body := doConsumeRequestBody{
		GUID: code,
	}
	if !strings.Contains(body.GUID, ":") {
		if config.GlobalConfig.DirectionOut {
			body.GUID += ":1"
		} else {
			body.GUID += ":0"
		}
	}
	log.Printf("sending body %v\n", body)
	var res doConsumeResponseBody
	err := callAPIBase("doConsume", "POST", body, nil, &res)
	if err != nil {
		return "", err
	}
	if res.Error.Code != 0 {
		return "", errors.New(res.Error.Message)
	}
	return res.Data.VerboseResult, nil
}
