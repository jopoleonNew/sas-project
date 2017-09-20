package model

import (
	"strconv"

	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI/campaigns"
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
			cs[i].Status = "UNKNOWN_STATUS"
		}
	}
	return cs
}

func AdaptVKAds(in vkontakteAPI.Ads) []Ad {
	ads := make([]Ad, len(in.Response))
	for i, a := range in.Response {
		ads[i].ID, _ = strconv.Atoi(a.ID)
		ads[i].CampID = a.CampaignID
		ads[i].Name = a.Name
		if a.CostType == 0 {
			//цена за переход в копейках.
			ads[i].CPC, _ = strconv.Atoi(a.Cpc)
		}
		if a.CostType == 1 {
			//цена за 1000 показов в копейках.
			ads[i].CPM, _ = strconv.Atoi(a.Cpm)
		}
		//0 — кампания остановлена
		//1 — кампания запущена
		//2 — кампания удалена
		switch a.Status {
		case 0:
			ads[i].Status = "STOPPED"
		case 1:
			ads[i].Status = "ACCEPTED"
		case 2:
			ads[i].Status = "DELETED"
		default:
			ads[i].Status = "UNKNOWN_STATUS"
		}
	}
	return ads
}
