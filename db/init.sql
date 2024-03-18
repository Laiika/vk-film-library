create table if not exists users
(
    id       int generated always as identity primary key,
    name     text unique not null,
    password text not null,
    role     text not null
);

create table if not exists actors
(
    id       int generated always as identity primary key,
    name     text not null,
    gender   text not null,
    birthday date not null
);

create table if not exists films
(
    id          int generated always as identity primary key,
    name        text not null,
    description text not null,
    created_at  date not null,
    rating      int not null
);

create table if not exists films_actors
(
    id       int generated always as identity primary key,
    film_id  int not null,
    actor_id int not null,

    foreign key (actor_is) references actors(id) on delete cascade,
    foreign key (film_id) references films(id) on delete cascade
);