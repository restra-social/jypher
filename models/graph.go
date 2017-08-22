package models

// Graph Struct
type Graph struct {
	Nodes Node  `json:"nodes,omitempty"`
	Edges Edges `json:"edges,omitempty"`
}

// Node Struct
type Node struct {
	ID         string                   `json:"id,omitempty"`
	Lebel      string                   `json:"lebel,omitempty"`
	Properties []map[string]interface{} `json:"properties,omitempty"`
}

// Edges Struct
type Edges struct {
	Source     string                   `json:"source,omitempty"`
	Target     string                   `json:"target,omitempty"`
	Relation   string                   `json:"relation,omitempty"`
	Properties []map[string]interface{} `json:"properties,omitempty"`
}

// JSONInfo struct contains the unmarshal JSON and the set of rules
// that might apply during Graph Model decoding
// along with the resource name and id
type JSONInfo struct {
	DecodedJSON map[string]interface{}
	Rules       map[string]interface{}
	ID          string
	Master      string
}
