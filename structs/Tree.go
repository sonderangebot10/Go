package structs

type Leaf struct {
	Node       Node
	TotalPrice float64
	Totals     TotalReq
	Parent     *Leaf
	Children   []*Leaf
}
