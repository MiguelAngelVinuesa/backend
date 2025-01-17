package elasticsearch

// Field contains index field metadata.
type Field struct {
	Name     string
	DataType string
	Index    bool
}

// Index contains index metadata.
type Index struct {
	Name   string
	Fields []Field
}

type putIndexField struct {
	Type  string `json:"type"`
	Index bool   `json:"index"`
}

type putIndexReq struct {
	Mappings struct {
		Properties map[string]putIndexField `json:"properties"`
	} `json:"mappings"`
}

type postDeleteResp struct {
	Deleted int `json:"deleted"`
}
