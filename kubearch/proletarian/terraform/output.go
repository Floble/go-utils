package terraform

import (
	"strings"
	"os"
)

type Output struct {
	ID string
	Variables map[string]string
}

func NewOutput(id string) *Output {
	output := new(Output)

	output.ID = id
	output.Variables = make(map[string]string, 0)

	return output
}

func (output *Output) AddVariable(variable string) {
	key := output.ID + "_"
	value := strings.Replace(key, "_", ".", -1)
	key += variable
	value += variable

	output.Variables[key] = value
}

func (output *Output) Export() error {
	file, err := os.OpenFile("outputs.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(output.stringBuilder()); err != nil {
		return nil
	}

	return nil
}

func (output *Output) stringBuilder() string {
	export := ""

	for key, value := range output.Variables {
		export += "output \"" + strings.Replace(key, "[0]", "", -1) + "\" {"
		export += "\n  value = \"${module." + value + "}\""
		export += "\n}\n\n"
	}

	return export
}