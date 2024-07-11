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

CREATE TABLE IF NOT EXISTS skutags
(
    tag             varchar(255)    not null primary key references tags(id) on delete cascade,
    sku             uuid            not null primary key references skus(id) on delete cascade
);