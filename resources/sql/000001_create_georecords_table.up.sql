CREATE TABLE IF NOT EXISTS georecords (
    ip_address    inet not null primary key,
    country_code  text,
    country       text,
    city          text,
    latitude      numeric,
    longitude     numeric,
    mystery_value numeric
);
