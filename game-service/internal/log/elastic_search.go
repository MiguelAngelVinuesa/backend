package log

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
)

const monitorKey = "gs"

func ReportError(message string) {
	if config.MonitorHost == "" {
		return
	}

	data, _ := json.Marshal(map[string]any{consts.FieldType: monitorKey, consts.FieldMessage: message})

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/report-error", config.MonitorHost), bytes.NewBuffer(data))
	req.Header.Set(consts.ContentType, consts.ApplicationJSON)
	req.Header.Set(consts.XApiKey, config.MonitorKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger.Error(consts.MsgMonitorFailed, consts.FieldError, err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logger.Error(consts.MsgMonitorFailed, consts.FieldError, consts.MsgInvalidStatusCode, consts.FieldStatus, resp.StatusCode)
	}
}
