package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ClientResponse struct {
	Code   uint        `json:"code"`
	Msg    string      `json:"message"`
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func Request(url, method string, data map[string]interface{}, optionalHeader ...map[string]string) *ClientResponse {
	jsonStrBytes, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStrBytes))
	if err != nil {
		return &ClientResponse{Code: 500, Msg: err.Error()}
	}

	req.Header.Set("Content-type", "application/json")
	header := make(map[string]string)
	for _, val := range optionalHeader {
		header = val
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ClientResponse{Code: 500, Msg: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return &ClientResponse{Code: 500, Msg: "http request error :" + strconv.Itoa(resp.StatusCode)}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ClientResponse{Code: 500, Msg: err.Error()}
	}
	result := &ClientResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &ClientResponse{Code: 500, Msg: err.Error()}
	}
	return result
}

type NotifReq struct {
	UserId  uint   `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Apikey  string `json:"api_key"`
}

type Message struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

func SendNotif(user_id uint, title, message, typee, value string) (err error) {
	var req = NotifReq{
		UserId:  user_id,
		Title:   title,
		Message: message,
		Apikey:  os.Getenv("NOTIF_APIKEY"),
	}

	var reqb, _ = json.Marshal(req)
	request, err := http.NewRequest("POST", os.Getenv("NOTIF_URL"), bytes.NewBuffer(reqb))
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func SendRequest(url, token, contentType, method, messageCode string, params []byte) (result []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(params))
	if err != nil {
		return result, err
	}

	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-type", contentType)
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if len(messageCode) > 0 {
		req.Header.Set("message_code", messageCode)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	return result, err
}

func Iam_request(messageCode int, params []byte, token string) (result []byte, err error) {
	req, _ := http.NewRequest("POST", os.Getenv("APP_IAM_URL"), bytes.NewBuffer(params))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("code", strconv.Itoa(messageCode))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	return body, err
}
