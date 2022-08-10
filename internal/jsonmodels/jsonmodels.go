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

type JSONBulkInput struct {
	CrlID string `json:"correlation_id"`
	Long string `json:"original_url,omitempty"`
	Short string `json:"short_url,omitempty"`
	URL string `json:"-"`
}

type JSONBulkOutput struct {
	CrlID string `json:"correlation_id"`
	Short string `json:"short_url"`
	Long string `json:"-"`
}
