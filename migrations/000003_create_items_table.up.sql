CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    name varchar(30),
    note varchar(100),
    category_id int,
    image_url varchar,
    created_at timestamp,
    updated_at timestamp,
    CONSTRAINT fk_category
        FOREIGN KEY (category_id)
            REFERENCES categories(id)
);