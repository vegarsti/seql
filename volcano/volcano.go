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

type constantSelect struct {
	input Node
	i     int
	d     string
}

func (s *constantSelect) Start() {
	s.input.Start()
}

func (s *constantSelect) Next() (Row, bool) {
	for {
		row, ok := s.input.Next()
		if !ok {
			// We've exhausted our input, so we're exhausted too.
			return nil, false
		}

		if row[s.i] == s.d {
			// This row passed the test, so emit it.
			return row, true
		}
	}
}

func ConstantSelect(input Node, i int, d string) Node {
	return &constantSelect{
		input: input,
		i:     i,
		d:     d,
	}
}

type equalsSelect struct {
	input Node
	i, j  int
}

func (s *equalsSelect) Start() {
	s.input.Start()
}

func (s *equalsSelect) Next() (Row, bool) {
	for {
		row, ok := s.input.Next()
		if !ok {
			// We've exhausted our input, so we're exhausted too.
			return nil, false
		}

		if row[s.i] == row[s.j] {
			// This row passed the test, so emit it.
			return row, true
		}
	}
}

func EqualsSelect(input Node, i, j int) Node {
	return &equalsSelect{
		input: input,
		i:     i,
		j:     j,
	}
}

type project struct {
	input Node
	cols  []int
}

func (p *project) Start() {
	p.input.Start()
}

func (p *project) Next() (Row, bool) {
	row, ok := p.input.Next()
	if !ok {
		// We've exhausted our input, so we're exhausted too.
		return nil, false
	}

	newRow := make(Row, len(p.cols))
	for i, col := range p.cols {
		newRow[i] = row[col]
	}
	return newRow, ok
}

func Project(input Node, cols []int) Node {
	return &project{
		input: input,
		cols:  cols,
	}
}
