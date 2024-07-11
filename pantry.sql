CREATE TABLE IF NOT EXISTS skus
(
    id              uuid            not null unique primary key default gen_random_uuid(),
    sku_name        text            not null unique,
    sku_quantity    smallint        not null
);

CREATE TABLE IF NOT EXISTS  tags
(
    id              varchar(255)    not null unique primary key
);

CREATE TABLE IF NOT EXISTS sku_tags
(
    tag_id             varchar(255)    not null references tags(id) on delete cascade,
    sku_id             uuid            not null references skus(id) on delete cascade,
    PRIMARY KEY (tag_id, sku_id)
);