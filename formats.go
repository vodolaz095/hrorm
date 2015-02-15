package hrorm

// Metadata is related information to entities returned from
// GetOne and Query requests. It have fields like current page, number of elements
// in collection that match the query given,  the fields readable and writable for
// single entity
type Metadata struct {
	ModelName      string            `json:"modelName"`
	CanWrite       bool              `json:"canWrite"`
	CanDelete      bool              `json:"canDelete"`
	FieldsReadable []string          `json:"fieldsReadable"`
	FieldsWritable []string          `json:"fieldsWritable"`
	Filter         map[string]string `json:"filter"`
	Page           uint32            `json:"page"`
	Sort           string            `json:"sort"`
	ItemsPerPage   uint32            `json:"itemsPerPage"`
	NumberOfPages  uint32            `json:"numberOfPages"`
	Count          uint32            `json:"count"`
}

// Response error allows to parse erroneous responses in future
type ResponseError struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

type responseParsed struct {
	Status   string          `json:"status"`
	Code     uint            `json:"code"`
	Error    []ResponseError `json:"errors"`
	Mtd      Metadata        `json:"metadata"`
	Data     interface{}     `json:"data"`
	Id       string          `json:"id"`
	Location string          `json:"location"`
}
