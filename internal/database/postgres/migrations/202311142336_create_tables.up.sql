begin transaction;

create type order_status as enum ('new', 'registered', 'processing', 'processed', 'invalid');

create table if not exists users (
    user_id bigserial primary key,
    login varchar(255) unique not null,
    password varchar(255) not null,
    current integer default 0,
    withdrawn integer default 0
);

create index if not exists users_login_idx on users(user_id, login);

create table if not exists orders (
    number varchar(255) primary key,
    user_id bigint not null references users(user_id),
    status order_status not null,
    accrual integer default 0,
    created_at timestamp not null
);

create index if not exists orders_number_idx on orders(number);

create table if not exists withdrawals (
    user_id bigint not null references users(user_id),
    order_number varchar(255) not null,
    sum integer not null,
    updated_at timestamp not null
);

commit;