package commons

type JStreeResponse struct {
	Id          string           `json:"id"`
	Key         string           `json:"key"`
	Text        string           `json:"text"`
	Title       string           `json:"title"`
	Subtitle    string           `json:"subtitle"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Icon        string           `json:"icon"`
	Tipe        int              `json:"tipe"`
	JStreeState JStreeState      `json:"state"`
	Children    []JStreeResponse `json:"children"`
}

type JStreeState struct {
	Opened     bool `json:"opened" form:"opened"`
	Disabled   bool `json:"disabled" form:"disabled"`
	Selected   bool `json:"selected" form:"selected"`
	Checked    bool `json:"checked" form:"checked"`
	HasSibling bool `json:"has_sibling" form:"has_sibling"`
}
