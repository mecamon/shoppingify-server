CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    name varchar(30),
    note varchar(100),
    category_id int,
    image_url varchar,
    is_active boolean,
    created_at bigint,
    updated_at bigint,
    CONSTRAINT fk_category
        FOREIGN KEY (category_id)
            REFERENCES categories(id)
);