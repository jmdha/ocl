package db

import "time"

func RequestAdd(
	ip string,
	method string,
	path string,
	query string,
	agent string,
	duration time.Duration,
) error {
	var err error

	_, err = db.Exec(`
		insert into request (timestamp, ip, method, path, query, agent, duration)
		values (?, ?, ?, ?, ?, ?, ?)`,
		time.Now().UnixNano(),
		ip,
		method,
		path,
		query,
		agent,
		duration.Nanoseconds(),
	)

	return err
}

func RequestCount() (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		select count(*)
		from request
	`).Scan(&out)

	return out, err
}

func RequestCount24H() (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		select count(*)
		from request
		where timestamp >= (unixepoch() - 86400) * 1000
	`).Scan(&out)

	return out, err
}

func RequestIPCount() (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		select count(distinct ip)
		from request
	`).Scan(&out)

	return out, err
}

func RequestIPCount24H() (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		select count(distinct ip)
		from request
		where timestamp >= (unixepoch() - 86400) * 1000
	`).Scan(&out)

	return out, err
}
