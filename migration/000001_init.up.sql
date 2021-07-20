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



INSERT INTO users (name, surname, username) VALUES('Ivan','SMITH','admin1');
INSERT INTO users (name, surname, username) VALUES('Oleg','BROWN','admin2');
INSERT INTO users (name, surname, username) VALUES('Kiril','WILSON','admin3');
INSERT INTO users (name, surname, username) VALUES('Sergei','THOMSON','admin4');
INSERT INTO users (name, surname, username) VALUES('Mishail','ROBERTSON','admin5');


INSERT INTO wallets (amount,currency,user_id) VALUES(1000.15,'€',1);
INSERT INTO wallets (amount,currency,user_id) VALUES(1020.15,'$',1);
INSERT INTO wallets (amount,currency,user_id) VALUES(1500.05,'€',2);
INSERT INTO wallets (amount,currency,user_id) VALUES(1009.45,'$',2);
INSERT INTO wallets (amount,currency,user_id) VALUES(1210.09,'€',3);
INSERT INTO wallets (amount,currency,user_id) VALUES(1251.19,'$',3);
INSERT INTO wallets (amount,currency,user_id) VALUES(2000.1,'€',4);
INSERT INTO wallets (amount,currency,user_id) VALUES(999.99,'$',4);
INSERT INTO wallets (amount,currency,user_id) VALUES(1014.66,'€',5);
INSERT INTO wallets (amount,currency,user_id) VALUES(1154.15,'$',5);

INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(10.15,'$', 'test', false, '2021-01-22 19:10:25', 1);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(11.00,'€', 'test', false, '2021-02-15 19:10:25', 1);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(11.15,'$', 'test', true, '2021-03-02 19:10:25', 2);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(111.19,'€', 'test', true, '2021-03-09 19:10:25', 2);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(119.10,'$', 'test', false, '2021-02-17 19:10:25', 3);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1000.00,'€', 'test', true, '2021-05-10 19:10:25', 3);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(2000.00,'$', 'test', false, '2021-01-07 19:10:25', 4);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(5000.00,'€', 'test', true, '2021-02-01 19:10:25', 4);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(666.6,'$', 'test', true, '2021-05-22 04:10:25', 5);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1024.15,'€', 'test', true, '2021-01-26 19:10:25', 5);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'$', 'test', true, '2021-02-24 19:10:25', 6);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'€', 'test', true, '2021-04-21 19:10:25', 6);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'$', 'test', true, '2021-03-08 19:10:25', 7);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'€', 'test', true, '2021-01-16 19:10:25', 7);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'$', 'test', true, '2021-02-27 19:10:25', 8);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'€', 'test', true, '2021-05-21 19:10:25', 8);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'$', 'test', true, '2021-06-13 19:10:25', 9);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'€', 'test', true, '2021-01-19 19:10:25', 9);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'$', 'test', true, '2021-02-10 19:10:25', 10);
INSERT INTO transaction_history (amount, currency, description, done, datetime, wallet_id) VALUES(1154.15,'€', 'test', true, '2021-04-26 19:10:25', 10);

