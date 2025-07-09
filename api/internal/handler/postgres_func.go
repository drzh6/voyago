package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
)

func RowsToJSON(rows pgx.Rows) ([]byte, error) {
	defer rows.Close()

	columns := rows.FieldDescriptions()
	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[string(col.Name)] = values[i]
		}

		result = append(result, rowMap)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return json.Marshal(result)
}
