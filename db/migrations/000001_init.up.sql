create table requests (
	id        integer primary key,
	timestamp integer not null,
	ip        text not null,
	method    text not null,
        path      text not null,
	query     text not null,
	agent     text not null,
	duration  integer not null
);

create table imports (
	id        integer primary key,
	timestamp integer not null,
	status    text not null check (status in ('pending', 'uploading', 'processing', 'done', 'failed')),
	error     text,
	processed integer default 0
);

create table wow_event_version (
	id        integer primary key,
	import_id integer not null,
	line      integer not null,
	timestamp integer not null,
	log       integer not null,
	advanced  integer not null,
	major     integer not null,
	minor     integer not null,
	patch     integer not null,
	project   integer not null,

	foreign key(import_id) references imports(id)
);

create table wow_event_encounter_start (
	id             integer primary key,
	import_id      integer not null,
	line           integer not null,
	timestamp      integer not null,
	encounter_id   integer not null,
	encounter_name text not null,
	difficulty_id  integer not null,
	group_size     integer not null,
	instance_id    integer not null,

	foreign key(import_id) references imports(id)
);

create table wow_event_encounter_end (
	id             integer primary key,
	import_id      integer not null,
	line           integer not null,
	timestamp      integer not null,
	encounter_id   integer not null,
	encounter_name text not null,
	difficulty_id  integer not null,
	group_size     integer not null,
	success        integer not null,
	duration       integer not null,

	foreign key(import_id) references imports(id)
);
