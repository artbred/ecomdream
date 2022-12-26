create table payments(
    id varchar unique primary key,

    session_id varchar null,
    payment_intent_id varchar null,

    amount_paid bigint null,
    promocode_id varchar null,

    email varchar null,

    plan_id int not null references plans(id),
    version_id varchar null references versions(id),

    created_at timestamp with time zone not null,
    paid_at timestamp with time zone null
);

create or replace function payment_is_paid (varchar) returns boolean as $$
select exists (
    select 1 from payments where id = $1 and paid_at is not null
);
$$ language sql;



