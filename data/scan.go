package data

var ignore ignoreScan

type ignoreScan struct{}

func (is ignoreScan) Scan(src interface{}) error {
	return nil
}

func ScanArgs(dv *float64, dl *string, vi, li, size int) []interface{} {
	args := make([]interface{}, size)
	for i := range args {
		switch i {
		case vi:
			args[i] = dv
		case li:
			args[i] = dl
		default:
			args[i] = ignore
		}
	}
	return args
}

func Indexes(columns []string, x, y string) (vi, li int) {
	li, vi = 0, 1
	if x != "" || y != "" {
		for i := range columns {
			switch columns[i] {
			case x:
				li = i
			case y:
				vi = i
			}
		}
	}
	if len(columns) == 1 {
		vi = 0
	}
	return
}
