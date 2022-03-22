create table museum
(
    id                  serial primary key,
    name                text not null,
    info                json not null default '{}',
    description         text not null,
    image               text not null,
    image_height        int not null,
    image_width         int not null,
    popular             bigint not null default 0
);

insert into museum(name, info, description, image, image_height, image_width) 
values('Museum', '{"country":"Russia","city":"Moscow","address";"2-ya Baumanskaya, 5","year":"2022","director":"Gordin Michael"}'::json, 'Most beautiful museum in Moscow', 'default.jpg', 835, 600);

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
values(1, 'Exhibition', 'Most beautiful exhibition in Moscow', '{"date from":"2022-10-08","date to":"2022-10-11"}'::json, 'default.jpg', 835, 600);

create table picture
(
    id          serial primary key,
    exh_id      int not null,
    name        text not null,
    image       text not null,
    description text not null,
    info        json not null default '{}',
    height      int not null,
    width       int not null
);

insert into picture(exh_id, name, image, description, info, height, width) 
values(1, 'Cat', 'default.jpg', 'First picture in the app', '{"author":"GregoryBS","year":"2021","technique":"Computer graphics"}'::json, 835, 600);