-- +migrate Up

create table states
(
    id           text primary key not null,
    operation    text unique      not null,
    confirmation text             not null
);

create table transitions
(
    tx    text primary key not null,
    state text             not null,
    chain text             not null
);

create index transitions_index on transitions (state);

-- +migrate Down
drop table states;
drop table transitions;