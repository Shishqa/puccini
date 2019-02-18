package v2018_08_31

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.v2018_08_31")

var Readers = make(map[string]tosca.Reader)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Readers["Condition"] = ReadCondition
	Readers["Output"] = ReadOutput
	Readers["Parameter"] = ReadParameter
	Readers["ParameterGroup"] = ReadParameterGroup
	Readers["Resource"] = ReadResource
	Readers["Template"] = ReadTemplate
	Readers["Value"] = ReadValue

	for name, sourceCode := range FunctionSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	// for name, sourceCode := range ConstraintClauseSourceCode {
	// 	nativeArgumentIndexes, _ := ConstraintClauseNativeArgumentIndexes[name]
	// 	DefaultScriptNamespace[name] = &tosca.Script{
	// 		SourceCode:            js.Cleanup(sourceCode),
	// 		NativeArgumentIndexes: nativeArgumentIndexes,
	// 	}
	// }
}