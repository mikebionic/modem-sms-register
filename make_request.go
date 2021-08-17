package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func make_request(url_address string, phone_number string, message_text string) error {
	postBody, _ := json.Marshal(map[string]string{
		"phone_number": phone_number,
		"message_text": message_text,
	})
	responseBody := bytes.NewBuffer(postBody)
	fmt.Println(responseBody)

	client := http.Client{}

	req, err := http.NewRequest("POST", url_address, responseBody)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer Token"},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
