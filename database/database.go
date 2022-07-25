package database

import (
	"database/sql"
	"fmt"
	"procedure-run/procedure"

	"github.com/desertbit/grumble"
	_ "github.com/go-sql-driver/mysql"
)

type database struct {
	app         *grumble.App
	databaseMap map[string]string
}

func New(app *grumble.App) *database {
	return &database{app: app, databaseMap: make(map[string]string)}
}

func (d *database) addCommands() {
	d.app.SetPrompt("存储过程应用 » ")
	d.app.Commands().RemoveAll()
	d.app.AddCommand(&grumble.Command{
		Name: "create",
		Help: "创建数据库",
		Args: func(a *grumble.Args) {
			a.String("databaseName", "数据库别名")
			a.String("collection", "数据库链接")
		},
		Run: func(c *grumble.Context) error {
			databaseName := c.Args.String("databaseName")
			collection := c.Args.String("collection")
			if _, b := d.databaseMap[databaseName]; !b {
				info, err := validateCollection(collection)
				if err != nil {
					c.App.Println(fmt.Sprintf("%v，错误：%v", info, err))
					return nil
				}
				d.databaseMap[databaseName] = collection
				c.App.Println("创建成功")
			} else {
				c.App.Println("该数据库别名重复")
			}
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
			delete(d.databaseMap, databaseName)
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
			for k, v := range d.databaseMap {
				c.App.Println(fmt.Sprintf("%v", k), v)
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
			if v, b := d.databaseMap[databaseName]; b {
				info, err := validateCollection(v)
				if err != nil {
					c.App.Println(fmt.Sprintf("%v，错误：%v", info, err))
					return nil
				}
				procedure.NewProcedure(c.App, databaseName, func() bool {
					d.Run()
					return true
				}).Run()
			} else {
				c.App.Println("该数据库别名不存在")
			}
			return nil
		},
	})
}

func validateCollection(c string) (string, error) {
	db, err := sql.Open("mysql", c)
	if err != nil {
		//c.App.Println("该数据库连接失败，错误：", err)
		return "该数据库连接失败", err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		//c.App.Println("该数据库连接失败1，错误：", err)
		return "该数据库连接失败", err
	}
	return "", nil
}

func (d *database) Run() {
	d.addCommands()
}
