CREATE DATABASE IF NOT EXISTS todo;

USE todo;

CREATE TABLE user (
    id int,
    name char(50),
    email char(50),
    PRIMARY KEY (id),
);

CREATE TABLE tasks (
    id int,
    summary varchar(100),
    description varchar(500),
    user_id int,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);
