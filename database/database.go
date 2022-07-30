package database

import (
	"fmt"
	"procedure-run/connector"
	"procedure-run/procedure"
	"procedure-run/store"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/liushuochen/gotable"
)

type database struct {
	app *grumble.App
}

func New(app *grumble.App) *database {
	return &database{app: app}
}

func (d *database) addCommands() {
	d.app.SetPrompt("存储过程应用 » ")
	d.app.Commands().RemoveAll()
	d.app.AddCommand(&grumble.Command{
		Name: "create",
		Help: "创建数据库",
		Flags: func(f *grumble.Flags) {
			f.Bool("s", "ssh", false, "是否启用SSH连接")
		},
		Args: func(a *grumble.Args) {
			a.String("databaseName", "数据库别名")
		},
		Run: func(c *grumble.Context) error {
			databaseName := c.Args.String("databaseName")
			_, err := store.GetDatabase(databaseName)
			if err == nil {
				c.App.Config().ErrorColor.Println("该数据库别名重复")
				return nil
			}
			ssh := c.Flags.Bool("ssh")
			collection := Ask(ssh)
			collection.IsSSH = ssh
			if collection.DbType == "" {
				collection.DbType = connector.MYSQL
			}
			if collection.Host == "" {
				collection.Host = "127.0.0.1"
			}
			if collection.Port == "" {
				collection.Port = "3306"
			}
			if collection.User == "" {
				collection.User = "root"
			}
			if collection.Password == "" {
				return nil
			}
			if collection.DbName == "" {
				return nil
			}
			if collection.IsSSH {
				if collection.SSHHost == "" {
					collection.SSHHost = "127.0.0.1"
				}
				if collection.SSHPort == "" {
					collection.SSHPort = "22"
				}
				if collection.SSHUser == "" {
					collection.SSHUser = "root"
				}
			}

			connector := connector.Database(collection)
			if connector == nil {
				c.App.Config().ErrorColor.Println(collection.DbType + "该类型数据库创建失败")
				return nil
			}
			err = connector.ValidateCollection()
			if err != nil {
				c.App.Config().ErrorColor.Println(err.Error())
				return nil
			}
			store.AddDatabase(databaseName, collection)
			c.App.Println("创建成功")

			return nil
		},
	})
	d.app.AddCommand(&grumble.Command{
		Name: "delete",
		Help: "删除数据库",
		Args: func(a *grumble.Args) {
			a.String("databaseName", "数据库名称")
		},
		Run: func(c *grumble.Context) error {
			databaseName := c.Args.String("databaseName")
			//delete(d.databaseMap, databaseName)
			store.DeleteDatabase(databaseName)
			c.App.Println("删除成功")
			return nil
		},
	})
	d.app.AddCommand(&grumble.Command{
		Name: "ls",
		Help: "列出数据库",
		Args: func(a *grumble.Args) {

		},
		Run: func(c *grumble.Context) error {
			databases := store.FindAll()
			l := len(databases)
			if l > 0 {
				table, err := gotable.Create("数据库别名", "数据库类型", "地址", "端口", "是否SSH")
				if err != nil {
					return nil
				}
				for _, database := range databases {
					isSSHName := "否"
					if database.DatabaseCollection.IsSSH {
						isSSHName = "是"
					}
					table.AddRow([]string{database.DatabaseAlias,
						string(database.DatabaseCollection.DbType),
						database.DatabaseCollection.Host,
						database.DatabaseCollection.Port,
						isSSHName})
				}
				fmt.Println(table)
			} else {
				color.New(color.FgGreen, color.Bold).Println("空")
			}
			return nil
		},
	})
	d.app.AddCommand(&grumble.Command{
		Name: "use",
		Help: "使用数据库",
		Args: func(a *grumble.Args) {
			a.String("databaseName", "数据库名称")
		},
		Run: func(c *grumble.Context) error {
			databaseName := c.Args.String("databaseName")
			database, err := store.GetDatabase(databaseName)
			if err == nil {
				connector := connector.Database(database.DatabaseCollection)
				err = connector.ValidateCollection()
				if err != nil {
					c.App.Config().ErrorColor.Println(err.Error())
					return nil
				}
				procedure.NewProcedure(c.App, databaseName, database.DatabaseCollection, func() bool {
					d.Run()
					return true
				}).Run()
			} else {
				c.App.Config().ErrorColor.Println("该数据库别名不存在")
			}
			return nil
		},
	})
}

func (d *database) Run() {
	d.addCommands()
}
