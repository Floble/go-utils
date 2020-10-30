package terraform

import (
	"os"
	"strings"
	"strconv"
)

type Module struct {
	ID, Provider, Source string
	Variables map[string]interface{}
}

func NewModule(id, provider, source string) *Module {
	module := new(Module)

	module.ID = id
	module.Provider = provider
	module.Source = "./factors_of_production/infrastructure/" + strings.ToLower(provider) + "/" + strings.ToLower(source)
	module.Variables = make(map[string]interface{})

	return module
}

func (module *Module) AddVariable(key string, value interface{}) {
	module.Variables[key] = value
}

func (module *Module) Export() error {
	file, err := os.OpenFile("main.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(module.stringBuilder()); err != nil {
		return nil
	}

	return nil
}

func (module *Module) stringBuilder() string {
	export := "\n\nmodule \"" + module.ID + "\" {\n  source = \"" + module.Source + "\""

	for variable, value := range module.Variables {
		switch value.(type) {
		case string:
			if value.(string) != "" {
				export = export + "\n  " + strings.ToLower(variable) + " = \"" + value.(string) + "\""
			}
		case []string:
			elements := value.([]string)
			export = export + "\n  " + strings.ToLower(variable) + " = "
			for i, element := range elements {
				if i == 0 && element != "" {
					export = export + "[\"" + element + "\""
					if i == len(elements)-1 {
						export = export + "]"
					}
				} else if i == len(elements)-1 && element != "" {
					export = export + ",\"" + element + "\"]"
				} else if element != "" {
					export = export + ",\"" + element + "\""
				}
			}
		case int:
			export = export + "\n  " + strings.ToLower(variable) + " = " + value.(string)
		case bool:
			export = export + "\n  " + strings.ToLower(variable) + " = \"" + strconv.FormatBool(value.(bool)) + "\""
		default:
			export = export + ""
		}
	}

	export = export + "\n}"

	return export
}