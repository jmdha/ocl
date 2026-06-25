create table import (
	id        integer primary key,
	timestamp integer not null,
	path      text    not null,
	status    text    not null,

	check (status in ('pending', 'processing', 'done', 'failed'))
);

create table character (
	id   integer primary key,
	name text not null unique
);

create table spell (
	id   integer primary key,
	name text not null
);

create table encounter (
	id             integer primary key,
	import_id      integer not null,

	encounter_id   integer not null,
	encounter_name text    not null,
	start          integer not null,
	end            integer,
	success        integer,

	check (success is not null or success in (0, 1))
);

create table challenge (
	id             integer primary key,
	import_id      integer not null,

	zone_name      text not null,
	instance_id    integer not null,
	challenge_id   integer not null,
	level          integer not null,
	start          integer not null,
	end            integer,
	success        integer,

	check (level between 0 and 100),
	check (success is not null or success in (0, 1))
);

create table event_damage (
	id         integer primary key,
	import_id  integer not null,

	timestamp  integer not null,
	source_id  integer not null,
	target_id  integer not null,
	spell_id   integer not null,
	amount     integer not null
);

create index import_timestamp       on import(timestamp);
