package main

import (
	"github.com/dormunis/gitd/cli"
)

func main() {
	// github.com/dormunis/gitd review purge
	// TODO: get all older than 1 week tasks, go over them iteractively, update task manager and archive manager
	// TODO: add option to control the age of the tasks to be purged
	// TODO: add archive manager (obsidian note)
	// TODO: add customizable setting for append/new archive note for archived tasks
	// TODO: add customizable setting whether or not to document archived tasks
	// TODO: add customziable option to format archived tasks (e.g. markdown, etc.)
	// TODO: add customizable option to format the filename of the archived document (e.g. date, etc.)

	// github.com/dormunis/gitd next
	// TODO: add get next action by filters (sorted by priority), e.g.: github.com/dormunis/gitd next --low

	cli.Execute()
}
