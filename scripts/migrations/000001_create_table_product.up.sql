create table if not exists products
(
    id         serial primary key,
    name       text                                not null,
    price      numeric                             not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    created_by varchar(100)                        not null,
    updated_by varchar(100)                        not null
);