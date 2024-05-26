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

func (d *Diff) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d Directus) GetDiff(s *Snapshot, force bool) (*Diff, error) {
	bodyBytes, err := d.GetRawDiff(s, "json", force)
	if err != nil {
		return nil, err
	}

	if bodyBytes == nil {
		return nil, nil
	}

	var diff Diff
	err = json.Unmarshal(bodyBytes, &diff)
	if err != nil {
		return nil, err
	}

	if len(diff.Diff.Collections) == 0 && len(diff.Diff.Fields) == 0 && len(diff.Diff.Relations) == 0 {
		return nil, nil
	}

	return &diff, nil
}

func (d Directus) GetRawDiff(s *Snapshot, format string, force bool) ([]byte, error) {
	if format != "json" && format != "yaml" && format != "xml" && format != "csv" {
		return nil, fmt.Errorf("Invalid format %s. Use json, yaml, xml or csv", format)
	}
	url := fmt.Sprintf("%s/schema/diff?export=%s&force=true&access_token=%s", d.Url, format, d.token)
	if !force {
		url = fmt.Sprintf("%s/schema/diff?export=%s&access_token=%s", d.Url, format, d.token)
	}

	b, err := s.Marshal()
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 204 {
		return nil, nil
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get diff with status code %d", res.StatusCode)
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
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

	if res.StatusCode != 200 && res.StatusCode != 204 {
		return fmt.Errorf("Failed to apply diff with status code %d", res.StatusCode)
	}

	return nil
}
