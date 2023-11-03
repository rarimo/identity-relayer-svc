-- +migrate Up

create table gists
(
    id           text primary key not null,
    operation    text unique      not null,
    confirmation text             not null
);

create table gist_transitions
(
    tx    text primary key not null,
    gist  text             not null,
    chain text             not null
);

create index gist_transitions_index on gist_transitions (gist);

-- +migrate Down
drop table gists;
drop table gist_transitions;