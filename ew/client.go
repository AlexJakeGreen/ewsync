package ew

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const endpoint string = "https://outlook.office365.com/EWS/Exchange.asmx"

type EW struct {
	client   *http.Client
	login    string
	password string
}

func NewEW(login, password string) *EW {
	return &EW{
		client:   &http.Client{},
		login:    login,
		password: password,
	}
}

func (c *EW) DoCall(body *bytes.Buffer) ([]byte, error) {
	// fmt.Println("===== Sent ======")
	// fmt.Println(body.String())
	req, err := http.NewRequest("POST", endpoint, body)
	req.SetBasicAuth(c.login, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Println("===== Received ======")
	// fmt.Println(string(res))
	return res, nil
}
