package db

import (
	"bytes"
	"time"
)

func ImportAdd() (int64, error) {
	res, err := db.Exec(
		`insert into imports (timestamp, status) values (?, 'uploading')`,
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func ImportMarkFailed(id int64, reason string) error {
	_, err := db.Exec(
		`update imports set status = 'failed', error = ? where id = ?`,
		reason,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func ImportMarkPending(id int64, data bytes.Buffer) error {
	_, err := db.Exec(
		`update imports set status = 'pending', data = ? where id = ?`,
		data.Bytes(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func ImportMarkDone(id int64) error {
	_, err := db.Exec(
		`update imports set status = 'done' where id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func ImportClaim() (int64, []byte, error) {
	res := db.QueryRow(`
		update imports
		set status = 'processing'
		where id = (
			select id
			from imports
			where status = 'pending'
			order by id
			limit 1
		)
		returning id, data
	`)

	var id int64
	var data []byte
	err := res.Scan(&id, &data)
	if err != nil {
		return 0, nil, err
	}

	return id, data, nil
}
