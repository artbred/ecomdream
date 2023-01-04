package models

import (
	"database/sql/driver"
	"ecomdream/src/pkg/storages/postgres"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type Image struct {
	ID       string `db:"id" json:"-"`
	PromptID string `db:"prompt_id" json:"-"`

	CdnURL string `db:"cdn_url" json:"cdn_url"`
	Width  int  `db:"width" json:"width"`
	Height int  `db:"height" json:"height"`

	IsUpscaled bool      `db:"is_upscaled" json:"is_upscaled"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type PromptImages []Image

func (t *PromptImages) Value() (driver.Value, error) {
	b, err := json.Marshal(t); if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(b)), nil
}

func (t *PromptImages) Scan (src interface{}) error {
	var source []byte

	if reflect.TypeOf(src) == nil {
		return nil
	}

	switch src.(type) {
		case string:
			source = []byte(src.(string))
		case []byte:
			source = src.([]byte)
		default:
			return errors.New("incompatible type for PromptImages")
	}

	return json.Unmarshal(source, t)
}

func (i *Image) Create() (err error) {
	conn := postgres.Connection()

	query := `INSERT INTO images (id, prompt_id, cdn_url, width, height, is_upscaled, created_at) VALUES
             (:id, :prompt_id, :cdn_url, :width, :height, false, now())`

	_, err = conn.NamedExec(query, i)
	if err != nil {
		logrus.WithError(err).Errorf("can't add image to db, %+v", i)
	}

	return
}
