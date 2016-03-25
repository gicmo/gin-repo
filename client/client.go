package client

import (
	"net/http"
	"net/url"

	"encoding/json"
	"fmt"
	. "github.com/G-Node/gin-repo/common"
	"io/ioutil"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (client *Client) LookupUserByFingerprint(fingerprint string) (*User, error) {
	web := &http.Client{}

	params := url.Values{}
	params.Add("key", fingerprint)
	url := fmt.Sprintf("%s/intern/user/lookup?%s", client.address, params.Encode())

	res, err := web.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}