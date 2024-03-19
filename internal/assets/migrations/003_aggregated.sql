-- +migrate Up

create table aggregateds
(
    gist         text unique not null,
    state_root   text unique not null,
    operation    text unique not null,
    confirmation text        not null,
    primary key (gist, state_root)
);

create table aggregated_transitions
(
    tx         text primary key not null,
    gist       text             not null,
    state_root text             not null,
    chain      text             not null
);

create index aggregated_transitions_index on aggregated_transitions (gist, state_root);

-- +migrate Down
drop table aggregateds;
drop table aggregated_transitions;