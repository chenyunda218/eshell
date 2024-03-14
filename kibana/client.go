package kibana

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewClient(username, password, baseURL string) *Client {
	return &Client{
		Username: username,
		Password: password,
		BaseURL:  baseURL,
		cookies:  make(map[string]string),
	}
}

type Client struct {
	Username string
	Password string
	BaseURL  string
	cookies  map[string]string
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	ProviderType string      `json:"providerType"`
	ProviderName string      `json:"providerName"`
	CurrentURL   string      `json:"currentURL"`
	Params       loginParams `json:"params"`
}

func (self *Client) Login() {
	payload := loginRequest{
		ProviderType: "basic",
		ProviderName: "basic",
		CurrentURL:   "/",
		Params: loginParams{
			Username: self.Username,
			Password: self.Password,
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", self.BaseURL+"/internal/security/login", bytes.NewBuffer(body))
	req.Header.Add("kbn-xsrf", "reporting")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	cookies := res.Cookies()
	for _, cookie := range cookies {
		self.cookies[cookie.Name] = cookie.Value
	}
}

type GetDataViewsResponse struct {
	DataViews []DataView `json:"data_view"`
}

func (self *Client) DataViews() []DataView {
	req, _ := http.NewRequest("GET", self.BaseURL+"/api/data_views", nil)
	req.Header.Add("kbn-xsrf", "reporting")
	req.Header.Set("Content-Type", "application/json")
	for k, v := range self.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	client := http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	var dataViewsResponse GetDataViewsResponse
	json.Unmarshal(body, &dataViewsResponse)
	return dataViewsResponse.DataViews
}
