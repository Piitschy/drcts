package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Diff struct {
	Hash string     `json:"hash"`
	Diff Difference `json:"diff"`
}

type Difference struct {
	Collections []map[string]*json.RawMessage `json:"collections"`
	Fields      []map[string]*json.RawMessage `json:"fields"`
	Relations   []map[string]*json.RawMessage `json:"relations"`
}

func (d Directus) GetDiff(s *Snapshot, force bool) (*Diff, error) {
	url := fmt.Sprintf("%s/schema/diff?export=json&access_token=%s", d.Url, d.token)
	if force {
		url = fmt.Sprintf("%s/schema/diff?export=json&force=true&access_token=%s", d.Url, d.token)
	}

	b, err := s.Marshal()
	if err != nil {
		return nil, err
	}
	fmt.Println("diff1:[]byte", string(b))

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 204 {
		return &Diff{}, nil
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get diff with status code %d", res.StatusCode)

	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("diff:string([]byte)", string(bodyBytes))

	var diff Diff
	err = json.Unmarshal(bodyBytes, &diff)
	if err != nil {
		return nil, err
	}

	return &diff, nil
}

func (d Directus) ApplyDiff(diff *Diff) error {
	url := fmt.Sprintf("%s/schema/apply?access_token=%s", d.Url, d.token)
	b, err := json.Marshal(diff)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 || res.StatusCode != 204 {
		return fmt.Errorf("Failed to apply diff with status code %d", res.StatusCode)
	}

	return nil
}
