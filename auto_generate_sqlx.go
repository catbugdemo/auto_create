package auto

import (
	"bytes"
	"fmt"
	"github.com/catbugdemo/errors"
	"strings"
	"text/template"
)

func AutoGenerateSqlx(way Way) (string, error) {
	if err := way.init(); err != nil {
		return "", errors.WithStack(err)
	}

	json, err := way.formatJSONSQL()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return json, nil
}

func (n *Normal) formatJSONSQL() (string, error) {
	columns := FindColumns(n.Driver, n.DataSource, n.TableName)
	var str = `
package {{.package}}

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
	"time"
)

type {{.struct_name}} struct {
	${type_struct}
}

func (m *{{.struct_name}}) Table() string {
	return "{{.table_name}}"
}

func (m *{{.struct_name}}) Condition(arg ...string) string {
${auto_where}
}

func (m *{{.struct_name}}) Count(db *sqlx.DB, arg ...string) (int, error) {
	return utils.Count(db, m, arg...)
}

func (m *{{.struct_name}}) Insert(db *sqlx.DB) error {
	return utils.Insert(db, m)
}

func (m *{{.struct_name}}) Find(db *sqlx.DB, arg ...string) ([]{{.struct_name}}, error) {
	var list []{{.struct_name}}
	if err := utils.Find(db, &list, arg...); err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

func (m *{{.struct_name}}) First(db *sqlx.DB, arg ...string) error {
	// limit , offset
	arg = append(arg, "1")
	find, err := m.Find(db, arg...)
	if err != nil {
		return errors.WithStack(err)
	}
	*m = find[0]
	return nil
}

func (m *{{.struct_name}}) Update(db *sqlx.DB, arg ...string) error {
	return utils.Update(db, m, arg...)
}

func (m *{{.struct_name}}) Delete(db *sqlx.DB, arg ...string) error {
	return utils.Delete(db, m, arg...)
}

func (m *{{.struct_name}}) IfExist(db *sqlx.DB, arg string) error {
	count, err := m.Count(db, arg)
	if err != nil {
		return errors.WithStack(err)
	}
	if count == 0 {
		return utils.ErrRecordNotFound
	}
	return nil
}

func (m *{{.struct_name}}) FindByCount(db *sqlx.DB, arg ...string) ([]SmbVocherDetailLog, int, error) {
	var conditon string
	if len(arg) > 0 {
		conditon = arg[0]
	}
	count, err := m.Count(db, conditon)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	if count == 0 {
		return nil, 0, nil
	}
	list, err := m.Find(db, arg...)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func (m *{{.struct_name}}) FirstById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}
	return m.First(db, fmt.Sprintf("where id=%d", id))
}

func (m *{{.struct_name}}) UpdateById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}
	return m.Update(db, fmt.Sprintf("where id=%d", id))
}

func (m *{{.struct_name}}) DeleteById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}
	return m.Delete(db, fmt.Sprintf("where id=%d", id))
}

`
	str = strings.ReplaceAll(str, "${type_struct}", getTypeStruct(columns))
	str = strings.ReplaceAll(str, "${auto_where}", n.autoWhere(columns))
	tt := template.Must(template.New("model").Parse(str))
	vals := map[string]string{
		"package":     n.Package,
		"struct_name": n.Info["struct_name"].(string),
		"table_name":  n.TableName,
	}
	var buf bytes.Buffer
	if err := tt.Execute(&buf, vals); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 自动生成条件
func (n Normal) autoWhere(columns []Column) string {
	var tmp strings.Builder
	query := `
	var params []string 
	params = append(params,arg...)
`
	tmp.WriteString(query)

	for _, column := range columns {
		if in(column.ColumnType, []string{"date", "datetime", "timestamp", "timestamp with time zone"}) {
			continue
		}
		fmt.Print()
		var str = `
	if m.${value_name} != ${type} {
		params = append(params,fmt.Sprint("${tmp_name}='",m.${value_name},"'"))
	}
`
		str = strings.ReplaceAll(str, "${value_name}", underLineToHump(column.ColumnName))
		str = strings.ReplaceAll(str, "${type}", CheckType(column.ColumnType))
		str = strings.ReplaceAll(str, "${tmp_name}", column.ColumnName)
		tmp.WriteString(str)
	}

	var result = `
	if len(params) == 0 {
		return ""
	}
	return fmt.Sprintf("where %s",strings.Join(params," and "))
`
	tmp.WriteString(result)

	return tmp.String()
}

func CheckType(s string) string {
	// postgres
	{
		if in(s, []string{"double precision", "double"}) {
			return "0"
		}
		if in(s, []string{"bigint", "bigserial", "integer", "smallint", "serial", "big serial"}) {
			return "0"
		}
		if in(s, []string{"character varying", "varchar"}) {
			return "\"\""
		}
		if in(s, []string{"bool", "boolean"}) {
			return "true"
		}
		if in(s, []string{"bytea", "jsonb"}) {
			return "nil"
		}
	}
	return ""
}
