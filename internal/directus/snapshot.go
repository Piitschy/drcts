package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Snapshot struct {
	Version     *int                          `json:"version"`
	Directus    *string                       `json:"directus"`
	Vendor      *string                       `json:"vendor"`
	Collections []map[string]*json.RawMessage `json:"collections"`
	Fields      []map[string]*json.RawMessage `json:"fields"`
	Relations   []map[string]*json.RawMessage `json:"relations"`
}

func (s *Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// FilterCollections filters the snapshot to only include the collections specified in the argument
// If a collection is not found in the snapshot, it is ignored
func (s *Snapshot) FilterCollections(c ...string) error {
	var newCollections []map[string]*json.RawMessage
	for _, collection := range s.Collections {
		for _, name := range c {
			if _, ok := collection[name]; ok {
				newCollections = append(newCollections, collection)
				break
			}
		}
	}
	s.Collections = newCollections
	return nil
}

func (d *Directus) GetSnapshot() (*Snapshot, error) {
	bodyBytes, err := d.GetRawSnapshot("json")

	if err != nil {
		return nil, err
	}

	// fmt.Println("string([]byte)", string(bodyBytes))

	var s Snapshot
	err = json.Unmarshal(bodyBytes, &s)
	// fmt.Println("s", s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *Directus) GetRawSnapshot(format string) ([]byte, error) {
	if format != "json" && format != "yaml" && format != "xml" && format != "csv" {
		return nil, fmt.Errorf("Invalid format %s. Use json, yaml, xml or csv", format)
	}
	url := fmt.Sprintf("%s/schema/snapshot?export=%s&access_token=%s", d.Url, format, d.token)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		return nil, fmt.Errorf("Failed to get snapshot. No permissions! Check your token...")
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get snapshot with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
