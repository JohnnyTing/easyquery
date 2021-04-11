package json_util

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
)

type NamingStrategyExtension struct {
	jsoniter.DummyExtension
	Translate func(string) string
}

func (extension *NamingStrategyExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		if unicode.IsLower(rune(binding.Field.Name()[0])) || binding.Field.Name()[0] == '_' {
			continue
		}
		tag, hastag := binding.Field.Tag().Lookup("json")
		if hastag {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] == "-" {
				continue // hidden field
			}
			if tagParts[0] != "" {
				continue // field explicitly named
			}
		}
		binding.ToNames = []string{extension.Translate(binding.Field.Name())}
		binding.FromNames = []string{extension.Translate(binding.Field.Name())}
	}
}

// 小写驼峰命名
func InitialLower(name string) string {
	return strcase.ToLowerCamel(name)
}
