package volcano

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

type Node interface {
	// Start is called to initialize any state that this node needs to execute.
	Start()

	// Next returns the next row in the Node's result set. If there are no more
	// rows to return, the second return value will be false, otherwise, it will
	// be true.
	Next() (Row, bool)
}

type scanRelation struct {
	data Relation

	// idx points at which row we are to return next
	idx int
}

// Start is a no-op for the scanRelation; no initalization is needed.
func (s *scanRelation) Start() {}

// Next returns the next row.
func (s *scanRelation) Next() (Row, bool) {
	if s.idx >= len(s.data.rows) {
		// We've read all the data
		return nil, false
	}
	s.idx++
	return s.data.rows[s.idx-1], true
}

func ScanRelation(data Relation) Node {
	return &scanRelation{
		data: data,
		idx:  0,
	}
}
