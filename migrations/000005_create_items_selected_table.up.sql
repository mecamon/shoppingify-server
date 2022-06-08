CREATE TABLE items_selected(
    id SERIAL PRIMARY KEY,
    item_id int,
    quantity int,
    is_completed bool,
    list_id int,
    created_at bigint,
    updated_at bigint,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items(id),
    CONSTRAINT fk_list FOREIGN KEY (list_id) REFERENCES lists(id)
);