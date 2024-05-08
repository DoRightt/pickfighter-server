package services

import (
	"fmt"
	"net/http"
	"fightbettr.com/fb-server/pkg/httplib"
	"fightbettr.com/fb-server/pkg/version"
	"time"

	"github.com/spf13/viper"
)

// HealthCheck handles HTTP requests for health checks.
// It returns a JSON response containing information about the application's health,
// version, uptime, and the status of registered modules.
func (h *ApiHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	result := httplib.SuccessfulResultMap()
	delete(result, "success")
	result["app_name"] = viper.GetString("app.name")
	result["app_dev_version"] = version.DevVersion
	result["app_git_version"] = version.GitVersion
	result["app_build_commit"] = version.BuildCommit
	result["app_build_date"] = version.BuildDate
	result["app_run_date"] = version.RunDate
	result["app_time_alive"] = time.Now().Unix() - version.RunDate
	result["healthy"] = true
	result["message"] = fmt.Sprintf("[%s] Bee-beep-bop... I'm working fine!", h.ServiceName)

	var names []string
	for name := range h.Services {
		names = append(names, name)
	}

	result["modules"] = names
	httplib.ResponseJSON(w, result)
}
