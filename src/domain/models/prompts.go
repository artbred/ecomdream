package models

import (
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/postgres"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type Prompt struct {
	ID string `db:"id" json:"-"`
	VersionID string `db:"version_id" json:"-"`

	PromptText string `db:"prompt_text" json:"prompt"`
	PromptNegative *string `db:"prompt_negative" json:"negative_prompt"`
	AmountImages int `db:"amount_images" json:"num_outputs"`
	Width int `db:"width" json:"-"`
	Height int `db:"height" json:"-"`
	PromptStrength float64 `db:"prompt_strength" json:"-"`
	InferenceSteps int `db:"inference_steps" json:"-"`
	GuidanceScale float64 `db:"guidance_scale" json:"-"`
	Seed int `db:"seed" json:"-"`

	PredictionID string `db:"prediction_id" json:"-"`
	PredictionTime *float64 `db:"prediction_time" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	FinishedAt *time.Time `db:"finished_at" json:"-"`

	PromptImages PromptImages `db:"images" json:"images"`
}

func (p *Prompt) TransferToReplicateBody(version *Version) (body []byte, err error) {
	req := replicate.Request{
		Version: *version.ModelID,
		Input: nil,
	}

	bodyInput, err := json.Marshal(p); if err != nil {
		logrus.WithError(err).Error("can't marshal prompt")
		return
	}

	err = json.Unmarshal(bodyInput, &req.Input); if err != nil {
		logrus.WithError(err).Error("can't unmarshal prompt")
	}

	body, err = json.Marshal(req); if err != nil {
		logrus.WithError(err).Error("can't marshal prompt to final bytes")
	}

	return
}

func (p *Prompt) MarkAsFinished() (err error) {
	conn := postgres.Connection()

	_, err = conn.Exec("UPDATE prompts SET finished_at=now(), prediction_time=$1 WHERE id=$2", p.PredictionTime, p.ID)
	if err != nil {
		logrus.WithError(err).Errorf("Can't mark prompt %s as finished", p.ID)
	}

	return
}

func (p *Prompt) Create() (err error) {
	conn := postgres.Connection()

	query := `INSERT INTO prompts (id, version_id, prompt_text, prompt_negative, amount_images, width, height,
			  prompt_strength, inference_steps, guidance_scale, prediction_id, created_at) VALUES
			  (:id, :version_id, :prompt_text, :prompt_negative, :amount_images, :width, :height,
			  :prompt_strength, :inference_steps, :guidance_scale, :prediction_id, now())`

	_, err = conn.NamedExec(query, p)
	if err != nil {
		logrus.WithError(err).Errorf("Can't create prompt %+v", p)
	}

	return
}

func GetCompletedPromptsForVersion(versionID string) (prompts []Prompt, err error) {
	conn := postgres.Connection()

	query := `SELECT *, (
		SELECT json_agg(i) AS images FROM (
			SELECT * FROM images GROUP BY images.id HAVING bool_or(prompt_id=prompts.id)
		) i
	) FROM prompts WHERE version_id=$1 AND finished_at IS NOT NULL ORDER BY prompts.created_at DESC`

	err = conn.Select(&prompts, query, versionID); if err != nil {
		logrus.WithError(err).Errorf("Can't get prompts for version %s", versionID)
	}

	return
}

func GetRunningPrompts() (prompts []*Prompt, err error) {
	conn := postgres.Connection()

	query := `SELECT * FROM prompts WHERE finished_at IS NULL`

	err = conn.Select(&prompts, query)
	if err != nil {
		logrus.WithError(err).Error("Can't get running prompts")
	}

	return
}
