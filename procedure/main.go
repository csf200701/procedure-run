package procedure

import (
	"database/sql"
	"fmt"

	"github.com/desertbit/grumble"
)

type procedure struct {
	app           *grumble.App
	databaseAlias string
	collection    string
	quit          func() bool
}

func NewProcedure(app *grumble.App, databaseAlias string, collection string, quit func() bool) *procedure {
	return &procedure{app: app, databaseAlias: databaseAlias, collection: collection, quit: quit}
}

func (p *procedure) addCommands() {
	p.app.SetPrompt(fmt.Sprintf("存储过程应用 » %v：", p.databaseAlias))
	p.app.Commands().RemoveAll()
	p.app.AddCommand(&grumble.Command{
		Name: "call",
		Help: "调用存储过程",
		Args: func(a *grumble.Args) {
			a.String("name", "存储过程名称")
			//a.StringList("values", "")
			a.StringList("values", "参数列表")
		},
		Run: func(c *grumble.Context) error {
			name := c.Args.String("name")
			values := c.Args.StringList("values")
			fmt.Println("name：", name)
			fmt.Println("values：", values)

			db, err := sql.Open("mysql", p.collection)
			if err != nil {
				c.App.Println("该数据库连接失败，错误：", err)
				return nil
			}
			db.Close()
			var sqlStr = "call %v("
			var vList []interface{} = make([]interface{}, 0)
			for i, v := range values {
				if i == 0 {
					sqlStr += "%v"
				} else {
					sqlStr += ",%v"
				}
				vList[i] = v
			}
			sqlStr += ")"
			c.App.Println(fmt.Sprintf(sqlStr, vList...))
			rows, err := db.Query(fmt.Sprintf(sqlStr, vList...))
			if err != nil {
				c.App.Println("执行该存储过程失败，错误：", err)
				return nil
			}
			defer rows.Close()
			for rows.Next() {
				//rows.Scan()
			}
			return nil
		},
	})
	p.app.AddCommand(&grumble.Command{
		Name: "quit",
		Help: "退出数据库",
		Args: func(a *grumble.Args) {

		},
		Run: func(c *grumble.Context) error {
			if p.quit != nil {
				p.quit()
			}
			return nil
		},
	})
}

func (p *procedure) Run() {
	p.addCommands()
}
