begin transaction;

create table if not exists users (
    user_id uuid default gen_random_uuid() primary key,
    login varchar(255) not null,
    password varchar(255) not null
);

create table if not exists orders (
    number varchar(255) primary key,
    user_id uuid not null references users(user_id),
    status varchar(255) not null,
    accrual integer,
    uploaded_at timestamp not null
);

create table if not exists user_accounts (
    user_id uuid not null references users(user_id),
    current integer default 0,
    withdrawn integer default 0
);

create table if not exists withdrawals (
    user_id uuid not null references users(user_id),
    order_number varchar(255) not null,
    sum integer not null,
    processed_at timestamp not null
);

commit;