package helper

type Select2Item struct {
	Id          string      `json:"id"`
	Value       string      `json:"value"`
	Text        string      `json:"text"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Selected    bool        `json:"selected"`
	Disabled    bool        `json:"disabled"`
	HasChildren bool        `json:"has_children"`
	Children    interface{} `json:"children"`
	ParentId    string      `json:"parent_id"`
}

type Select2Request struct {
	Q     string
	Page  int
	Size  int
	Field []string
}

type Select2Response struct {
	Results []Select2Item `json:"results"`
	Count   int           `json:"count"`
}
