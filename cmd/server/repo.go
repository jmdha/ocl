package main

import "database/sql"

type Metrics struct {
	Aggregate MetricsAggregate `json:"aggregate"`
	Routes    []MetricsRoute   `json:"routes"`
}

type MetricsAggregate struct {
	Calls      uint `json:"calls"`
	Calls24    uint `json:"calls24"`
	Visitors   uint `json:"visitors"`
	Visitors24 uint `json:"visitors24"`
}

type MetricsRoute struct {
	Method string  `json:"method"`
	Path   string  `json:"path"`
	Calls  uint    `json:"calls"`
	Avg    float64 `json:"avg"`
}

func repoMetrics() (Metrics, error) {
	var aggregate MetricsAggregate
	var routes []MetricsRoute
	var err error

	aggregate, err = repoMetricsAggregate()
	if err != nil {
		return Metrics{}, err
	}

	routes, err = repoMetricsRoutes()
	if err != nil {
		return Metrics{}, err
	}

	return Metrics{
		Aggregate: aggregate,
		Routes:    routes,
	}, nil
}

func repoMetricsAggregate() (MetricsAggregate, error) {
	var calls uint
	var calls24 uint
	var visitors uint
	var visitors24 uint
	var err error

	calls, err = repoMetricsCalls()
	if err != nil {
		return MetricsAggregate{}, err
	}

	calls24, err = repoMetricsCalls24()
	if err != nil {
		return MetricsAggregate{}, err
	}

	visitors, err = repoMetricsVisitors()
	if err != nil {
		return MetricsAggregate{}, err
	}

	visitors24, err = repoMetricsVisitors24()
	if err != nil {
		return MetricsAggregate{}, err
	}

	return MetricsAggregate{
		Calls:      calls,
		Calls24:    calls24,
		Visitors:   visitors,
		Visitors24: visitors24,
	}, err
}

func repoMetricsRoutes() ([]MetricsRoute, error) {
	var data []MetricsRoute
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`
		select method, path, count(*), avg(duration) from requests 
		group by method, path
		order by count(*) desc;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row MetricsRoute
		err = rows.Scan(&row.Method, &row.Path, &row.Calls, &row.Avg)
		if err != nil {
			return nil, err
		}
		row.Avg = row.Avg / 1e6
		data = append(data, row)
	}

	return data, nil
}

func repoMetricsCalls() (uint, error) {
	var count uint
	var err error

	err = DB.QueryRow(`
		select count(*) from requests
	`).Scan(&count)

	return count, err
}

func repoMetricsCalls24() (uint, error) {
	var count uint
	var err error

	err = DB.QueryRow(`
		select count(*) from requests
		where timestamp >= datetime('now', '-1 day')
	`).Scan(&count)

	return count, err
}

func repoMetricsVisitors() (uint, error) {
	var count uint
	var err error

	err = DB.QueryRow(`
		select count(distinct ip) from requests
	`).Scan(&count)

	return count, err
}

func repoMetricsVisitors24() (uint, error) {
	var count uint
	var err error

	err = DB.QueryRow(`
		select count(distinct ip) from requests
		where timestamp >= datetime('now', '-1 day')
	`).Scan(&count)

	return count, err
}
