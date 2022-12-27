create table prompts(
    id varchar unique primary key,

    version_id varchar not null references versions(id),

    prompt_text varchar not null,
    prompt_negative varchar null,
    amount_images int not null,
    width int not null,
    height int not null,
    prompt_strength double precision not null,
    inference_steps int not null,
    guidance_scale double precision not null,
    seed bigint null,

    prediction_id varchar not null,
    prediction_time double precision null,

    created_at timestamp with time zone not null,
    finished_at timestamp with time zone null
);
