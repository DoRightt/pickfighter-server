package services

import (
	"fmt"
	"net/http"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/version"
	"time"

	"github.com/spf13/viper"
)

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

	httplib.ResponseJSON(w, result)
}
