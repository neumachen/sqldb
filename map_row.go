package dbsql

import (
	"database/sql"
)

// MappedRow represents a row of mapped column-value pairs.
type MappedRow map[Column]any

// HasColumn returns true if the MappedRow contains the given column.
func (m MappedRow) HasColumn(column Column) bool {
	_, ok := m[column]
	return ok
}

// Count returns the number of columns in the MappedRow.
func (m MappedRow) Count() int {
	return len(m)
}

func (m MappedRow) Get(column Column) (any, bool) {
	v, ok := m[column]
	return v, ok
}

// MapRow maps the columns and values of the given sql.Row to a MappedRow.
func MapRow(row *sql.Row, columns Columns) (MappedRow, error) {
	values := make([]any, columns.Count())
	for i := range values {
		values[i] = new(any)
	}

	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}

	rowData := make(MappedRow)
	for i, column := range columns {
		rowData[column] = *(values[i].(*any))
	}

	return rowData, nil
}
