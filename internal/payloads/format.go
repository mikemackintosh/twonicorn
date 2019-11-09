package payloads

import (
	"fmt"
	"html/template"
	"strings"
)

var transformTemplates = template.FuncMap{
	"join":           strings.Join,
	"formatDatetime": formatDatetime,
	"toJiraTable":    toJiraTable,
	"toJiraPanel":    toJiraPanel,
}

func formatDatetime(date string) string {
	return strings.Replace(date, " ", "T", -1)
}

func toJiraTable(data map[string]interface{}) string {
	c := make([]string, 0, len(data))
	for k, v := range data {
		c = append(c, fmt.Sprintf("||%s|%v|", k, v))
	}
	return strings.Join(c, "\n")
}

func toJiraPanel(title, data interface{}) string {
	return fmt.Sprintf("{panel:title=%s}%#v{panel}", title, data)
}
