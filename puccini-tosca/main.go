package main

import (
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/puccini-tosca/commands"

	_ "github.com/tliron/commonlog/simple"
)

func main() {
	util.ExitOnSIGTERM()
	commands.Execute()
	util.Exit(0)
}
