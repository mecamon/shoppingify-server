CREATE TABLE categories(
    id SERIAL PRIMARY KEY,
    name varchar(30),
    user_id int,
    created_at timestamp,
    updated_at timestamp,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);