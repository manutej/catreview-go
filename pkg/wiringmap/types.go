package wiringmap

// Document is a wiringmap.v0 instance.
type Document struct {
	SchemaVersion string     `json:"schema_version"`
	Meta          Meta       `json:"meta"`
	Junctions     []Junction `json:"junctions,omitempty"`
	Units         []Unit     `json:"units"`
	Edges         []Edge     `json:"edges"`
}

type Meta struct {
	Repo  string `json:"repo"`
	Ref   string `json:"ref"`
	Slice string `json:"slice"`
	Notes string `json:"notes,omitempty"`
}

type Unit struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Role  string `json:"role"`
	Ports []Port `json:"ports"`
}

type Port struct {
	ID       string `json:"id"`
	Symbol   string `json:"symbol"`
	Kind     string `json:"kind"`
	Exported bool   `json:"exported"`
}

type Junction struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	JunctionKind  string `json:"junction_kind"`
}

type Edge struct {
	ID       string `json:"id"`
	From     string `json:"from"`
	To       string `json:"to"`
	EdgeKind string `json:"edge_kind"`
}
