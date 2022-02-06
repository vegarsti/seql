package spc

import (
	"bytes"
	"fmt"
)

type Row []string

type Relation struct {
	colNames []string
	rows     []Row
}

func NewRelation(colNames []string, rows []Row) Relation {
	return Relation{
		colNames: colNames,
		rows:     rows,
	}
}

func (r Relation) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v\n----\n", r.colNames)
	for _, row := range r.rows {
		fmt.Fprintf(&buf, "%v\n", row)
	}
	return buf.String()
}

// ConstantSelect filters rel to only the rows for which the i-th column is equal
// to d.
func ConstantSelect(rel Relation, i int, d string) Relation {
	// Create a slice to store the new result
	result := make([]Row, 0)
	// Iterate over the old row set, adding a row to the new row set if it passes
	// the test.
	for _, row := range rel.rows {
		if row[i] == d {
			result = append(result, row)
		}
	}
	return Relation{
		// The output has the same colNames as the input.
		colNames: rel.colNames,
		rows:     result,
	}
}

// EqualsSelect filters rel to only the rows for which the i-th column is equal
// to the j-th column
func EqualsSelect(rel Relation, i, j int) Relation {
	// Create a slice to store the new result
	result := make([]Row, 0)
	// Iterate over the old row set, adding a row to the new row set if it passes
	// the test.
	for _, row := range rel.rows {
		if row[i] == row[j] {
			result = append(result, row)
		}
	}
	return Relation{
		// The output has the same colNames as the input.
		colNames: rel.colNames,
		rows:     result,
	}
}
