package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
	vkhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/vkontakte"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	log.Println("GetToken used")
	switch vars := mux.Vars(r); vars["source"] {
	case "yandex":
		ctx := context.WithValue(r.Context(), "source", "Яндекс Директ")
		context.WithValue(ctx, "YandexRole", r.FormValue("accrole"))
		yandexhandlers.AddYandexAccount(w, r.WithContext(ctx))
	case "vkontakte":
		ctx := context.WithValue(r.Context(), "source", "Вконтакте")
		vkhandlers.AddVKAccount(w, r.WithContext(ctx))
	case "adwords":
		logrus.Warn("Something from AdWords: \n,", r)
	case "youtube":
		//ctx := context.WithValue(r.Context(), "source", "YouTube")
		fmt.Fprintf(w, "YouTube account are not availiable now: %s", vars["source"])
		return

	case "":
		log.Println("AddAccount Error: no source")
		return
	default:
		fmt.Fprintf(w, "Unknow account source: %s", vars["source"])
		return
	}
}
