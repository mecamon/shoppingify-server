CREATE TABLE items_selected(
    id SERIAL PRIMARY KEY,
    item_id int,
    quantity int,
    list_id int,
    created_at timestamp,
    updated_at timestamp,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items(id),
    CONSTRAINT fk_list FOREIGN KEY (list_id) REFERENCES lists(id)
);