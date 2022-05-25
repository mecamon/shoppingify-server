CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(40),
    lastname varchar(40),
    email varchar(40),
    password varchar,
    is_active bool,
    is_visitor bool,
    login_code UUID,
    created_at bigint,
    updated_at bigint
);