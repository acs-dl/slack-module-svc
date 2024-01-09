-- +migrate Up

create table if not exists responses
(
    id          uuid primary key,
    status      text                        not null,
    error       text,
    description text,
    payload     jsonb,
    created_at  timestamp without time zone not null default current_timestamp
);

create table if not exists users
(
    id         bigint,
    username   TEXT,
    real_name  TEXT,
    slack_id   TEXT PRIMARY KEY,
    updated_at timestamp with time zone not null default current_timestamp,
    created_at timestamp with time zone not null default current_timestamp
);



create index if not exists users_id_idx on users (id);
create index if not exists users_username_idx on users (username);
create index if not exists users_slackid_idx on users (slack_id);

create table if not exists conversations
(
    title          text   not null,
    id             TEXT   not null,
    members_amount bigint not null,

    unique (id)
    );
create index if not exists conversations_id_idx on conversations (id);

create table if not exists links
(
    id   serial primary key,
    link text not null,
    unique (link)
    );
create index if not exists links_link_idx on links (link);

create table if not exists permissions
(
    request_id   text                     not null,
    workspace    text,
    slack_id     TEXT,
    username     TEXT,
    link         text                     not null,
    access_level TEXT                     not null,
    bill         BOOLEAN                  not null,
    created_at   timestamp with time zone not null default current_timestamp,
    updated_at   timestamp with time zone not null default current_timestamp,
    submodule_id text                     not null,
    unique (slack_id, submodule_id),
    foreign key (slack_id) references users (slack_id) on delete cascade on update cascade
    );

create index if not exists permissions_slackid_idx on permissions (slack_id);
create index if not exists permissions_submodule_id_idx on permissions (submodule_id);

-- +migrate Down
DROP TABLE IF EXISTS permissions;
DROP INDEX IF EXISTS links_link_idx;
DROP TABLE IF EXISTS links;
DROP INDEX IF EXISTS conversations_id_idx;
DROP TABLE IF EXISTS conversations;
DROP INDEX IF EXISTS users_slackid_idx;
DROP INDEX IF EXISTS users_username_idx;
DROP INDEX IF EXISTS users_id_idx;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS responses;


