package database

import (
	"net"
	"procedure-run/connector"

	// "gopkg.in/AlecAivazis/survey.v1"
	survey "github.com/AlecAivazis/survey/v2"
)

func init() {
	survey.SelectQuestionTemplate = `
	{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
	{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
	{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
	{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
	{{- else}}
		{{- "  "}}{{- color "cyan"}}[使用箭头移动、空格选择、键入筛选{{- if and .Help (not .ShowHelp)}}, {{ .Config.HelpInput }} for more help{{end}}]{{color "reset"}}
		{{- "\n"}}
		{{- range $ix, $choice := .PageEntries}}
		{{- if eq $ix $.SelectedIndex }}{{color $.Config.Icons.SelectFocus.Format }}{{ $.Config.Icons.SelectFocus.Text }} {{else}}{{color "default"}}  {{end}}
		{{- $choice.Value}}
		{{- color "reset"}}{{"\n"}}
		{{- end}}
	{{- end}}`

	survey.ErrorTemplate =
		`{{color .Icon.Format }}{{ .Icon.Text }} 对不起, 校验失败: {{ .Error.Error }}{{color "reset"}}
`
}

func Ask(ssh bool) *connector.Collection {
	collection := new(connector.Collection)
	survey.Ask(Questions(Type), collection)
	if collection.DbType == connector.ORACLE {
		survey.Ask(Questions(Host, Port, User, Password, OracleDb), collection)
	} else {
		survey.Ask(Questions(Host, Port, User, Password, Db), collection)
	}
	if ssh {
		survey.Ask(Questions(SSHHost, SSHPort, SSHUser, SSHPassword), collection)
	}
	return collection
}

func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}
