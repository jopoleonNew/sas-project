package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	vkhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/vkontakte"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"

	"github.com/gorilla/mux"
)

func GetAuthLink(w http.ResponseWriter, r *http.Request) {

	log.Println("GetAuthLink used")
	switch vars := mux.Vars(r); vars["source"] {
	case "yandex":
		ctx := context.WithValue(r.Context(), "source", "Яндекс Директ")
		context.WithValue(ctx, "YandexRole", r.FormValue("accrole"))
		//return vkhandlers.VKauthorize()
		yandexhandlers.GetYandexAuthLink(w, r.WithContext(ctx))
	case "vkontakte":
		ctx := context.WithValue(r.Context(), "source", "Вконтакте")
		vkhandlers.GetVKAuthLink(w, r.WithContext(ctx))
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
