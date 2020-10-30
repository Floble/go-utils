package terraform

import (
	"os"
	"strings"
)

type Provider struct {
	Name string
	Variables map[string]interface{}
}

func NewProvider(name string) *Provider {
	provider := new(Provider)
	provider.Name = name
	provider.Variables = make(map[string]interface{})

	return provider
}

func (provider *Provider) AddVariable(key string, value interface{}) {
	provider.Variables[key] = value
}

func (provider *Provider) Export() error {
	file, err := os.OpenFile("main.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(provider.stringBuilder()); err != nil {
		return nil
	}

	return nil
}

func (provider *Provider) stringBuilder() string {
	export := "provider \"" + provider.Name + "\" {"

	for variable, value := range provider.Variables {
		switch value.(type) {
		case string:
			export = export + "\n  " + strings.ToLower(variable) + " = \"" + strings.ToLower(value.(string)) + "\""
		case int:
			export = export + "\n  " + strings.ToLower(variable) + " = " + strings.ToLower(value.(string))
		default:
			export = export + ""
		}
	}

	export = export + "\n}"

	return export
}