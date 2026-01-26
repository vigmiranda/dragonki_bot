package mercadopago

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Client struct {
	accessToken string
}

func NewClient() *Client {
	return &Client{
		accessToken: os.Getenv("MP_ACCESS_TOKEN"),
	}
}

func (c *Client) post(url string, payload any, result any) error {
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}
