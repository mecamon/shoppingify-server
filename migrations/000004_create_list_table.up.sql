CREATE TABLE lists(
    id SERIAL PRIMARY KEY,
    name varchar,
    is_completed bool,
    is_cancelled bool,
    user_id int,
    created_at timestamp,
    updated_at timestamp,
    completed_at timestamp,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);