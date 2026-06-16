create table request (
	id        integer primary key,
	timestamp integer not null,
	ip        text    not null,
	method    text    not null,
        path      text    not null,
	query     text    not null,
	agent     text    not null,
	duration  integer not null
);

create table import (
	id        integer primary key,
	timestamp integer not null,
	status    text    not null,
	error     text,
	data      blob,

	check (status in ('pending', 'uploading', 'processing', 'done', 'failed'))
);

create table encounter (
	id              integer primary key,
	import_id       integer not null,

	timestamp_start integer not null,
	timestamp_end   integer,
	encounter_id    integer not null,
	encounter_name  integer not null,
	group_size      integer not null,
	difficulty      integer not null,
	success         integer,

	check (success is null or success in (0, 1))
);

create table character (
	id   integer primary key,
	guid text not null unique,
	name text not null
);

create table event_damage (
	id        integer primary key,
	import_id integer not null,

	timestamp integer not null,
	source_id integer not null,
	target_id integer not null,
	spell_id  integer not null,
	raw       integer not null,
	amount    integer not null
);

create table event_heal (
	id        integer primary key,
	import_id integer not null,

	timestamp integer not null,
	source_id integer not null,
	target_id integer not null,
	spell_id  integer not null,
	amount    integer not null
);

create index import_timestamp on import(timestamp);
