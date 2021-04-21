CREATE DATABASE IF NOT EXISTS Recycle;

USE Recycle;

CREATE TABLE IF NOT EXISTS Orders
(
    id        bigint unsigned primary key not null,
    eventAt   DATETIME                    not null,
    address   text                        not null,
    size      integer                     not null, # maximum place count
    createdAt DATETIME                    not null,
    updatedAt DATETIME                    not null
);


CREATE TABLE IF NOT EXISTS OrdersUsers
(
    user_id  bigint,
    order_id bigint,
    CONSTRAINT unqUO UNIQUE (user_id, order_id)
);