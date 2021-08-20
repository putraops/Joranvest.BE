package helper

type Select2Item struct {
	Id          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	Selected    bool
	Disabled    bool
}

type Select2Request struct {
	Q     string
	Page  int
	Field []string
}

type Select2Response struct {
	Results []Select2Item `json:"results"`
	Count   int           `json:"count"`
}
