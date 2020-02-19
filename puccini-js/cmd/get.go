package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&output, "output", "o", "", "output to file (default is stdout)")
}

var getCmd = &cobra.Command{
	Use:   "get [NAME] [[Clout PATH or URL]]",
	Short: "Get JavaScript scriptlet from Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]

		var path string
		if len(args) == 2 {
			path = args[1]
		}

		clout, err := ReadClout(path)
		common.FailOnError(err)

		scriptlet, err := js.GetScriptlet(scriptletName, clout)
		common.FailOnError(err)

		if !terminal.Quiet {
			err = formatpkg.WriteOrPrint(scriptlet, format, terminal.Stdout, pretty, output)
			common.FailOnError(err)
		}
	},
}
