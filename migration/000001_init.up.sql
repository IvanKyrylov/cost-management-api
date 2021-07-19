CREATE TABLE users (
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null,
    username varchar(255) not null unique
);

CREATE TABLE wallets (
    id serial not null unique,
    amount numeric(1000, 4) not null,
    currency varchar(255) not null,
    user_id int references users(id) on delete cascade not null
);

CREATE TABLE transaction_history (
    id serial not null unique,
    amount numeric(1000, 4) not null,
    currency varchar(255) not null,
    description varchar(255),
    done boolean not null default false,
    datetime TIMESTAMP not null,
    wallet_id int references wallets(id) on delete cascade not null
);