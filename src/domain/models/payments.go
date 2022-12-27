package models

import (
	"database/sql"
	"ecomdream/src/pkg/storages/postgres"
	"github.com/sirupsen/logrus"
	"time"
)

type Payment struct {
	ID string `db:"id"`

	SessionID       *string `db:"session_id"`
	PaymentIntentID *string `db:"payment_intent_id"`

	PlanID int  `db:"plan_id"`
	Plan   Plan `db:"plans"`

	VersionID *string  `db:"version_id"`
	Version   *Version `db:"versions"`

	PromocodeID *string `db:"promocode_id"`
	AmountPaid  *int64  `db:"amount_paid"`
	Email       *string `db:"email"`

	CreatedAt time.Time  `db:"created_at"`
	PaidAt    *time.Time `db:"paid_at"`
}

func (p *Payment) Create() (err error) {
	conn := postgres.Connection()

	query := `INSERT INTO payments (id, plan_id, created_at, version_id, promocode_id) VALUES (:id, :plan_id, now(), :version_id, :promocode_id)`

	_, err = conn.NamedExec(query, p)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (p *Payment) MarkAsPaid() (err error) {
	conn := postgres.Connection()

	query := `UPDATE payments SET session_id=$1, payment_intent_id=$2, amount_paid=$3, email=$4, paid_at=now() WHERE id=$5`

	_, err = conn.Exec(query, p.SessionID, p.PaymentIntentID, p.AmountPaid, p.Email, p.ID)
	if err != nil {
		logrus.WithError(err).Errorf("Can't mark payment %s as paid", p.ID)
	}

	return
}

func (p *Payment) SetVersionID(versionID string) (err error) {
	conn := postgres.Connection()
	p.VersionID = &versionID

	query := `UPDATE payments SET version_id=$2 WHERE id=$1`

	_, err = conn.Exec(query, p.ID, p.VersionID)
	if err != nil {
		logrus.WithError(err).Errorf("Can't set version id for payment %s", p.ID)
	}

	return
}

func (p *Payment) Delete() {
	conn := postgres.Connection()

	_, err := conn.Exec("DELETE FROM payments WHERE id=$1", p.ID)
	if err != nil {
		logrus.Errorf("Can't delete payment with id %s", p.ID)
	}
}

func GetPayment(id string) (payment *Payment, err error) {
	conn := postgres.Connection()
	payment = &Payment{}

	query := `SELECT
		plans.id "plans.id", plans.is_init "plans.is_init", plans.price "plans.price", plans.plan_name "plans.plan_name", plans.plan_description "plans.plan_description",
		plans.feature_amount_images "plans.feature_amount_images", plans.feature_amount_image_to_prompt "plans.feature_amount_image_to_prompt",
		plans.is_deprecated "plans.is_deprecated", plans.created_at "plans.created_at",
		payments.*
	FROM payments INNER JOIN plans ON plans.id=payments.plan_id WHERE payments.id=$1`

	err = conn.Get(payment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Can't get payment with id %s", id)
		return
	}

	// TODO left join
	if payment.VersionID != nil {
		payment.Version, err = GetVersion(*payment.VersionID)
	}

	return
}
