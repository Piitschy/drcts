package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type authResponse struct {
	Data Auth `json:"data"`
}

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int    `json:"expires"`
}

func (d Directus) GetAuth(email, password string) (*Auth, error) {
	url := fmt.Sprintf("%s/auth/login", d.Url)

	body := map[string]string{
		"email":    email,
		"password": password,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to login with status code %d", res.StatusCode)
	}

	var data authResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}

// Login this directus instance with email and password
func (d *Directus) Login(email, password string) error {
	auth, err := d.GetAuth(email, password)
	if err != nil {
		return err
	}
	if auth.AccessToken == "" {
		return fmt.Errorf("Failed to login: empty token")
	}
	d.SetToken(auth.AccessToken)
	return nil
}
