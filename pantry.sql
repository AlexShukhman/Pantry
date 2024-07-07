CREATE TABLE IF NOT EXISTS PantrySKUs
(
    id              integer         not null primary key,
    sku_name        text            not null unique,
    sku_quantity    smallint
);
