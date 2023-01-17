CREATE VIEW version_frequent_buyers AS
SELECT
    versions.id as version_id,
    versions.model_id,
    versions.trainer_version,
    versions.prediction_id,
    versions.identifier,
    versions.class,
    versions.instance_prompt,
    versions.class_prompt,
    versions.instance_data,
    versions.max_train_steps,
    versions.model,
    versions.created_at,
    versions.pushed_at,
    versions.deleted_at,
    plans.plan_name,
    plans.plan_description,
    plans.feature_amount_images,
    plans.feature_amount_image_to_prompt,
    payments.amount_paid,
    payments.promocode_id,
    payments.created_at as payment_created_at,
    payments.paid_at as payment_paid_at
FROM
    versions
        LEFT JOIN payments ON versions.id = payments.version_id
        LEFT JOIN plans ON payments.plan_id = plans.id
WHERE
    payments.email IN (
        SELECT email FROM payments
        GROUP BY email
        HAVING COUNT(DISTINCT version_id) > 2
    )
