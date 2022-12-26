create table prompts(
    id varchar unique primary key,

    version_id varchar not null references versions(id),
    prompt_text varchar not null,

    prediction_id varchar not null,
    prediction_time double precision null,

    created_at timestamp with time zone not null,
    finished_at timestamp with time zone null
);
