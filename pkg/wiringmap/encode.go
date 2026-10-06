package wiringmap

import (
	"encoding/json"
)

// MarshalJSON encodes doc with stable formatting for committed examples.
func MarshalJSON(doc *Document, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(doc, "", "  ")
	}
	return json.Marshal(doc)
}
