package models

import (
	"database/sql"
	"ecomdream/src/pkg/storages/postgres"
	"github.com/sirupsen/logrus"
	"time"
)

type Features struct {
	FeatureAmountImages        int `db:"feature_amount_images" json:"feature_amount_images"`
	FeatureAmountImageToPrompt int `db:"feature_amount_image_to_prompt" json:"feature_amount_image_to_prompt"` // TODO https://replicate.com/methexis-inc/img2prompt
}

type Plan struct {
	ID int `db:"id" json:"id"`

	IsInit bool  `db:"is_init" json:"is_init"`
	Price  int64 `db:"price" json:"price"`

	PlanName string `db:"plan_name" json:"plan_name"`
	PlanDescription string `db:"plan_description" json:"plan_description"`

	IsDeprecated bool `db:"is_deprecated" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`

	Features
}

func GetAvailablePlans() (plans []Plan, err error) {
	conn := postgres.Connection()

	query := `SELECT * FROM plans WHERE is_deprecated=false ORDER BY feature_amount_images ASC`

	err = conn.Select(&plans, query)
	if err != nil {
		logrus.WithError(err).Error("can't get available plans")
	}

	return
}

func GetPlan(id int) (plan *Plan, err error) {
	query := `SELECT * FROM plans WHERE id=$1`
	plan = &Plan{}

	conn := postgres.Connection()

	err = conn.Get(plan, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Can't get plan %d", id)
		return
	}

	return
}
