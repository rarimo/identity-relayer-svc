-- +migrate Up

create table states
(
    id           text primary key not null,
    operation    text unique      not null,
    confirmation text             not null
);

-- +migrate Down
drop table states;