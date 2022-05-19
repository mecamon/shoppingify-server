CREATE TABLE top_categories(
    id SERIAL PRIMARY KEY,
    user_id int,
    category_id int,
    sum_quantity int,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id)
);