CREATE TABLE top_items(
    id SERIAL PRIMARY KEY,
    user_id int,
    item_id int,
    sum_quantity int,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items(id)
);