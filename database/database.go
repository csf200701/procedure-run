package database

import (
	"database/sql"
	"fmt"
	"procedure-run/procedure"
	"procedure-run/store"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
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
		Args: func(a *grumble.Args) {
			a.String("databaseName", "数据库别名")
		},
		Run: func(c *grumble.Context) error {
			collection := Ask()
			databaseName := c.Args.String("databaseName")
			_, err := store.GetDatabase(databaseName)
			if err == nil {
				c.App.Config().ErrorColor.Println("该数据库别名重复")
				return nil
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
			collectionStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", collection.User, collection.Password, collection.Host, collection.Port, collection.DbName)
			info, err := validateCollection(collectionStr)
			if err != nil {
				c.App.Config().ErrorColor.Println(fmt.Sprintf("%v，错误：%v", info, err))
				return nil
			}
			store.AddDatabase(databaseName, collectionStr)
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
				table, err := gotable.Create("数据库别名", "数据库连接串")
				if err != nil {
					return nil
				}
				for _, database := range databases {
					table.AddRow([]string{database.DatabaseAlias, database.DatabaseCollection})
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
				info, err := validateCollection(database.DatabaseCollection)
				if err != nil {
					c.App.Config().ErrorColor.Println(fmt.Sprintf("%v，错误：%v", info, err))
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

func validateCollection(c string) (string, error) {
	db, err := sql.Open("mysql", c)
	if err != nil {
		return "该数据库连接失败", err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return "该数据库连接失败", err
	}
	return "", nil
}

func (d *database) Run() {
	d.addCommands()
}
