package iamgo

import (
	"encoding/json"
	"fmt"

	"github.com/khulnasoft-lab/jfather"
)

type String struct {
	inner string
	r     Range
}

func (s *String) UnmarshalJSONWithMetadata(node jfather.Node) error {
	s.r.StartLine = node.Range().Start.Line
	s.r.EndLine = node.Range().End.Line
	return node.Decode(&s.inner)
}

func (d String) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.inner)
}

type Strings struct {
	inner innerStrings
	r     Range
}

func (s *Strings) UnmarshalJSONWithMetadata(node jfather.Node) error {
	s.r.StartLine = node.Range().Start.Line
	s.r.EndLine = node.Range().End.Line
	if err := node.Decode(&s.inner); err == nil {
		return nil
	}
	s.inner = []string{""}
	return node.Decode(&s.inner[0])
}

type innerStrings []string

func (d Strings) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.inner)
}

func (v *innerStrings) UnmarshalJSONWithMetadata(node jfather.Node) error {
	var raw interface{}
	if err := node.Decode(&raw); err != nil {
		return err
	}
	switch actual := raw.(type) {
	case []string:
		*v = actual
	case []interface{}:
		var output []string
		for _, raw := range actual {
			output = append(output, fmt.Sprintf("%s", raw))
		}
		*v = output
	case string:
		*v = []string{actual}
	case bool:
		*v = []string{fmt.Sprintf("%t", actual)}
	default:
		return fmt.Errorf("cannot use %T type for multi-value JSON field", actual)
	}
	return nil
}
