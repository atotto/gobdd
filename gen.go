package gobdd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"code.google.com/p/go.tools/imports"

	"github.com/atotto/gobdd/internal"
	"github.com/muhqu/go-gherkin"
	"github.com/muhqu/go-gherkin/nodes"
)

func featurecomment(feature nodes.FeatureNode) string {
	return "// " + strings.Replace(feature.Description(), "\n", "\n// ", -1)
}

func scenarioname(scenario nodes.ScenarioNode) string {
	return strings.Replace(strings.Title(scenario.Title()), " ", "", -1)
}

func scenariocomment(scenario nodes.ScenarioNode) string {
	node := scenario.Comment()
	if node != nil {
		return node.Comment()
	} else {
		return scenarioname(scenario)
	}
}

func scenarioToTestCode(scenario nodes.ScenarioNode) string {
	var buf bytes.Buffer

	switch scenario.(type) {
	case nodes.OutlineNode:
		outline := scenario.(nodes.OutlineNode)
		if len(outline.Examples().Table().Rows()) < 2 {
			panic("Scenario Outline: not enought examples.")
		}
		header := outline.Examples().Table().Rows()[0]
		examples := outline.Examples().Table().Rows()[1:]
		buf.WriteString("tests := []struct{\n")
		for i, name := range header {
			buf.WriteString(strings.ToLower(name))
			buf.WriteString(" ")
			kind := internal.InspectKind(examples[0][i])
			buf.WriteString(kind.String())
			buf.WriteString("\n")
		}
		buf.WriteString("}{\n")
		for i, example := range examples {
			buf.WriteString("{")
			for j, name := range example {
				kind := internal.InspectKind(name)
				if kind == reflect.String {
					examples[i][j] = "\"" + name + "\""
				}
			}
			buf.WriteString(strings.Join(example, ", "))
			buf.WriteString("},\n")
		}
		buf.WriteString("}\n\n")

		buf.WriteString("for _, tt := range tests {")

		re := regexp.MustCompile("<(.+?)>")
		for _, step := range scenario.Steps() {
			buf.WriteString("\n")
			out := step.Text()
			mach := re.FindAllStringSubmatch(out, -1)
			for i, val := range mach {
				out = strings.Replace(out, val[0], fmt.Sprintf("_val%d_", i+1), 1)
			}
			out = step.StepType() + "_" + strings.Replace(strings.Title(out), " ", "", -1)
			buf.WriteString(out)
			buf.WriteString("(t")
			for i, val := range mach {
				if i != len(mach) {
					buf.WriteRune(',')
				}
				buf.WriteString("tt.")
				name := val[1]
				buf.WriteString(strings.ToLower(name))
			}
			buf.WriteRune(')')
		}
		buf.WriteString("}")
	case nodes.ScenarioNode:
		re := regexp.MustCompile("\"(.+?)\"")
		for _, step := range scenario.Steps() {
			buf.WriteString("\n")
			out := step.Text()
			mach := re.FindAllStringSubmatch(out, -1)
			for i, val := range mach {
				out = strings.Replace(out, val[0], fmt.Sprintf("_val%d_", i+1), 1)
			}
			out = step.StepType() + "_" + strings.Replace(strings.Title(out), " ", "", -1)
			buf.WriteString(out)
			buf.WriteString("(t")
			for i, val := range mach {
				if i != len(mach) {
					buf.WriteRune(',')
				}
				name := val[1]
				kind := internal.InspectKind(name)
				if kind == reflect.String {
					name = "\"" + name + "\""
				}
				buf.WriteString(name)
			}
			buf.WriteRune(')')
		}
	case nodes.BackgroundNode:
	default:
		panic("not supported")
	}
	return buf.String()
}

var tp *template.Template

func init() {
	funcMap := template.FuncMap{
		"featurecomment":  featurecomment,
		"scenarioname":    scenarioname,
		"scenariocomment": scenariocomment,
		"scenario":        scenarioToTestCode,
	}

	tp = template.Must(template.New("feature_test").Funcs(funcMap).Parse(testTemplate))
}

func Gen(featurefile string) ([]byte, error) {
	buf, err := ioutil.ReadFile(featurefile)
	if err != nil {
		return nil, err
	}

	feature, err := gherkin.ParseGherkinFeature(string(buf))
	if err != nil {
		return nil, err
	}

	var code bytes.Buffer
	err = tp.Execute(&code, feature)
	if err != nil {
		return nil, err
	}

	out, err := imports.Process("tmp_test.go", code.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

const testTemplate = `
package feature

import "testing"

{{featurecomment .}} 

{{range $scenario := .Scenarios}}
// {{scenariocomment $scenario}}
func Test_{{scenarioname $scenario}}(t *testing.T) {
     {{scenario $scenario}}
}
{{end}}
`
