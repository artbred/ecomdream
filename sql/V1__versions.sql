create table versions (
    id varchar unique primary key not null,

    model_id varchar unique null,
    prediction_id varchar unique not null,

    identifier varchar not null,
    class varchar not null,

    instance_prompt varchar not null,
    class_prompt varchar not null,
    instance_data varchar not null,
    max_train_steps int not null,

    model varchar null,

    created_at timestamp with time zone not null,
    pushed_at timestamp with time zone null
);

create or replace function version_is_pushed (varchar) returns boolean as $$
select exists (
    select 1 from versions where id = $1 and pushed_at is not null
);
$$ language sql;
