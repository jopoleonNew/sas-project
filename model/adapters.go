package model

import (
	"github.com/nk2ge5k/goyad/campaigns"
)

func AdaptYandexCampaings(in campaigns.GetResponse) (out []Campaign) {

	cs := make([]Campaign, len(in.Campaigns))
	for i, camp := range in.Campaigns {
		cs[i].ID = int(camp.Id)
		cs[i].Status = string(camp.Status)
		cs[i].Name = camp.Name
		cs[i].Type = string(camp.Type)
	}
	return cs
}
func AdaptVKCampaings() {}
