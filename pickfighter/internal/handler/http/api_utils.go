package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"pickfighter.com/pickfighter/pkg/version"
	"pickfighter.com/pkg/httplib"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	result := httplib.SuccessfulResultMap()
	delete(result, "success")

	result["app_dev_version"] = version.DevVersion
	result["app_name"] = viper.GetString("app.name")
	result["app_run_date"] = version.RunDate
	result["app_time_alive"] = time.Now().Unix() - version.RunDate
	result["healthy"] = true
	result["message"] = fmt.Sprintf("[%s]: I'm fine, thanks!", viper.GetString("app.name"))

	authStatus := h.ctrl.GetAuthServiceHealthStatus()
	result["auth-service"] = authStatus

	eventStatus := h.ctrl.GetEventServiceHealthStatus()
	result["event-service"] = eventStatus

	fightersStatus := h.ctrl.GetFightersServiceHealthStatus()
	result["fighters-service"] = fightersStatus

	httplib.ResponseJSON(w, result)
}
