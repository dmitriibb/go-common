package model

type Recipe struct {
	Name          string    `json:"name" bson:"name"`
	Description   string    `json:"description" bson:"description"`
	Ingredients   []string  `json:"ingredients" bson:"ingredients"`
	TimeToWaitSec int64     `json:"timeToWaitSec" bson:"timeToWaitSec"`
	SubStages     []*Recipe `json:"subStages" bson:"subStages"`
}
