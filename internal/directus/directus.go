package directus

import (
	"fmt"
	"net/http"
)

type Directus struct {
	Url   string
	token string
}

func NewDirectus(url, token string) (*Directus, error) {
	d := &Directus{
		Url:   url,
		token: token,
	}
	// err := d.TestConnection()
	// if err != nil {
	// 	return nil, err
	// }
	return d, nil
}

func (d Directus) TestConnection() error {
	res, err := http.Get(d.Url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Connection failed with status code %d", res.StatusCode)
	}
	return nil
}

func (d *Directus) SetToken(token string) {
	d.token = token
}
