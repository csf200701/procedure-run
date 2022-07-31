package database

import (
	"fmt"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"
)

type QuestionType int

const (
	Type QuestionType = iota
	Host
	Port
	User
	Password
	Db
	OracleDb
	SSHHost
	SSHPort
	SSHUser
	SSHPassword
)

var questionMap = map[QuestionType]*survey.Question{
	Type: {
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
	Host: {
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
	Port: {
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
	User: {
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
	Password: {
		Name:   "password",
		Prompt: &survey.Password{Message: "密码?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("密码不能为空")
			}
			return nil
		},
	},
	Db: {
		Name:   "db",
		Prompt: &survey.Input{Message: "数据库名?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("数据库名不能为空")
			}
			return nil
		},
	},
	OracleDb: {
		Name:   "db",
		Prompt: &survey.Input{Message: "数据库名?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("数据库名不能为空")
			}
			return nil
		},
	},
	SSHHost: {
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
	SSHPort: {
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
	SSHUser: {
		Name:   "ssh_user",
		Prompt: &survey.Input{Message: "SSH用户名（root）?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return nil
			}
			return nil
		},
	},
	SSHPassword: {
		Name:   "ssh_password",
		Prompt: &survey.Password{Message: "SSH密码?"},
		Validate: func(ans interface{}) error {
			return nil
		},
	},
}

func Questions(qtypes ...QuestionType) []*survey.Question {
	l := len(qtypes)
	questions := make([]*survey.Question, 0, l)
	if l > 0 {
		for i := 0; i < l; i++ {
			questions = append(questions, questionMap[qtypes[i]])
		}
	}
	return questions
}
