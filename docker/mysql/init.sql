CREATE DATABASE IF NOT EXISTS Recycle;

USE Recycle;

CREATE TABLE IF NOT EXISTS Orders
(
    id        bigint unsigned primary key not null,
    address   text                        not null, # creator address
    maxSize   integer                     not null, # maximum place count
    eventAt   DATETIME                    not null,
    createdAt DATETIME                    not null,
    updatedAt DATETIME                    not null
);

CREATE TABLE IF NOT EXISTS OrdersUsers
(
    user_id  bigint,
    order_id bigint,
    address  text not null, # subscribers address
    CONSTRAINT unqUO UNIQUE (user_id, order_id)
);
