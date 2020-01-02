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

alter table users add constraint users_uc_email unique (email);