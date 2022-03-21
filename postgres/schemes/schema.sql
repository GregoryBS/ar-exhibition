create table museum
(
    id                  serial primary key,
    name                text not null,
    country             text not null,
    city                text not null,
    address             text,
    year                int,
    description         text not null,
    director            text,
    image               text not null,
    image_height        int not null,
    image_width         int not null,
    popular             bigint not null default 0
);

insert into museum(name, country, city, address, year, description, director, image, image_height, image_width) 
values('Museum', 'Russia', 'Moscow', '2-ya Baumanskaya, 5', 2022, 'Most beautiful museum in Moscow', 'Gordin Michael', 'default.jpg', 835, 600);

create table exhibition
(
    id           serial primary key,
    museum_id    int not null,
    name         text not null,
    description  text not null,
    image        text not null,
    image_height int not null,
    image_width  int not null,
    date_from    timestamp with time zone,
    date_to      timestamp with time zone,
    popular      bigint not null default 0
);

insert into exhibition(museum_id, name, description, date_from, date_to, image, image_height, image_width) 
values(1, 'Exhibition', 'Most beautiful exhibition in Moscow', '2022-10-08', '2022-10-11', 'default.jpg', 835, 600);

create table picture
(
    id          serial primary key,
    exh_id      int not null,
    name        text not null,
    technique   text not null,
    image       text not null,
    author      text not null,
    year        int,
    height      int not null,
    width       int not null
);

insert into picture(exh_id, name, technique, image, author, year, height, width) 
values(1, 'Cat', 'Computer graphics', 'default.jpg', 'GregoryBS', 2021, 835, 600);