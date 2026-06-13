package db

import (
	"database/sql"
	"time"
)

func ImportAdd() (int64, error) {
	var res sql.Result
	var out int64
	var err error

	res, err = db.Exec(`
		insert into import(timestamp, status)
		values (?, 'uploading')`,
		time.Now().UnixNano(),
	)
	if err != nil {
		return 0, err
	}

	out, err = res.LastInsertId()

	return out, err
}

func ImportSetPending(id int64, buf []byte) error {
	var err error

	_, err = db.Exec(`
		update import set status = 'pending', data = ?
		where id = ?`,
		buf,
		id,
	)

	return err
}

func ImportSetFailed(id int64, reason string) error {
	var err error

	_, err = db.Exec(`
		update import set status = 'failed', error = ?
		where id = ?`,
		reason,
		id,
	)

	return err
}

func ImportSetDone(id int64) error {
	var err error

	_, err = db.Exec(`
		update import set status = 'done'
		where id = ?`,
		id,
	)

	return err
}

func ImportClaim() (int64, []byte, error) {
	var id int64
	var buf []byte
	var err error

	err = db.QueryRow(`
		update import set status = 'processing'
		where id = (
			select id
			from import
			where status = 'pending' or 
			      status = 'processing'
			limit 1
		)
		returning id, data
	`).Scan(&id, &buf)

	return id, buf, err
}

func ImportStatusCount(status string) (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		select count(*)
		from import
		where status = ?`,
		status,
	).Scan(&out)

	return out, err
}
