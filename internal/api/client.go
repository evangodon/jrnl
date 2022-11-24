package api

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/evangodon/jrnl/internal/cfg"
)

type Res struct {
	Status int
	Body   []byte
}

type Client struct {
	Config cfg.Config
}

func (c *Client) MakeRequest(method string, path string, bodyParams io.Reader) (Res, error) {
	apiConfig := c.Config.API
	url := apiConfig.BaseURL + path
	req, err := http.NewRequest(method, url, bodyParams)
	if err != nil {
		panic("Failed to setup http request")
	}

	req.Header.Add("X-API-Key", apiConfig.Key)

	r := Res{}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("jrnl server not running")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		r.Status = res.StatusCode
		return r, err
	}

	if res.StatusCode >= 300 {
		panic("request error: " + string(body))
	}
	res.Body.Close()

	r.Status = res.StatusCode
	r.Body = body

	return r, nil
}
