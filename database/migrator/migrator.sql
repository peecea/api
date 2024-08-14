drop database if exists peec;
set sql_mode = '';
create database peec;
use peec;

create table user
(
    id          int auto_increment
        primary key,
    created_at  datetime     default CURRENT_TIMESTAMP     not null,
    updated_at  datetime     default CURRENT_TIMESTAMP     not null,
    deleted_at  datetime     default '0000-00-00 00:00:00' null,
    name        varchar(500) default ''                    null,
    family_name varchar(500) default ''                    null,
    nick_name   varchar(100) default ''                    null,
    email       varchar(100) default ''                    null,
    matricule   varchar(32)  default ''                    null,
    age         int          default 0                     null,
    birth_date  datetime     default '0000-00-00 00:00:00' null,
    sex         int          default 0                     null,
    status      int          default 0                     null,
    constraint user_pk
        unique (email)
);

create table authorization
(
    id         int auto_increment
        primary key,
    created_at datetime default CURRENT_TIMESTAMP     not null,
    updated_at datetime default CURRENT_TIMESTAMP     not null,
    deleted_at datetime default '0000-00-00 00:00:00' null,
    user_id    int      default 0                     null,
    level      int      default 0                     null,
    constraint authorization_pk
        unique (user_id, level),
    constraint authorization_user_id_fk
        foreign key (user_id) references user (id)
            on update cascade on delete cascade
);

create table password
(
    id           int auto_increment
        primary key,
    created_at   datetime      default CURRENT_TIMESTAMP     not null,
    updated_at   datetime      default CURRENT_TIMESTAMP     not null,
    deleted_at   datetime      default '0000-00-00 00:00:00' null,
    user_id      int           default 0                     null,
    psw          varchar(1000) default ''                    null,
    content_hash varchar(500)  default ''                    null,
    constraint password_user_id_fk
        foreign key (user_id) references user (id)
);

create index password_user_id_index
    on password (user_id);

create table media
(
    id           int auto_increment
        primary key,
    created_at   datetime     default CURRENT_TIMESTAMP     not null,
    updated_at   datetime     default CURRENT_TIMESTAMP     not null,
    deleted_at   datetime     default '0000-00-00 00:00:00' null,
    file_name    varchar(500) default ''                    null,
    extension    varchar(10)  default ''                    null,
    xid          varchar(500) default ''                    null,
    user_id      int          default 0                     null,
    content_type int          default 0                     null
);

create table code
(
    id           int auto_increment
        primary key,
    created_at   datetime      default CURRENT_TIMESTAMP     not null,
    updated_at   datetime      default CURRENT_TIMESTAMP     not null,
    deleted_at   datetime      default '0000-00-00 00:00:00' null,
    user_id      int           default 0                     null,
    verification_code          int default 0                    null
);

create index code_user_val
    on code(user_id,verification_code);

create table address
(
    id           int auto_increment      primary key not null,
    created_at   datetime      default CURRENT_TIMESTAMP    ,
    updated_at   datetime      default CURRENT_TIMESTAMP    ,
    deleted_at   datetime      default '0000-00-00 00:00:00' ,
    country varchar(100) default '' ,
    city varchar(100) default '',
    latitude float default 0,
    longitude float default 0,
    street varchar(100) ,
    full_address varchar(600),
    xid          varchar(500) default ''
);

create table user_address
(
    id           int auto_increment      primary key not null,
    created_at   datetime      default CURRENT_TIMESTAMP    ,
    updated_at   datetime      default CURRENT_TIMESTAMP    ,
    deleted_at   datetime      default '0000-00-00 00:00:00' ,
    user_id int unique,
    address_id int unique,
    address_type varchar(100) default '',
    foreign key (user_id) references user(id),
    foreign key (address_id) references address(id)
);

create table thumb
(
    id           int auto_increment
        primary key,
    created_at   datetime     default CURRENT_TIMESTAMP     not null,
    updated_at   datetime     default CURRENT_TIMESTAMP     not null,
    deleted_at   datetime     default '0000-00-00 00:00:00' null,
    file_name    varchar(500) default ''                    null,
    extension    varchar(10)  default ''                    null,
    media_xid          varchar(500) default ''                    null,
    content_type int          default 0                     null
);

alter table user add profile_image_xid varchar(500) default '' after status ;

alter table thumb
    rename media_thumb,
    add column xid varchar(500) default '',
    drop column content_type ,
    drop column file_name ;


create table user_media_detail
(
   id     int auto_increment       primary key,
   created_at   datetime     default CURRENT_TIMESTAMP ,
   updated_at   datetime     default CURRENT_TIMESTAMP ,
   deleted_at   datetime     default '0000-00-00 00:00:00',
   owner_id int default 0,
   document_type int default 0,
   document_xid varchar(100) unique default ''
);

alter table user_media_detail
    add constraint user_media_detail_user_id_fk
        foreign key (owner_id) references user (id);


alter table media drop column user_id;


create table qr_code_registry (
    id int primary key auto_increment ,
    created_at   datetime     default CURRENT_TIMESTAMP ,
    deleted_at   datetime     default '0000-00-00 00:00:00',
    user_id int unique default 0,
    xid varchar(100) unique default '',
    is_used  boolean ,
    foreign key (user_id) references user(id)
);

create table education(
    id int primary key auto_increment,
    created_at   datetime     default CURRENT_TIMESTAMP ,
    updated_at   datetime     default CURRENT_TIMESTAMP ,
    deleted_at   datetime     default '0000-00-00 00:00:00',
    name varchar(500) default ''
);

create table subject (
    id int primary key auto_increment,
    created_at   datetime     default CURRENT_TIMESTAMP ,
    updated_at   datetime     default CURRENT_TIMESTAMP ,
    deleted_at   datetime     default '0000-00-00 00:00:00',
    education_level_id int default 0,
    name varchar(500) default '',
    subject_code varchar(500) default ''
);

alter table subject
    add constraint subject_education_id_fk
        foreign key (education_level_id) references education (id);

insert into education (name)
values      ('L1'),
            ('L2'),
            ('L3'),
            ('LP1'),
            ('LP2'),
            ('LP3'),
            ('M1'),
            ('M2'),
            ('MP1'),
            ('MP2'),
            ('MVR'),
            ('D1'),
            ('D2');
--
--
--  INDIRECTION TABLE BETWEEN SUBJECT AND USER
--
--
create table user_education_level_subject
(
    id int primary key auto_increment,
    created_at   datetime     default CURRENT_TIMESTAMP ,
    updated_at   datetime     default CURRENT_TIMESTAMP ,
    deleted_at   datetime     default '0000-00-00 00:00:00',
    user_id int default 0 unique,
    subject_id int default 0 unique
);


alter table user_education_level_subject
    add constraint user_education_level_subject_user_id_fk
        foreign key (user_id) references user (id);


alter table user_education_level_subject
    add constraint user_education_level_subject_subject_id_fk
        foreign key (subject_id) references subject (id);

