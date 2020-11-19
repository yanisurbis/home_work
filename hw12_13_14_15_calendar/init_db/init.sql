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

create table notifications
(
    event_id    int,
    user_id     int,
    start_at    timestamp,
    event_title varchar(100),
    constraint fk_event FOREIGN KEY (event_id) REFERENCES events (id),
    constraint fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);
