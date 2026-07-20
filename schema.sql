create table if not exists requests (
	id     integer primary key,
	ip     string,
	method string,
	uri    string,
	start  integer not null,
	end    integer
);

create table if not exists upload (
	id        integer primary key,
	requestID integer not null,
	userID    integer,
	size      integer
);

create table if not exists logfile (
	id     integer primary key,
	userID integer
);

create table if not exists run (
	id     integer primary key
);
