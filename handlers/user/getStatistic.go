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

func GetStatistic(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStatistic used")
	switch vars := mux.Vars(r); vars["source"] {
	case "yandex", "Яндекс Директ":
		ctx := context.WithValue(r.Context(), "source", "Яндекс Директ")
		context.WithValue(ctx, "YandexRole", r.FormValue("accrole"))
		yandexhandlers.GetStatSliceHandler(w, r.WithContext(ctx))
	case "vkontakte":
		ctx := context.WithValue(r.Context(), "source", "Вконтакте")
		vkhandlers.CollectVKStatistic(w, r.WithContext(ctx))
	case "youtube":
		//ctx := context.WithValue(r.Context(), "source", "YouTube")
		fmt.Fprintf(w, "YouTube accounts are not availiable now: %s", vars["source"])
		return
	case "":
		log.Println("AddAccount Error: no source")
		return
	default:
		fmt.Fprintf(w, "Unknow account source: %s", vars["source"])
		return
	}
}

// GetAccountStat redirecting to handler for getting single account statistic
func GetAccountStat(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAccountStat used")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	switch vars := mux.Vars(r); vars["source"] {
	case "yandex", "Яндекс Директ", "Yandex":
		ctx := context.WithValue(r.Context(), "source", "Яндекс Директ")
		context.WithValue(ctx, "YandexRole", r.FormValue("accrole"))
		yandexhandlers.CollectYandexStatistic(w, r.WithContext(ctx))
	case "vkontakte", "Вконтакте", "Vkontakte":
		ctx := context.WithValue(r.Context(), "source", "Вконтакте")
		vkhandlers.CollectVKStatistic(w, r.WithContext(ctx))
	case "adwords", "Adwords", "AdWords":
		fmt.Fprintf(w, "AdWords accounts are not availiable now: %s", vars["source"])
		return
	case "youtube":
		//ctx := context.WithValue(r.Context(), "source", "YouTube")
		fmt.Fprintf(w, "YouTube accounts are not availiable now: %s", vars["source"])
		return

	case "":
		logrus.Error("AddAccount Error: no source")
		return
	default:
		fmt.Fprintf(w, "Unknow account source: %s", vars["source"])
		return
	}
}
