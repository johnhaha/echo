package echo

type SubCtx struct {
	Value
}

func (c *SubCtx) Parser(data interface{}) error {
	err := c.GetJsonData(data)
	return err
}

type JobHandler func(*SubCtx)

const (
	BoolTrue  = "True"
	BoolFalse = "False"
)
