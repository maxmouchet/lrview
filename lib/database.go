package lib

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func QueryAll[T interface{}](db *sqlx.DB, query string, args ...interface{}) ([]T, error) {
	objects := make([]T, 0)
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var object T
		err := rows.StructScan(&object)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func QueryOne[T interface{}](db *sqlx.DB, query string, args ...interface{}) (*T, error) {
	objects, err := QueryAll[T](db, query, args...)
	if err != nil {
		return nil, err
	}
	if len(objects) != 1 {
		return nil, fmt.Errorf("got %d rows, expected exactly one", len(objects))
	}
	return &objects[0], nil
}
