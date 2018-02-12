package utils

// Paginator : paginator struct
type Paginator struct {
	Pagenum int
	Limit   int
}

// ToParameter return start and limit from Paginator
func (pa *Paginator) ToParameter() (int, int) {
	limit := pa.Limit
	start := (pa.Pagenum - 1) * limit
	if start < 0 {
		start = 0
	}
	return start, limit
}
