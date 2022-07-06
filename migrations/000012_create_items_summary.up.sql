CREATE TABLE items_summary(
    id SERIAL PRIMARY KEY,
    user_id int,
    month int,
    year int,
    quantity int,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
)