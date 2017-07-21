package model

import (
	"github.com/nk2ge5k/goyad/campaigns"
	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

func AdaptYandexCampaings(in campaigns.GetResponse) (out []Campaign) {

	cs := make([]Campaign, len(in.Campaigns))
	for i, camp := range in.Campaigns {
		cs[i].ID = int(camp.Id)
		cs[i].Status = string(camp.Status)
		cs[i].Name = camp.Name
		cs[i].Type = string(camp.Type)
		cs[i].Owner = string(camp.ClientInfo)
	}
	return cs
}
func AdaptVKCampaings(in vkontakteAPI.AdsCampaigns, owner string) (out []Campaign) {
	cs := make([]Campaign, len(in.Response))
	for i, camp := range in.Response {
		cs[i].ID = int(camp.ID)
		cs[i].Name = camp.Name
		cs[i].Type = string(camp.Type)
		cs[i].Owner = owner
		//0 — кампания остановлена
		//1 — кампания запущена
		//2 — кампания удалена
		switch camp.Status {
		case 0:
			cs[i].Status = "STOPPED"
		case 1:
			cs[i].Status = "ACCEPTED"
		case 2:
			cs[i].Status = "DELETED"
		default:
			cs[i].Status = "UNKNOWNSTATUS"
		}
	}
	return cs
}
