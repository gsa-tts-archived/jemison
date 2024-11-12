package pages

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func PutData(base_url string, data string) {
	log.Println(base_url, data)

	req, err := http.NewRequest(http.MethodPut,
		base_url,
		bytes.NewBuffer([]byte(data)))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
}

func GetData(base_url string, path string) []byte {
	resp, err := http.Get(base_url + path)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return b
}
