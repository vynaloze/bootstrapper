package datasource

type Literal map[string]string

func NewLiteral() (Literal, error) {
	l := make(Literal)
	datasources = append(datasources, l)
	return l, nil
}

func (l Literal) Get(key string) (string, bool) {
	v, ok := l[key]
	return v, ok
}
