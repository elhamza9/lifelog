package server

import (
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

// JSONReqActivity is used to unmarshal a json activity
type JSONReqActivity struct {
	ID       domain.ActivityID `json:"id"`
	Label    string            `json:"label"`
	Desc     string            `json:"desc"`
	Place    string            `json:"place"`
	Time     time.Time         `json:"time"`
	Duration time.Duration     `json:"duration"`
	TagIds   []domain.TagID    `json:"tagIds"`
}

// ToDomain constructs and returns a domain.Activity from a JSONReqActivity
func (reqAct JSONReqActivity) ToDomain() domain.Activity {
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range reqAct.TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	// Call adding service
	return domain.Activity{
		ID:       reqAct.ID,
		Label:    reqAct.Label,
		Desc:     reqAct.Desc,
		Place:    reqAct.Place,
		Time:     reqAct.Time,
		Duration: reqAct.Duration,
		Tags:     tags,
	}
}

// JSONRespDetailActivity is used to marshal an activity to json
type JSONRespDetailActivity struct {
	ID       domain.ActivityID     `json:"id"`
	Label    string                `json:"label"`
	Desc     string                `json:"desc"`
	Place    string                `json:"place"`
	Time     time.Time             `json:"time"`
	Duration time.Duration         `json:"duration"`
	Expenses []JSONRespListExpense `json:"expenses"`
	Tags     []domain.Tag          `json:"tags"`
}

// JSONRespListActivity is used to marshal an activity to json
type JSONRespListActivity struct {
	ID       domain.ActivityID `json:"id"`
	Label    string            `json:"label"`
	Desc     string            `json:"desc"`
	Place    string            `json:"place"`
	Time     time.Time         `json:"time"`
	Duration time.Duration     `json:"duration"`
}

// From constructs a JSONRespDetailActivity object from a domain.Activity object
func (respAct *JSONRespDetailActivity) From(act domain.Activity, expenses []domain.Expense) {
	(*respAct).ID = act.ID
	(*respAct).Label = act.Label
	(*respAct).Place = act.Place
	(*respAct).Desc = act.Desc
	(*respAct).Time = act.Time
	(*respAct).Duration = act.Duration
	respExpenses := make([]JSONRespListExpense, len(expenses))
	var respExp JSONRespListExpense
	for i, exp := range expenses {
		respExp.From(exp)
		respExpenses[i] = respExp
	}
	(*respAct).Expenses = respExpenses
	(*respAct).Tags = act.Tags
}

// From constructs a JSONRespListActivity object from a domain.Activity object
func (respAct *JSONRespListActivity) From(act domain.Activity) {
	(*respAct).ID = act.ID
	(*respAct).Label = act.Label
	(*respAct).Place = act.Place
	(*respAct).Desc = act.Desc
	(*respAct).Time = act.Time
	(*respAct).Duration = act.Duration
}
