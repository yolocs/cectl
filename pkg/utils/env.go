package utils

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
)

const (
	CeExtPrefix = "CE_IN_EXT_"

	CeEnvID          = "CE_IN_ID"
	CeEnvSource      = "CE_IN_SOURCE"
	CeEnvType        = "CE_IN_TYPE"
	CeEnvSubject     = "CE_IN_SUBJECT"
	CeEnvTime        = "CE_IN_TIME"
	CeEnvDataSchema  = "CE_IN_DATASCHEMA"
	CeEnvContentType = "CE_IN_CONTENTTYPE"
	CeEnvData        = "CE_IN_DATA"
)

func EvnsFromEvent(e *event.Event) []string {
	envs := []string{}
	envs = append(envs, env(CeEnvID, e.ID()))
	envs = append(envs, env(CeEnvSource, e.Source()))
	envs = append(envs, env(CeEnvType, e.Type()))
	envs = append(envs, env(CeEnvSubject, e.Subject()))
	envs = append(envs, env(CeEnvTime, e.Time().Unix()))
	envs = append(envs, env(CeEnvDataSchema, e.DataSchema()))
	envs = append(envs, env(CeEnvContentType, e.DataContentType()))
	envs = append(envs, env(CeEnvData, base64.StdEncoding.EncodeToString(e.Data())))
	for k, v := range e.Extensions() {
		// Only support strings.
		if str, ok := v.(string); ok {
			envs = append(envs, env(CeExtPrefix+strings.ToUpper(k), str))
		}
	}
	return envs
}

func env(key string, val interface{}) string {
	return fmt.Sprintf("%s=%v", key, val)
}
