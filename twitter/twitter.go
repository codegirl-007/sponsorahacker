package twitter

import (
	"encoding/json"
	"net/http"
)

type TwitterClient interface {
	GetUsername(userId string) (string, error)
	GetFollowing(userId string) ([]User, error)
}

type Client struct {
	httpClient  *http.Client
	bearerToken string
	baseURL     string
}

type User struct {
	ID       string
	Username string
	Name     string
}

func NewClient(bearerToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		bearerToken: bearerToken,
		baseURL:     "https://api.twitter.com",
	}
}

func (c *Client) GetUsername(userId string) (string, error) {
	url := c.baseURL + "/users/" + userId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var userResponse struct {
		Data User
	}

	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return "", err
	}

	return userResponse.Data.Username, nil
}

func (c *Client) GetFollowing(userId string) ([]User, error) {
	url := c.baseURL + "/users/" + userId + "/following"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	var followingResponse struct {
		Data []User
	}

	if err := json.NewDecoder(resp.Body).Decode(&followingResponse); err != nil {
		return nil, err
	}

	return followingResponse.Data, nil
}
