CREATE TABLE lists(
    id SERIAL PRIMARY KEY,
    name varchar,
    is_completed bool,
    is_cancelled bool,
    user_id int,
    created_at bigint,
    updated_at bigint,
    completed_at bigint,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);