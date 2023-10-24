CREATE TABLE spaceship
(
    id     int primary key not null auto_increment,
    name   varchar(255)    not null,
    class  varchar(255)    not null,
    crew   int             not null,
    image  varchar(255)    not null,
    value  decimal(13, 2)  not null,
    status varchar(255)    not null
);

CREATE TABLE armament
(
    id   int primary key not null auto_increment,
    name varchar(255)    not null
);

CREATE TABLE spaceship_armament
(
    id           int primary key not null auto_increment,
    spaceship_id int             not null,
    armament_id  int             not null
);
