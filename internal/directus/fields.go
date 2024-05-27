package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func UnmarshalField(data []byte) (*Field, error) {
	var r Field
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *Field) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Field struct {
	Collection string `json:"collection"`
	Field      string `json:"field"`
	Type       string `json:"type"`
	Meta       Meta   `json:"meta"`
	Schema     Schema `json:"schema"`
}

type Meta struct {
	ID             int64       `json:"id"`
	Collection     string      `json:"collection"`
	Field          string      `json:"field"`
	Special        interface{} `json:"special"`
	Interface      string      `json:"interface"`
	Options        interface{} `json:"options"`
	Display        interface{} `json:"display"`
	DisplayOptions interface{} `json:"display_options"`
	Readonly       bool        `json:"readonly"`
	Hidden         bool        `json:"hidden"`
	Sort           int64       `json:"sort"`
	Width          string      `json:"width"`
	Translations   interface{} `json:"translations"`
	Note           string      `json:"note"`
}

type Schema struct {
	Name             string      `json:"name"`
	Table            string      `json:"table"`
	DataType         string      `json:"data_type"`
	DefaultValue     interface{} `json:"default_value"`
	MaxLength        interface{} `json:"max_length"`
	NumericPrecision int64       `json:"numeric_precision"`
	NumericScale     int64       `json:"numeric_scale"`
	IsNullable       bool        `json:"is_nullable"`
	IsPrimaryKey     bool        `json:"is_primary_key"`
	HasAutoIncrement bool        `json:"has_auto_increment"`
	ForeignKeyColumn interface{} `json:"foreign_key_column"`
	ForeignKeyTable  interface{} `json:"foreign_key_table"`
	Comment          interface{} `json:"comment"`
}

func (d *Directus) CreateField(field Field) error {
	b, err := field.Marshal()
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("%s/fields", d.Url), "application/json", bytes.NewBuffer(b))

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to create field: %s", res.Status)
	}
	return nil
}

func (d *Directus) GetFieldOfCollection(collection string) (*Field, error) {
	url := fmt.Sprintf("%s/fields/%s", d.Url, collection)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get field: %s", res.Status)
	}

	var f Field
	err = json.NewDecoder(res.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
