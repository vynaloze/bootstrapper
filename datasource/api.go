package datasource

type Datasource interface {
	Get(key string) (string, bool)
}

var datasources = make([]Datasource, 0)

func Find(key string) (string, bool) {
	for _, d := range datasources {
		v, ok := d.Get(key)
		if ok {
			return v, ok
		}
	}
	return "", false
}
