package main

import (
	"procedure-run/database"

	grumble "github.com/desertbit/grumble"
	color "github.com/fatih/color"
)

var app = grumble.New(&grumble.Config{
	Name:                  "存储过程应用",
	Description:           "连接数据库，并执行存储过程程序",
	Prompt:                "存储过程应用 » ",
	PromptColor:           color.New(color.FgYellow, color.Bold),
	HelpHeadlineColor:     color.New(color.FgYellow),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,
	//HelpHeadlineColor:     color.New(color.FgHiBlue),
	//HelpHeadlineUnderline: true,
	//HelpSubCommands: true,
	HistoryFile: ".procedure-run-history",
	Flags: func(f *grumble.Flags) {

	},
})

func init() {
	app.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println("                   _   _     ")
		a.Println(" ___ ___ _ _ _____| |_| |___ ")
		a.Println("| . |  _| | |     | . | | -_|")
		a.Println("|_  |_| |___|_|_|_|___|_|___|")
		a.Println("|___|                        ")
		a.Println()
	})
}

func main() {
	app.SetDefaultPrompt()

	database.New(app).Run()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
