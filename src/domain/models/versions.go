package models

import (
	"database/sql"
	"ecomdream/src/pkg/storages/postgres"
	"github.com/sirupsen/logrus"
	"time"
)

type Version struct {
	ID string `db:"id"`

	ModelID      *string `db:"model_id"`
	PredictionID string  `db:"prediction_id"`

	Identifier string `db:"identifier"`
	Class      string `db:"class"`

	InstancePrompt string `db:"instance_prompt"`
	ClassPrompt    string `db:"class_prompt"`
	InstanceData   string `db:"instance_data"`

	MaxTrainStep int64  `db:"max_train_steps"`
	Model        string `db:"model"`

	CreatedAt time.Time  `db:"created_at"`
	PushedAt  *time.Time `db:"pushed_at"`

	AmountImagesGenerated int `db:"amount_images_generated"`
}

func (v *Version) Create(payment *Payment) (err error) {
	conn := postgres.Connection()

	tx := conn.MustBegin()
	defer func() {
		err = tx.Commit()
		if err != nil {
			logrus.Error(err)
		}
	}()

	queryVersion := `INSERT INTO versions (id, prediction_id, identifier, class, instance_prompt, class_prompt, instance_data, max_train_steps, model, created_at) VALUES (
    	:id, :prediction_id, :identifier, :class, :instance_prompt, :class_prompt, :instance_data, :max_train_steps, :model, now())`

	_, err = tx.NamedExec(queryVersion, v)
	if err != nil {
		logrus.Error(err)
		return
	}

	payment.VersionID = &v.ID
	queryPayment := `UPDATE payments SET version_id=$2 WHERE id=$1`

	_, err = tx.Exec(queryPayment, payment.ID, payment.VersionID)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func GetVersion(id string) (version *Version, err error) {
	conn := postgres.Connection()
	version = &Version{}

	query := `SELECT
    	versions.*,
       	count(images.id) as amount_images_generated FROM versions
		INNER JOIN prompts ON prompts.id=versions.id
		INNER JOIN images ON images.prompt_id=prompts.id
		WHERE versions.id=$1 GROUP BY versions.id`

	err = conn.Get(version, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Can't get version with id %s", id)
		return
	}

	return
}

func GetRunningVersions() (versions []Version, err error) {
	conn := postgres.Connection()

	err = conn.Select(&versions, "SELECT * FROM versions WHERE pushed_at IS NULL")
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (v *Version) MarkAsPushed(modelID *string) (err error) {
	conn := postgres.Connection()
	v.ModelID = modelID

	_, err = conn.Exec("UPDATE versions SET pushed_at=now(), model_id=$1 WHERE id=$2", v.ModelID, v.ID)
	if err != nil {
		logrus.WithError(err).Errorf("Can't mark version %s as pushed", v.ID)
	}

	return
}

func (v *Version) GetUnifiedFeatures() (features Features, err error) {
	conn := postgres.Connection()

	query := `SELECT
		SUM(feature_amount_images) as feature_amount_images,
		bool_or(feature_prompts_database) as feature_prompts_database,
		SUM(feature_upscale_limit) as feature_upscale_limit
		FROM payments INNER JOIN plans ON payments.plan_id = plans.id WHERE payments.version_id=$1 AND payments.paid_at IS NOT NULL`

	err = conn.Get(&features, query, v.ID)
	if err != nil {
		logrus.WithError(err).Errorf("Can't get unified features for version %s", v.ID)
	}

	return
}
