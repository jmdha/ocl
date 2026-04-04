create table if not exists requests (
	id        integer primary key autoincrement,
	method    text not null,
	path      text not null,
	query     text not null,
	ip        text not null,
	agent     text not null,
	duration  integer not null,
	status    integer not null,
	timestamp datetime not null
);

create table if not exists users (
	id        integer primary key autoincrement,
	name      text not null,
	email     text not null,
	timestamp datetime not null
);

create table if not exists uploads (
	id        integer primary key autoincrement,
	user      integer,
	ip        text not null,
	timestamp datetime not null
);

create table if not exists logs (
	id        integer primary key autoincrement,
	user      integer,
	timestamp datetime not null
);

create table if not exists reports (
	id        integer primary key autoincrement,
	timestamp datetime not null
);

create table if not exists report_logs (
	report integer not null,
	log    integer not null primary key
);
