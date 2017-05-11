package model

type YadCampaign struct {
	ID     int    `bson:"id"`
	Name   string `bson:"name"`
	Status string `bson:"status"`
}
