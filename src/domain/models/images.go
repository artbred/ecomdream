package models

import (
	"ecomdream/src/pkg/storages/postgres"
	"github.com/sirupsen/logrus"
	"time"
)

type Image struct {
	ID       string `db:"id"`
	PromptID string `db:"prompt_id"`

	CdnURL string `db:"cdn_url"`
	Width  int  `db:"width"`
	Height int  `db:"height"`

	IsUpscaled bool      `db:"is_upscaled"`
	CreatedAt  time.Time `db:"created_at"`
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
