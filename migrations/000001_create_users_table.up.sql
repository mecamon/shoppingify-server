CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(40),
    lastname varchar(40),
    email varchar(40),
    is_active bool,
    is_visitor bool,
    logging_code UUID,
    created_at timestamp,
    updated_at timestamp
);