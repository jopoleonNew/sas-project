package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	vkhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/vkontakte"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"

	"github.com/gorilla/mux"
	"gogs.itcloud.pro/SAS-project/sas/handlers/adWords"
)

func GetAuthLink(w http.ResponseWriter, r *http.Request) {

	switch vars := mux.Vars(r); vars["source"] {
	case "yandex":
		ctx := context.WithValue(r.Context(), "source", "Яндекс Директ")
		//context.WithValue(ctx, "YandexRole", r.FormValue("accrole"))
		yandexhandlers.GetYandexAuthLink(w, r.WithContext(ctx))
	case "vkontakte":
		ctx := context.WithValue(r.Context(), "source", "Вконтакте")
		vkhandlers.GetVKAuthLink(w, r.WithContext(ctx))
	case "adwords", "Adwords", "AdWords":
		//ctx := context.WithValue(r.Context(), "source", "YouTube")
		ctx := context.WithValue(r.Context(), "source", "AdWords")
		adWords.GetAdWordsAuthLink(w, r.WithContext(ctx))
		//logrus.Info("AdWords GetAuthLink")
		//fmt.Fprintf(w, "AdWords account are not availiable now: %s", vars["source"])
		return
	case "youtube":
		//ctx := context.WithValue(r.Context(), "source", "YouTube")
		fmt.Fprintf(w, "YouTube account are not availiable now: %s", vars["source"])
		return

	case "":
		log.Println("GetAuthLink Error: no source")
		return
	default:
		fmt.Fprintf(w, "Unknow account source: %s", vars["source"])
		return
	}

}
