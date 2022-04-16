create table museum
(
    id                  serial primary key,
    name                text not null,
    info                json not null default '{}',
    description         text not null default '',
    image               text not null default '',
    image_height        int not null default 0,
    image_width         int not null default 0,
    popular             bigint not null default 0
);

insert into museum(name, info, description, image, image_height, image_width) 
values('Museum', '{"Страна":"Россия","Город":"Москва","Адрес":"2-я Бауманская, 5","Год":"2022","Руководитель":"Гордин Михаил"}'::json, 'Most beautiful museum in Moscow', 'default.jpg', 835, 600);

create table exhibition
(
    id           serial primary key,
    museum_id    int not null,
    name         text not null,
    description  text not null,
    image        text not null,
    image_height int not null,
    image_width  int not null,
    info         json not null default '{}',
    popular      bigint not null default 0
);

insert into exhibition(museum_id, name, description, info, image, image_height, image_width) 
values(1, 'Exhibition', 'Most beautiful exhibition in Moscow', '{"Начало":"2022-10-08","Конец":"2022-10-11"}'::json, 'default.jpg', 835, 600);

create table picture
(
    id          serial primary key,
    exh_id      int not null default 0,
    name        text not null,
    image       text not null,
    description text not null,
    info        json not null default '{}',
    height      int not null,
    width       int not null,
    popular      bigint not null default 0
);

insert into picture(exh_id, name, image, description, info, height, width) 
values(1, 'Cat', 'default.jpg,notfound.jpg', 'First picture in the app', '{"Автор":"Человек","Год":"2021","Техника":"Компьютерная графика","Размер":"3 х 2"}'::json, 835, 600);

create table users 
(
    id       serial primary key,
    login    text unique not null,
    password bytea not null,
    museum   int not null
);