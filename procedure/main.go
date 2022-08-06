package procedure

import (
	"database/sql"
	"fmt"
	"procedure-run/connector"
	"strings"

	"github.com/cheekybits/genny/generic"
	"github.com/desertbit/grumble"
	"github.com/liushuochen/gotable"
)

type procedure struct {
	app           *grumble.App
	databaseAlias string
	collection    *connector.Collection
	quit          func() bool
}

type Item generic.Type

func NewProcedure(app *grumble.App, databaseAlias string, collection *connector.Collection, quit func() bool) *procedure {
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
			var sqlStr = "call `%v`("
			var vList []interface{} = make([]interface{}, len(values))
			vplus := 0

			for i, v := range values {
				if i > 0 {
					sqlStr += ","
				}
				if strings.HasPrefix(v, "@") {
					sqlStr += v
					vplus += 1
				} else {
					sqlStr += "?"
					vList[i-vplus] = v
				}

			}
			sqlStr += ")"
			vList = vList[0 : len(values)-vplus]
			connector := connector.Database(p.collection)
			rows, err := connector.Query(fmt.Sprintf(sqlStr, name), vList...)
			if err != nil {
				c.App.Config().ErrorColor.Println("执行该存储过程失败，错误：", err)
				return nil
			}
			defer rows.Close()
			tablePrint(rows)
			return nil
		},
	})
	// p.app.AddCommand(&grumble.Command{
	// 	Name: "fetch",
	// 	Help: "获取变量",
	// 	Args: func(a *grumble.Args) {
	// 		a.StringList("names", "变量名称")
	// 	},
	// 	Run: func(c *grumble.Context) error {
	// 		names := c.Args.StringList("names")
	// 		db, err := sql.Open("mysql", p.collection)
	// 		if err != nil {
	// 			c.App.Config().ErrorColor.Println("该数据库连接失败，错误：", err)
	// 			return nil
	// 		}
	// 		err = db.Ping()
	// 		if err != nil {
	// 			c.App.Config().ErrorColor.Println("该数据库连接失败，错误：", err)
	// 			return nil
	// 		}

	// 		defer db.Close()

	// 		var sqlStr = "select "
	// 		for i, name := range names {
	// 			if i > 0 {
	// 				sqlStr += ","
	// 			}
	// 			sqlStr += "@" + name + " as " + name
	// 		}
	// 		fmt.Println(sqlStr)
	// 		rows, _ := db.Query(sqlStr)
	// 		defer rows.Close()
	// 		tablePrint(rows)
	// 		return nil
	// 	},
	// })
	p.app.AddCommand(&grumble.Command{
		Name: "ls",
		Help: "列出存储过程",
		Args: func(a *grumble.Args) {

		},
		Run: func(c *grumble.Context) error {
			connector := connector.Database(p.collection)
			// show procedure status
			rows, err := connector.Query("select routine_schema,routine_name,definer,last_altered,created from information_schema.ROUTINES")
			if err != nil {
				c.App.Config().ErrorColor.Println("获取该存储过程列表失败，错误：", err)
				return nil
			}
			defer rows.Close()
			table, err := gotable.Create("所属数据库", "存储过程名称", "创建时间")
			if err != nil {
				return nil
			}
			for rows.Next() {
				dbName, procedureName, procedureDefiner, modified, created := "", "", "", "", ""
				if err = rows.Scan(&dbName, &procedureName, &procedureDefiner, &modified, &created); err == nil {
					table.AddRow([]string{dbName, procedureName, created})
				}
			}
			fmt.Println(table)
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

func tablePrint(rows *sql.Rows) (err error) {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	table, err := gotable.Create(columns...)
	if err != nil {
		return err
	}

	vals := make([]interface{}, len(columns)) //建立接口
	valsp := make([]interface{}, len(vals))   //建立接口指针的接口
	//将接口转换为指针类型的接口
	for i := range vals {
		valsp[i] = &vals[i]
	}
	for rows.Next() {
		valArr := make([]string, len(vals))
		if err = rows.Scan(valsp...); err == nil {
			for i, v := range vals { //注意：此处用vals
				valArr[i] = fmt.Sprint(v)
				if v, ok := v.([]byte); ok {
					valArr[i] = string(v)
				}
			}
		}
		table.AddRow(valArr)
	}
	fmt.Println(table)
	return nil
}

func (p *procedure) Run() {
	p.addCommands()
}
