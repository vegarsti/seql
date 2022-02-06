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
