package db

func CharacterID(guid string, name string) (int64, error) {
	var out int64
	var err error

	err = db.QueryRow(`
		insert into character (guid, name)
		values (?, ?)
		on conflict(guid) do update
		set guid = excluded.guid
		returning id`,
		guid,
		name,
	).Scan(&out)

	return out, err
}
