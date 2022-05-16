
create table museum
(
    id             serial primary key,
    name           text not null,
    info           json not null default '{}',
    description    text not null default '',
    image          text not null default '',
    image_height   real not null default 0,
    image_width    real not null default 0,
    popular        bigint not null default 0,
    user_id        int not null default 0,
    mus_show       boolean not null default false
);

insert into museum(name, info, description, image, image_height, image_width, mus_show) 
values('Museum', '{"Страна":"Россия","Город":"Москва","Адрес":"2-я Бауманская, 5","Год":"2022","Руководитель":"Гордин Михаил"}'::json, 'Most beautiful museum in Moscow', 'default.jpg', 835, 600, true);

create table exhibition
(
    id           serial primary key,
    museum_id    int not null default 0,
    name         text not null,
    description  text not null,
    image        text not null default '',
    image_height real not null default 0,
    image_width  real not null default 0,
    info         json not null default '{}',
    popular      bigint not null default 0,
    user_id      int not null default 0,
    exh_show     boolean not null default false,
    mus_show     boolean not null default false
);

insert into exhibition(museum_id, name, description, info, image, image_height, image_width, exh_show, mus_show) 
values(1, 'Exhibition', 'Most beautiful exhibition in Moscow', '{"Начало":"2022-05-08","Конец":"2022-10-11"}'::json, 'default.jpg', 835, 600, true, true);

create table picture
(
    id          serial primary key,
    exh_id      int[] not null default '{}',
    name        text not null,
    image       text not null default '',
    description text not null,
    info        json not null default '{}',
    height      real not null default 0,
    width       real not null default 0,
    popular     bigint not null default 0,
    user_id     int not null default 0,
    video       text not null default '',
    video_size  text not null default '',
    pic_show    boolean not null default false,
    exh_show    boolean[] not null default '{}',
    mus_show    boolean not null default false
);

insert into picture(exh_id, name, image, description, info, height, width, pic_show, exh_show, mus_show) 
values('{1}', 'Cat', 'default.jpg,notfound.jpg', 'First picture in the app', '{"Автор":"Человек","Год":"2021","Техника":"Компьютерная графика","Размер":"3 х 2"}'::json, 835, 600, true, '{1}', true);

create table users 
(
    id       serial primary key,
    login    text unique not null,
    password bytea not null,
    admin    boolean not null default false
);

create table stats 
(
    id      bigserial primary key,
    port    int not null,
    method  text not null,
    status  int not null,
    urls    text not null,
    reqtime timestamp default now()::timestamp,
    perform int
);
