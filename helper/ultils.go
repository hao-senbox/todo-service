package helper

import (
	"context"
	"strconv"
	"todo-service/pkg/constants"
)

func GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)

	if lang, ok := ctx.Value(constants.AppLanguageKey).(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	return headers
}
