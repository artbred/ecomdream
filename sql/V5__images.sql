create table images(
    id varchar unique primary key,

    prompt_id varchar not null references prompts(id),
    cdn_url varchar not null,

    width int not null,
    height int not null,

    is_upscaled boolean not null default false,
    created_at timestamp with time zone not null
);
