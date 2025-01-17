package i18n

import (
	"fmt"
	"strings"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
)

func GetMessage(id tg.MessageKind, locale string) string {
	if len(locale) > 5 {
		locale = locale[:5]
	}
	locale = strings.Replace(strings.ToLower(locale), "_", "-", 1)

	code := id.String()
	key := fmt.Sprintf("%s.%s", locale, code)

	mutex.RLock()
	trans := messages[key]
	mutex.RUnlock()
	if trans != "" {
		return trans
	}

	if m, err := getTranslations(code); err == nil {
		mutex.Lock()
		for k, v := range m {
			messages[k] = v
		}
		trans = messages[key]
		mutex.Unlock()
		return trans
	}

	// failed to find message in i18n!
	trans = "Internal error: unsupported MessageKind"
	mutex.Lock()
	messages[key] = trans
	mutex.Unlock()
	return trans
}

func ProcessEvent(data map[string]any) {
	if category := conv.StringFromAny(data[consts.FieldCategory]); category == "message" {
		locale := conv.StringFromAny(data[consts.FieldLocale])
		code := conv.StringFromAny(data[consts.FieldCode])
		trans := conv.StringFromAny(data[consts.FieldTranslation])

		if locale != "" && code != "" && trans != "" {
			key := fmt.Sprintf("%s.%s.%s", locale, category, code)
			mutex.Lock()
			messages[key] = trans
			mutex.Unlock()
		}
	}
}

var (
	mutex    sync.RWMutex
	messages = make(map[string]string, 128)
)
