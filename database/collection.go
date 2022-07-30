package database

import (
	"fmt"
	"net"
	"strconv"

	"gopkg.in/AlecAivazis/survey.v1"
	surveycore "gopkg.in/AlecAivazis/survey.v1/core"
)

type Collection struct {
	Host     string `survey:"host"`
	Port     string `survey:"port"`
	User     string `survey:"user"`
	Password string `survey:"password"`
	DbName   string `survey:"db"`
}

var qs = []*survey.Question{
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
		// Prompt: &survey.Select{
		// 	Message: "Choose a color:",
		// 	Options: []string{"red", "blue", "green"},
		// 	Default: "red",
		// },
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
		Prompt: &survey.Input{Message: "密码?"},
		Validate: func(ans interface{}) error {
			if err := survey.Required(ans); err != nil {
				return fmt.Errorf("密码不能为空")
			}
			return nil
		},
		//Transform: survey.Title,
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
		//Transform: survey.Title,
	},
}

func Ask() *Collection {
	surveycore.ErrorTemplate = `{{color "red"}}{{ ErrorIcon }} 对不起, 校验失败: {{.Error}}{{color "reset"}}
	`
	collection := new(Collection)
	survey.Ask(qs, collection)
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
