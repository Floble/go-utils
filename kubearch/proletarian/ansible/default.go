package ansible

import (
	"io/ioutil"
	"strings"
	"regexp"
)

type Default struct {
	Path string
	Variables map[string]string
}

func NewDefault(path string) *Default {
	d := new(Default)
	d.Path = path
	d.Variables = make(map[string]string, 0)
	d.getVariables()

	return d
}

func (d *Default) getVariables() error {
	input, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		match, _ := regexp.MatchString(".+ = !kubearch{{ .+\\..+\\..+ }}", line)
		if match {
			content := strings.ReplaceAll(strings.Split(strings.ReplaceAll(strings.Split(line, "=")[1], " ", ""), "{{")[1], "}}", "")
			element, id, variable := d.ParseVariable(content)
			d.Variables[element + "." + id + "." + variable] = ""
		}
	}

	return nil
}

func (d *Default) ParseVariable(content string) (string, string, string) {
	element := strings.Split(content, ".")[0]
	id := strings.Split(content, ".")[1]
	variable := strings.Split(content, ".")[2]

	return element, id, variable
}

func (d *Default) SetVariableValue(variable, value string) {
	d.Variables[variable] = value
}

func (d *Default) ReplaceVariables() error {
	input, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		match, _ := regexp.MatchString(".+ = !kubearch{{ .+\\..+\\..+ }}", line)
		if match {
			ansibleVariable := strings.ReplaceAll(strings.Split(line, "=")[0], " ", "")
			content := strings.ReplaceAll(strings.Split(strings.ReplaceAll(strings.Split(line, "=")[1], " ", ""), "{{")[1], "}}", "")
			element, id, variable := d.ParseVariable(content)

			lines[i] = "#" + line + "\n" + strings.ToLower(ansibleVariable) + ": " + d.Variables[element + "." + id + "." + variable]
		}
	}

	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(d.Path, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}