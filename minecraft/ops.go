package minecraft

type Operator struct {
	UUID                string `json:"uuid"`
	Name                string `json:"name"`
	Level               uint8  `json:"level"`
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

type ServerOperators []Operator
