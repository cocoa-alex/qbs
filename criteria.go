package qbs

type Criteria struct {
	model     *Model
	condition *Condition
	orderBy   string
	orderDesc bool
	limit     int
	offset    int
	omitFields []string
}


func (c *Criteria) mergePkCondition(d Dialect) {
	idValue := c.model.idValue()
	var con *Condition
	if idValue != 0 {
		expr := d.Quote(c.model.Pk.Name) + " = ?"
		con = NewCondition(expr, idValue)
		con.AndCondition(c.condition)
	} else {
		con = c.condition
	}
	c.condition = con
}

// Conditions are structured in a way to define
// complex where clause easily.
type Condition struct {
	Expr string
	Args []interface{}
	Sub  *Condition
	IsOr bool
}

func NewCondition(expr string, args ...interface{}) *Condition {
	return &Condition{
		Expr: expr,
		Args: args,
	}
}

func NewInCondition(column string, values []interface {}) *Condition {
	expr := column + " IN ("
	for _ = range values {
		expr +="?, "
	}
	expr = expr[:len(expr) - 2]
	expr += ")"
	return &Condition{
		Expr: expr,
		Args: values,
	}
}

func (c *Condition) And(expr string, args ...interface{}) *Condition {
	if c.Sub != nil {
		c.Expr, c.Args = c.Merge()
	}
	c.Sub = NewCondition(expr, args...)
	c.IsOr = false
	return c
}

func (c *Condition) AndCondition(subCondition *Condition) *Condition {
	if c.Sub != nil {
		c.Expr, c.Args = c.Merge()
	}
	c.Sub = subCondition
	c.IsOr = false
	return c
}

func (c *Condition) Or(expr string, args ...interface{}) *Condition {
	if c.Sub != nil {
		c.Expr, c.Args = c.Merge()
	}
	c.Sub = NewCondition(expr, args...)
	c.IsOr = true
	return c
}

func (c *Condition) OrCondition(subCondition *Condition) *Condition {
	if c.Sub != nil {
		c.Expr, c.Args = c.Merge()
	}
	c.Sub = subCondition
	c.IsOr = true
	return c
}

func (c *Condition) Merge() (expr string, args []interface{}) {
	expr = c.Expr
	args = c.Args
	if c.Sub == nil {
		return
	}
	expr = "(" + expr + ")"
	if c.IsOr {
		expr += " OR "
	} else {
		expr += " AND "
	}
	subExpr, subArgs := c.Sub.Merge()
	expr += "(" + subExpr + ")"
	args = append(args, subArgs...)
	return expr, args
}
