package server

import "github.com/elhamza90/lifelog/internal/domain"

// JSONReqTag is used to unmarshal a json tag.
type JSONReqTag struct {
	ID   domain.TagID `json:"id"`
	Name string       `json:"name"`
}

// ToDomain constructs and returns a domain.Tag from a JSONReqTag.
func (reqTag JSONReqTag) ToDomain() domain.Tag {
	return domain.Tag{
		ID:   reqTag.ID,
		Name: reqTag.Name,
	}
}

// JSONRespDetailTag is used to marshal a tag to json.
type JSONRespDetailTag struct {
	ID   domain.TagID `json:"id"`
	Name string       `json:"name"`
}

// JSONRespListTag is used to marshal a tag in a list to json.
type JSONRespListTag struct {
	ID   domain.TagID `json:"id"`
	Name string       `json:"name"`
}

// From constructs a JSONRespDetailTag object from a domain.Tag object.
func (respExp *JSONRespDetailTag) From(tag domain.Tag) {
	(*respExp).ID = tag.ID
	(*respExp).Name = tag.Name
}

// From constructs a JSONRespListTag object from a domain.Tag object.
func (respExp *JSONRespListTag) From(tag domain.Tag) {
	(*respExp).ID = tag.ID
	(*respExp).Name = tag.Name
}
