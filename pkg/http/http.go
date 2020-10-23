package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

type Message interface{}

type Service interface {
	Address() string
}

type MessageEncodeUrl interface {
	Encode() (urlQuery string, err error)
}

type ServiceEndpoint struct {
	Service Service
	Endpoint string
	Method string
	ContentType string
}

func (e *ServiceEndpoint) Send(message Message) (response Message, err error) {
	var httpResponse *http.Response

	switch e.Method {
	case "GET":
		httpResponse, err = e.sendGet(message)
	case "POST", "PUT", "DELETE":
		httpResponse, err = e.sendWithBody(message)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response: %d", httpResponse.StatusCode)
		return
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("cannot unmarshal request")
	}
	return
}

func (e *ServiceEndpoint) sendWithBody(request Message) (response *http.Response, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}

	var queryParams string
	if urlEncode, ok := request.(MessageEncodeUrl); ok {
		queryParams, err = urlEncode.Encode()
		if err != nil {
			return
		}
	}

	client := http.Client{}
	req, err := http.NewRequest(e.Method, path.Join(e.getFullAddress(), queryParams), bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", e.ContentType)

	response, err = client.Do(req)
	return
}

func (e *ServiceEndpoint) sendGet(message Message) (response *http.Response, err error) {
	response, err = http.Get(e.getFullAddress())
	return
}

func (e *ServiceEndpoint) getFullAddress() string {
	return fmt.Sprintf("%s%s", e.Service.Address(), e.Endpoint)
}