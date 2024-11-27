package main

import "database/sql"

func getReleases(db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, name, year FROM releases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var year int
		err := rows.Scan(&id, &name, &year)
		if err != nil {
			return nil, err
		}
		items = append(items, map[string]interface{}{
			"id":   id,
			"name": name,
			"year": year,
		})
	}

	return items, nil
}
