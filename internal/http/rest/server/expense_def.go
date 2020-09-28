package server

import (
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

// JSONReqExpense is used to unmarshal a json expense
type JSONReqExpense struct {
	ID         domain.ExpenseID  `json:"id"`
	Label      string            `json:"label"`
	Time       time.Time         `json:"time"`
	Value      float32           `json:"value"`
	Unit       string            `json:"unit"`
	ActivityID domain.ActivityID `json:"activityId"`
	TagIds     []domain.TagID    `json:"tagIds"`
}

// ToDomain constructs and returns a domain.Expense from a JSONReqExpense
func (reqExp JSONReqExpense) ToDomain() domain.Expense {
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range reqExp.TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	return domain.Expense{
		ID:         reqExp.ID,
		Label:      reqExp.Label,
		Time:       reqExp.Time,
		Value:      reqExp.Value,
		Unit:       reqExp.Unit,
		ActivityID: reqExp.ActivityID,
		Tags:       tags,
	}
}

// JSONRespDetailExpense is used to marshal an expense to json
type JSONRespDetailExpense struct {
	ID         domain.ExpenseID  `json:"id"`
	Label      string            `json:"label"`
	Time       time.Time         `json:"time"`
	Value      float32           `json:"value"`
	Unit       string            `json:"unit"`
	ActivityID domain.ActivityID `json:"activityId"`
	Tags       []domain.Tag      `json:"tags"`
}

// JSONRespListExpense is used to marshal an expense to json
type JSONRespListExpense struct {
	ID         domain.ExpenseID  `json:"id"`
	Label      string            `json:"label"`
	Time       time.Time         `json:"time"`
	Value      float32           `json:"value"`
	Unit       string            `json:"unit"`
	ActivityID domain.ActivityID `json:"activityId"`
}

// From constructs a JSONRespDetailExpense object from a domain.Expense object
func (respExp *JSONRespDetailExpense) From(exp domain.Expense) {
	(*respExp).ID = exp.ID
	(*respExp).Label = exp.Label
	(*respExp).Time = exp.Time
	(*respExp).Value = exp.Value
	(*respExp).Unit = exp.Unit
	(*respExp).ActivityID = exp.ActivityID
	(*respExp).Tags = exp.Tags
}

// From constructs a JSONRespListExpense object from a domain.Expense object
func (respExp *JSONRespListExpense) From(exp domain.Expense) {
	(*respExp).ID = exp.ID
	(*respExp).Label = exp.Label
	(*respExp).Time = exp.Time
	(*respExp).Value = exp.Value
	(*respExp).Unit = exp.Unit
	(*respExp).ActivityID = exp.ActivityID
}
