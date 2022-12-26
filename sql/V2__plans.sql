create table plans (
   id serial unique primary key not null,
   is_init boolean not null,

   price int not null,
   plan_name varchar not null,
   plan_description varchar not null,

   feature_amount_images bigint not null,
   feature_amount_image_to_prompt bigint not null,

   is_deprecated boolean not null default false,
   created_at timestamp with time zone not null default now()
);

INSERT INTO plans (is_init, price, plan_name, plan_description, feature_amount_images, feature_amount_image_to_prompt) VALUES (
    true, 1300, 'basic_init', 'efewfew', 100, 0
);



