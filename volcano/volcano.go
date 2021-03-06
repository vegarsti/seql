package volcano

import (
	"bytes"
	"fmt"
	"sort"
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

type cross struct {
	left  Node
	right Node
	// leftBuffer contains all the rows from left
	leftBuffer []Row
	// row is the current row from right
	row Row
	// idx is the current pointer into leftBuffer
	idx int
}

func (c *cross) Start() {
	c.left.Start()
	c.right.Start()
	// Buffer up everything in left
	for row, ok := c.left.Next(); ok; row, ok = c.left.Next() {
		c.leftBuffer = append(c.leftBuffer, row)
	}
	// Start out as though we finished an iteratio
	c.idx = len(c.leftBuffer)
}

func Cross(left Node, right Node) Node {
	return &cross{
		left:       left,
		right:      right,
		leftBuffer: make([]Row, 0),
		row:        nil,
		idx:        0,
	}
}

func (c *cross) Next() (Row, bool) {
	// If we're done with left, reset it and grab a new row from right.
	// This is in a loop to neatly handle the case where left has zero rows.
	for c.idx >= len(c.leftBuffer) {
		row, ok := c.right.Next()
		if !ok {
			return nil, false
		}
		c.row = row
		c.idx = 0
	}
	leftRow := c.leftBuffer[c.idx]
	rightRow := c.row
	c.idx++
	return append(append(make(Row, 0), leftRow...), rightRow...), true
}

type union struct {
	left  Node
	right Node
}

func Union(left Node, right Node) Node {
	return &union{
		left:  left,
		right: right,
	}
}

func (u *union) Start() {
	u.left.Start()
	u.right.Start()
}

func (u *union) Next() (Row, bool) {
	row, ok := u.left.Next()
	if !ok {
		row, ok = u.right.Next()
		if !ok {
			return nil, false
		}
	}
	return row, true
}

type zip struct {
	left  Node
	right Node
}

func Zip(left Node, right Node) Node {
	return &zip{
		left:  left,
		right: right,
	}
}

func (z *zip) Start() {
	z.left.Start()
	z.right.Start()
}

func (z *zip) Next() (Row, bool) {
	leftRow, ok := z.left.Next()
	if !ok {
		return nil, false
	}
	rightRow, ok := z.right.Next()
	if !ok {
		return nil, false
	}
	return append(append(make(Row, 0), leftRow...), rightRow...), true
}

type inspect struct {
	node Node
}

func Inspect(node Node) Node {
	return &inspect{
		node: node,
	}
}

func (i *inspect) Start() {
	i.node.Start()
}

func (i *inspect) Next() (Row, bool) {
	row, ok := i.node.Next()
	if !ok {
		return nil, false
	}
	fmt.Println(row)
	return row, true
}

type intersect struct {
	left  Node
	right Node
	// leftBuffer contains all the rows from left
	leftBuffer []Row
}

func Intersect(left Node, right Node) Node {
	return &intersect{
		left:       left,
		right:      right,
		leftBuffer: make([]Row, 0),
	}
}

func (i *intersect) Start() {
	i.left.Start()
	i.right.Start()
	// Buffer up everything in left
	for row, ok := i.left.Next(); ok; row, ok = i.left.Next() {
		i.leftBuffer = append(i.leftBuffer, row)
	}
}

func (i *intersect) Next() (Row, bool) {
	for rightRow, ok := i.right.Next(); ok; rightRow, ok = i.right.Next() {
		for _, leftRow := range i.leftBuffer {
			emit := true
			for i := range leftRow {
				if leftRow[i] != rightRow[i] {
					emit = false
					break
				}
			}
			if emit {
				return rightRow, true
			}
		}
	}
	return nil, false
}

type distinct struct {
	node Node
	seen []Row
}

func Distinct(node Node) Node {
	return &distinct{
		node: node,
		seen: make([]Row, 0),
	}
}

func (d *distinct) Start() {
	d.node.Start()
}

func (d *distinct) Next() (Row, bool) {
	for row, ok := d.node.Next(); ok; row, ok = d.node.Next() {
		emit := true
		for _, seenRow := range d.seen {
			allEqual := true
			for i := range seenRow {
				if seenRow[i] != row[i] {
					allEqual = false
					break
				}
			}
			if allEqual {
				emit = false
			}
		}
		if emit {
			d.seen = append(d.seen, row)
			return row, true
		}
	}
	return nil, false
}

type order struct {
	node   Node
	byCols []int
	buffer []Row // holds all rows from node
	idx    int   // index into sorted buffer
}

func Order(node Node, byCols []int) Node {
	return &order{
		node:   node,
		byCols: byCols,
		buffer: make([]Row, 0),
	}
}

func (o *order) Start() {
	o.node.Start()
	// Exhaust the node so we can sort it
	for row, ok := o.node.Next(); ok; row, ok = o.node.Next() {
		o.buffer = append(o.buffer, row)
	}
	// Sort buffer
	sort.Slice(o.buffer, func(i, j int) bool {
		for _, col := range o.byCols {
			// Rows are equal in this column, use the next column
			if o.buffer[i][col] == o.buffer[j][col] {
				continue
			}
			return o.buffer[i][col] < o.buffer[j][col]
		}
		// No columns left in which the rows are different, so just return true
		return true
	})
}

func (o *order) Next() (Row, bool) {
	if o.idx >= len(o.buffer) {
		// No rows left
		return nil, false
	}
	o.idx++
	return o.buffer[o.idx-1], true
}

type hashJoin struct {
	left         Node // probe input
	right        Node // build input
	index        map[string][]Row
	leftJoinKey  int
	rightJoinKey int
	// leftRow is the current row from left
	leftRow Row
	// idx is the current pointer into the []Row slice in index[row[leftJoinKey]]
	idx int
}

func HashJoin(left Node, right Node, leftJoinKey int, rightJoinKey int) Node {
	return &hashJoin{
		left:         left,
		right:        right,
		index:        make(map[string][]Row),
		leftJoinKey:  leftJoinKey,
		rightJoinKey: rightJoinKey,
		leftRow:      nil,
		idx:          0,
	}
}

func (h *hashJoin) Start() {
	h.left.Start()
	h.right.Start()

	// Build hash table from build input (right)
	for row, ok := h.right.Next(); ok; row, ok = h.right.Next() {
		h.index[row[h.rightJoinKey]] = append(h.index[row[h.rightJoinKey]], row)
	}
	// Initialize h.leftRow to first row in left (possibly nil)
	row, _ := h.left.Next()
	h.leftRow = row
}

func (h *hashJoin) Next() (Row, bool) {
	for h.leftRow != nil {
		joinValue := h.leftRow[h.leftJoinKey]
		// there is a row to emit; emit it
		if h.idx < len(h.index[joinValue]) {
			rightRow := h.index[joinValue][h.idx]
			h.idx++
			return append(append(make(Row, 0), h.leftRow...), rightRow...), true
		}
		// no more rows to emit from join row; get next and reset index pointer
		row, _ := h.left.Next()
		h.leftRow = row
		h.idx = 0
	}
	return nil, false
}
