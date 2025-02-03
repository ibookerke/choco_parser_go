CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    phone TEXT,
    birthday varchar(50),
    orders_count int NOT NULL default 1,
    full_name TEXT
);