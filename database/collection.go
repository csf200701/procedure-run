package database

import (
	"fmt"
	"net"
	"procedure-run/connector"
	"strconv"

	// "gopkg.in/AlecAivazis/survey.v1"
	survey "github.com/AlecAivazis/survey/v2"
)

var qs = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "数据库类型（Mysql）?",
			Options: []string{"Mysql", "PostgreSQL", "Oracle"},
			Default: "Mysql",
		},
		Validate: func(ans interface{}) error {
			return nil
		},
	},
	{
		Name:   "host",
		Prompt: &survey.Input{Message: "主机地址（127.0.0.1）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			if str, ok := ans.(string); ok {
				_, i := ParseIP(str)
				if i == 0 {
					return fmt.Errorf("主机地址格式不正确")
				}
			} else {

				return fmt.Errorf("主机地址格式不正确")
			}
			return nil
		},
		Transform: survey.Title,
	},
	{
		Name:   "port",
		Prompt: &survey.Input{Message: "端口号（3306）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			if str, ok := ans.(string); ok {
				_, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return fmt.Errorf("端口号格式不正确")
				}
			} else {

				return fmt.Errorf("端口号格式不正确")
			}

			return nil
		},
		Transform: survey.Title,
	},
	{
		Name:   "user",
		Prompt: &survey.Input{Message: "用户名（root）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			return nil
		},
		//Transform: survey.ToLower,
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "密码?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("密码不能为空")
			}
			return nil
		},
	},
	{
		Name:   "db",
		Prompt: &survey.Input{Message: "数据库名?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("数据库名不能为空")
			}
			return nil
		},
	},
}

var qsssh = []*survey.Question{
	{
		Name:   "ssh_host",
		Prompt: &survey.Input{Message: "SSH主机地址（127.0.0.1）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			if str, ok := ans.(string); ok {
				_, i := ParseIP(str)
				if i == 0 {
					return fmt.Errorf("主机地址格式不正确")
				}
			} else {

				return fmt.Errorf("主机地址格式不正确")
			}
			return nil
		},
		Transform: survey.Title,
	},
	{
		Name:   "ssh_port",
		Prompt: &survey.Input{Message: "SSH端口号（22）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			if str, ok := ans.(string); ok {
				_, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return fmt.Errorf("端口号格式不正确")
				}
			} else {

				return fmt.Errorf("端口号格式不正确")
			}

			return nil
		},
		Transform: survey.Title,
	},
	{
		Name:   "ssh_user",
		Prompt: &survey.Input{Message: "SSH用户名（root）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			return nil
		},
	},
	{
		Name:   "ssh_password",
		Prompt: &survey.Password{Message: "SSH密码?"},
		Validate: func(ans interface{}) error {
			return nil
		},
	},
}

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
	// var newqs []*survey.Question
	// if ssh {
	// 	newqs = append(qs, qsssh...)
	// } else {
	// 	newqs = qs
	// }
	survey.Ask(qs, collection)
	if ssh {
		survey.Ask(qsssh, collection)
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
