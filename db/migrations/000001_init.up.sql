-- Denotes a http request
-- Used for auditing, metrics, etc
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

-- Denotes an upload
-- Used for tracking who uploaded what, processing queue
create table imports (
	id        integer primary key,
	timestamp integer not null,
	status    text not null check (status in ('pending', 'uploading', 'processing', 'done', 'failed')),
	error     text,
	data      blob
);

create table event_encounter_start (
	id            integer primary key,
	import_id     integer not null,
	line          integer not null,
	timestamp     integer not null,
	encounter_id  integer not null,
	difficulty_id integer not null,
	group_size    integer not null,

	foreign key(import_id) references imports(id)
);

create table event_encounter_end (
	id            integer primary key,
	import_id     integer not null,
	line          integer not null,
	timestamp     integer not null,
	encounter_id  integer not null,
	difficulty_id integer not null,
	group_size    integer not null,
	success       integer not null,
	duration      integer not null,

	foreign key(import_id) references imports(id)
);

create table event_damage (
	id        integer primary key,
	import_id integer not null,
	line      integer not null,
	timestamp integer not null,
	spell_id  integer not null,
	source_id integer not null,
	target_id integer not null,
	amount    float not null,

	foreign key(import_id) references imports(id)
);
