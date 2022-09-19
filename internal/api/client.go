package api

import (
	"io"
	"io/ioutil"
	"net/http"
)

type Res struct {
	Status int
	Body   []byte
}

func MakeRequest(method string, path string, bodyParams io.Reader) (Res, error) {
	// TODO: get base url from config file
	baseURL := "http://localhost:8080"
	url := baseURL + path
	req, _ := http.NewRequest(method, url, bodyParams)

	r := Res{}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		r.Status = res.StatusCode
		return r, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		r.Status = res.StatusCode
		return r, err
	}
	res.Body.Close()

	r.Body = body

	// TODO: check status here

	return r, nil
}
