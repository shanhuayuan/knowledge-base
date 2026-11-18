package knowledge

type QueryOption struct {
	KbId  string `json:"kbId" uri:"kbId"`
	Query string `json:"query"`
}
