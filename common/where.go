package common

// WhereClause is a struct representing a where clause.
// This is made to easily create WHERE clauses from parameters passed from a request.
type WhereClause struct {
	Clause string
	Params []interface{}
}

// Where adds a new WHERE clause to the WhereClause.
func (w *WhereClause) Where(clause, passedParam string, allowedValues ...string) *WhereClause {
	if passedParam == "" {
		return w
	}
	if len(allowedValues) != 0 && !contains(allowedValues, passedParam) {
		return w
	}
	// checks passed, if string is empty add "WHERE"
	if w.Clause == "" {
		w.Clause += "WHERE "
	} else {
		w.Clause += " AND "
	}
	w.Clause += clause
	w.Params = append(w.Params, passedParam)
	return w
}

// Where is the same as WhereClause.Where, but creates a new WhereClause.
func Where(clause, passedParam string, allowedValues ...string) *WhereClause {
	w := new(WhereClause)
	return w.Where(clause, passedParam, allowedValues...)
}
