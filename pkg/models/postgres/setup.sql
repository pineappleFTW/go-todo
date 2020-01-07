--  psql -f ./pkg/models/postgres/setup.sql -d todo

create table todos (
    id serial primary key not null,
    title varchar(255) not null,
    content text not null,
    created timestamp not null
);

create table users (
    id serial primary key not null,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    role int not null,
    active boolean not null default true,
    created timestamp not null
);

create table refresh_tokens (
    id serial primary key not null,
    identifier varchar(255) not null, 
    token varchar(255) not null,
    user_id int references users(id) not null,
    created timestamp not null,
    updated timestamp not null
);

alter table refresh_tokens add constraint refresh_tokens_uc_identifer unique(identifier);

alter table refresh_tokens add constraint refresh_tokens_uc_token unique(token);

alter table users add constraint users_uc_email unique (email);