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

// Project restrict rel to only the columns in cols
func Project(rel Relation, cols []int) Relation {
	// Restrict each row to cols.
	result := make([]Row, 0)
	for _, row := range rel.rows {
		newRow := make(Row, len(cols))
		for j, idx := range cols {
			newRow[j] = row[idx]
		}
		result = append(result, newRow)
	}
	// Compute the new set of colNames
	colNames := make([]string, len(cols))
	for i, idx := range cols {
		colNames[i] = rel.colNames[idx]
	}
	return Relation{
		colNames: colNames,
		rows:     result,
	}
}

// Cross returns a row for every pair of rows in rel1 and rel2.
func Cross(rel1, rel2 Relation) Relation {
	// Create a slice to store the new result, initialized to be length 0,
	// and with capacity as len(rel1) * rel(2)
	result := make([]Row, 0, len(rel1.colNames)*len(rel2.colNames))
	for _, row1 := range rel1.rows {
		for _, row2 := range rel2.rows {
			newRow := append(append(make(Row, 0), row1...), row2...)
			result = append(result, newRow)
		}
	}
	colNames := append(append(make([]string, 0), rel1.colNames...), rel2.colNames...)
	return Relation{
		colNames: colNames,
		rows:     result,
	}
}
