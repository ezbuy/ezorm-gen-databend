package databend

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/ezbuy/ezorm-gen-databend/internal/handler"
	"github.com/ezbuy/ezorm/v2/pkg/plugin"
	"github.com/iancoleman/strcase"
)

//go:embed templates/*.tpl
var fsTemplate embed.FS

var files = []string{
	"templates/create_table.tpl",
}

var (
	_ handler.SchemaHandler = (*CreateTable)(nil)
	_ handler.Printer       = (*CreateTable)(nil)
)

type DBField struct {
	name       string
	colType    string
	isNull     bool
	colDefault string
}

func (f DBField) GetName() string {
	return f.name
}

func (f DBField) GetType() string {
	return f.colType
}

func (f DBField) GetNull() string {
	if f.isNull {
		return "NULL"
	}
	return "NOT NULL"
}

func (f DBField) GetDefault() string {
	if f.colDefault != "" {
		return "DEFAULT " + f.colDefault
	}
	return ""
}

type CreateTable struct {
	Database string
	Table    string
	Fields   []DBField
	fileName string
}

func getDB(s plugin.Schema) string {
	return s["dbname"].(string)
}

func getTable(s plugin.Schema) string {
	return s["dbtable"].(string)
}

func getFields(s plugin.Schema) ([]DBField, error) {
	rawFields := s["fields"].([]any)
	resFields := make([]DBField, 0, len(rawFields))
	for _, f := range rawFields {
		var df DBField
		ff, ok := f.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("field is not map[string]any, raw type is %#+v", f)
		}
		for k, v := range ff {
			switch k {
			case "sqltype":
				df.colType = toBendSQL(v.(string))
			case "sqlcolumn":
				df.name = v.(string)
			case "flags":
				for _, flag := range v.([]any) {
					fstr := flag.(string)
					if fstr == "nullable" {
						df.isNull = true
					}
				}
			case "default":
				df.colDefault = v.(string)
			default:
				if k[:1] == strings.ToUpper(k[:1]) {
					df.name = strcase.ToSnake(k)
					df.colType = toBendSQL(v.(string))
				}
			}
		}
		resFields = append(resFields, df)
	}
	return resFields, nil
}

func toBendSQL(ormType string) string {
	switch ormType {
	case "uint8", "int8":
		return "TINYINT"
	case "uint16", "int16":
		return "SMALLINT"
	case "uint32", "int32", "int":
		return "INT"
	case "uint64", "int64":
		return "BIGINT"
	case "float32":
		return "FLOAT"
	case "float64":
		return "double"
	case "time.Time", "*time.Time", "timestamp", "timeint":
		return "TIMESTAMP"
	case "datetime":
		return "DATE"
	case "bool":
		return "BOOLEAN"
	default:
		return "VARCHAR"
	}
}

func (ct *CreateTable) Handle(ctx context.Context, s plugin.Schema) error {
	ct.Database = getDB(s)
	ct.Table = getTable(s)
	var err error
	ct.Fields, err = getFields(s)
	if err != nil {
		return err
	}
	return nil
}

func (ct *CreateTable) Print(ctx context.Context, dir string) error {
	t, err := template.New("databend_create_table").Funcs(
		sprig.TxtFuncMap(),
	).ParseFS(fsTemplate, files...)
	if err != nil {
		return err
	}
	ct.fileName = filepath.Join(dir, fmt.Sprintf("%s_%s_create_table.sql", ct.Database, ct.Table))
	w, err := os.OpenFile(ct.fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer w.Close()
	return t.ExecuteTemplate(w, "databend_create_table", ct)
}
