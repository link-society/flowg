package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type IndexMap map[string][]string

var _ pflag.Value = (*IndexMap)(nil)

func (im *IndexMap) String() string {
	return ""
}

func (im *IndexMap) Type() string {
	return "key=value"
}

func (im *IndexMap) Set(s string) error {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format for index map: %s, expected key=value", s)
	}

	key, val := parts[0], parts[1]
	(*im)[key] = append((*im)[key], val)

	return nil
}
