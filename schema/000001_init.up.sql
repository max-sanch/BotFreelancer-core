CREATE TABLE categories
(
    id   serial       not null unique,
    name varchar(256) not null unique
);

CREATE TABLE users
(
    id       serial       not null unique,
    tg_id    numeric      not null unique,
    username varchar(256) not null
);

CREATE TABLE user_settings
(
    id           serial                                          not null unique,
    user_id      integer references users (id) on delete cascade not null unique,
    is_safe_deal boolean                                         not null default false,
    is_budget    boolean                                         not null default false,
    is_term      boolean                                         not null default false
);

CREATE TABLE user_categories
(
    id              serial                                                  not null unique,
    user_setting_id integer references user_settings (id) on delete cascade not null,
    category_id     integer references categories (id) on delete cascade    not null
);

CREATE TABLE channels
(
    id       serial       not null unique,
    api_id   numeric      not null unique,
    api_hash varchar(256) not null,
    name     varchar(256) not null unique
);

CREATE TABLE channel_settings
(
    id           serial                                             not null unique,
    channel_id   integer references channels (id) on delete cascade not null unique,
    is_safe_deal boolean                                            not null default false,
    is_budget    boolean                                            not null default false,
    is_term      boolean                                            not null default false
);

CREATE TABLE channel_categories
(
    id                 serial                                                     not null unique,
    channel_setting_id integer references channel_settings (id) on delete cascade not null,
    category_id        integer references categories (id) on delete cascade       not null
);

CREATE TABLE freelance_tasks
(
    id           serial                                               not null unique,
    task_url     varchar(2048)                                        not null,
    title        varchar(256)                                         not null,
    body         text                                                 not null,
    category_id  integer references categories (id) on delete cascade not null,
    is_budget    boolean                                              not null default false,
    is_term      boolean                                              not null default false,
    is_safe_deal boolean                                              not null default false
);

CREATE TABLE last_parsed_tasks
(
    id       serial                   not null unique,
    datetime timestamp with time zone not null
);