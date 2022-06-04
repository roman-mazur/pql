package data

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/osquery/osquery-go"
)

type OsQuerySource struct {
	Path string

	osq *osquery.ExtensionManagerClient
}

func NewOsQuerySource(path string) (*OsQuerySource, error) {
	res := &OsQuerySource{Path: path}
	if err := res.connect(); err != nil {
		return nil, err
	}
	return res, nil
}

func (oqs *OsQuerySource) connect() error {
	if oqs.osq != nil {
		return nil
	}
	osq, err := osquery.NewClient(oqs.Path, 10*time.Second)
	if err != nil {
		return err
	}
	oqs.osq = osq
	return nil
}

func (oqs *OsQuerySource) Query(ctx context.Context, query string) (rs Set, err error) {
	defer func() {
		if err != nil {
			_ = oqs.Close()
		}
	}()

	if err := oqs.connect(); err != nil {
		return nil, err
	}

	colResp, err := oqs.osq.GetQueryColumns(query)
	if err != nil {
		return nil, err
	}
	if colResp.Status == nil || colResp.Status.Code != 0 {
		return nil, fmt.Errorf("unable to query columns: %v", colResp.Status)
	}

	rows, err := oqs.osq.QueryRows(query)
	if err != nil {
		return nil, err
	}

	var res osQueryRes
	res.values = rows
	res.l = len(rows)
	res.types = make(map[string]string, len(colResp.Response))
	for _, info := range colResp.Response {
		for name, typ := range info {
			res.columns = append(res.columns, name)
			res.types[name] = typ
			break
		}
	}
	return &res, nil
}

func (oqs *OsQuerySource) Close() error {
	if oqs.osq != nil {
		oqs.osq.Close()
		oqs.osq = nil
	}
	return nil
}

type osQueryRes struct {
	arrayValues
	types  map[string]string
	values []map[string]string
}

func (r *osQueryRes) Scan(dest ...interface{}) error {
	for i, d := range dest {
		cn := r.columns[i]
		sv := r.values[r.c][cn]

		if d == nil {
			continue
		}
		if scan, ok := d.(sql.Scanner); ok {
			err := scan.Scan(sv)
			if err != nil {
				return err
			}
			continue
		}

		switch typ := r.types[cn]; typ {
		case "TEXT":
			if dv, ok := d.(*string); ok {
				*dv = sv
			} else {
				return fmt.Errorf("cannot store TEXT in a non-string, row %d, col %d (%s)", r.c, i, cn)
			}
		default:
			switch dv := d.(type) {
			case *float64:
				var err error
				*dv, err = strconv.ParseFloat(sv, 64)
				if err != nil {
					return fmt.Errorf("cannot parse float value, row %d, col %d (%s)", r.c, i, cn)
				}
			case *string:
				*dv = r.values[r.c][cn]
			default:
				return fmt.Errorf("cannot store value in a %s, row %d, col %d (%s)", reflect.TypeOf(d), r.c, i, cn)
			}
		}
	}
	return nil
}
