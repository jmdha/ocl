package main

import (
	"log"
	"net/http"
)

func routeIndex(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeMetrics(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("select method, path, count(*), avg(duration) from requests group by method, path order by count(*) desc;")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type rowdata struct {
		Method string  `json:"method"`
		Path   string  `json:"path"`
		Calls  uint    `json:"calls"`
		Avg    float64 `json:"avg"`
	}
	var data []rowdata
	for rows.Next() {
		var rd rowdata
		err = rows.Scan(&rd.Method, &rd.Path, &rd.Calls, &rd.Avg)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		rd.Avg = rd.Avg / 1e6
		data = append(data, rd)
	}

	err = Templates.ExecuteTemplate(w, "metrics.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
