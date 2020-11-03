create database events;
create user yanis with encrypted password 'yanis';
grant all privileges on database events to yanis;

DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id   serial primary key,
    name text
);

insert into users (name)
values ('Ann'),
       ('Bob'),
       ('Carl');

create table events
(
    id          serial primary key,
    title       varchar(100),
    start_at    timestamp,
    end_at      timestamp,
    description text,
    user_id     int,
    notify_at   timestamp,
    constraint fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

insert into events (title, start_at, end_at, description, user_id, notify_at)
values ('Event 1', current_timestamp, current_timestamp, 'Description 1', 1, current_timestamp),
       ('Event 2', current_timestamp, current_timestamp, 'Description 2', 1, current_timestamp),
       ('Event 3', current_timestamp, current_timestamp, 'Description 3', 2, current_timestamp),
       ('Event 4', current_timestamp, current_timestamp, 'Description 4', 2, current_timestamp),
       ('Event 5', current_timestamp, current_timestamp, 'Description 5', 1,
        current_timestamp + interval '250 hours'),
       ('Event 6', current_timestamp + interval '48 hours', current_timestamp, 'Description 5', 1,
        current_timestamp + interval '250 hours');



