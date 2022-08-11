package jsonmodels

//JSON models
type JSONURL struct {
	URL string `json:"url"`
}

type JSONResult struct {
	Result string `json:"result"`
}

type JSONURLTouple struct {
	Short string `json:"short_url"`
	Long string `json:"original_url"`
}
