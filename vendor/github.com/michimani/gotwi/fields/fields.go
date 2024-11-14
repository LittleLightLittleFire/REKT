package fields

import "github.com/michimani/gotwi/internal/util"

type Field interface {
	String() string
}

type Fields interface {
	FieldsName() string
	Values() []string
}

func SetFieldsParams(m map[string]string, flist ...Fields) map[string]string {
	for _, f := range flist {
		if f == nil {
			continue
		}

		values := f.Values()
		if len(values) == 0 {
			continue
		}

		m[f.FieldsName()] = util.QueryValue(values)
	}

	return m
}
