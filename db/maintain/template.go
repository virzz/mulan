package maintain

import (
	"bytes"
	_ "embed"
	"text/template"
)

var (
	// embed template/pgsql_drop.sql
	dropPgSQL string
	// embed template/pgsql_create.sql
	createPgSQL string
	// embed template/mysql_drop.sql
	dropMySQL string
	// embed template/mysql_create.sql
	createMySQL string
)

type SqlTemplateData struct {
	Name   string
	User   string
	Pass   string
	Schema string
}

func parseSql(content string, data SqlTemplateData) (string, error) {
	buf := bytes.Buffer{}
	t, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
