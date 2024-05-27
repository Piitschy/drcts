package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshalCollection(data []byte) (*Collection, error) {
	var r Collection
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *Collection) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CollectionResponse struct {
	Data Collection `json:"data"`
}

type Collection struct {
	Collection string            `json:"collection"`
	Meta       CollectionMeta    `json:"meta"`
	Schema     CollectionSchema  `json:"schema"`
	Fields     []CollectionField `json:"fields,omitempty"`
}

type CollectionField struct {
	Field  string      `json:"field"`
	Type   string      `json:"type"`
	Meta   FieldMeta   `json:"meta"`
	Schema FieldSchema `json:"schema"`
}

type FieldMeta struct {
	Icon string `json:"icon"`
}

type FieldSchema struct {
	IsPrimaryKey bool `json:"is_primary_key"`
	IsNullable   bool `json:"is_nullable"`
}

type CollectionMeta struct {
	Collection            string        `json:"collection"`
	Icon                  *string       `json:"icon"`
	Note                  *string       `json:"note"`
	DisplayTemplate       *string       `json:"display_template"`
	Hidden                bool          `json:"hidden"`
	Singleton             bool          `json:"singleton"`
	Translations          []Translation `json:"translations"`
	ArchiveField          *string       `json:"archive_field"`
	ArchiveValue          string        `json:"archive_value"`
	UnarchiveValue        string        `json:"unarchive_value"`
	ArchiveAppFilter      bool          `json:"archive_app_filter"`
	SortField             *string       `json:"sort_field"`
	ItemDuplicationFields interface{}   `json:"item_duplication_fields"`
	Sort                  *int64        `json:"sort"`
	Accountability        *string       `json:"accountability,omitempty"`
	Color                 interface{}   `json:"color"`
	Group                 interface{}   `json:"group"`
	Collapse              *string       `json:"collapse,omitempty"`
	PreviewURL            interface{}   `json:"preview_url"`
	Versioning            *bool         `json:"versioning,omitempty"`
}

type Translation struct {
	Language    string  `json:"language"`
	Translation string  `json:"translation"`
	Singular    *string `json:"singular,omitempty"`
	Plural      *string `json:"plural,omitempty"`
}

type CollectionSchema struct {
	Name    string      `json:"name"`
	Comment interface{} `json:"comment"`
	SQL     *string     `json:"sql,omitempty"`
}

func (d *Directus) CreateCollection(c *Collection) error {
	bodyBytes, err := c.Marshal()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/collections?access_token=%s", d.Url, d.token)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	if res.StatusCode == 403 {
		return fmt.Errorf("Failed to create collection: %s - Try to login into directus first with directus.Login()", res.Status)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to create collection: %s", res.Status)
	}

	return nil
}

func (d *Directus) GetCollection(collection string) (*Collection, error) {
	url := fmt.Sprintf("%s/collections/%s?access_token=%s", d.Url, collection, d.token)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get collection: %s", res.Status)
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cr CollectionResponse
	err = json.Unmarshal(bodyBytes, &cr)
	if err != nil {
		return nil, err
	}
	return &cr.Data, nil
}
