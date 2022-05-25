CREATE TABLE categories(
    id SERIAL PRIMARY KEY,
    name varchar(30),
    user_id int,
    created_at bigint,
    updated_at bigint,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);