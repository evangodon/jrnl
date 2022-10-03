package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/evangodon/jrnl/internal/cfg"
)

type Res struct {
	Status int
	Body   []byte
}

func MakeRequest(method string, path string, bodyParams io.Reader) (Res, error) {
	apiConfig := cfg.GetConfig().API
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNotFound {
		fmt.Println(string(body))
		panic("request error")
	}
	res.Body.Close()

	r.Body = body

	return r, nil
}
