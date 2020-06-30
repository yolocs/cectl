package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
)

const (
	CeInEnvExtPrefix   = "CE_IN_EXT_"
	CeInEnvID          = "CE_IN_ID"
	CeInEnvSource      = "CE_IN_SOURCE"
	CeInEnvType        = "CE_IN_TYPE"
	CeInEnvSubject     = "CE_IN_SUBJECT"
	CeInEnvTime        = "CE_IN_TIME"
	CeInEnvDataSchema  = "CE_IN_DATASCHEMA"
	CeInEnvContentType = "CE_IN_CONTENTTYPE"
	CeInEnvData        = "CE_IN_DATA"

	CeOutEnvExt         = "CE_OUT_EXTS"
	CeOutEnvID          = "CE_OUT_ID"
	CeOutEnvSource      = "CE_OUT_SOURCE"
	CeOutEnvType        = "CE_OUT_TYPE"
	CeOutEnvSubject     = "CE_OUT_SUBJECT"
	CeOutEnvTime        = "CE_OUT_TIME"
	CeOutEnvDataSchema  = "CE_OUT_DATASCHEMA"
	CeOutEnvContentType = "CE_OUT_CONTENTTYPE"
	CeOutEnvData        = "CE_OUT_DATA"
)

func EvnsFromEvent(e *event.Event) []string {
	envs := []string{}
	envs = append(envs, envStr(CeInEnvID, e.ID()))
	envs = append(envs, envStr(CeInEnvSource, e.Source()))
	envs = append(envs, envStr(CeInEnvType, e.Type()))
	envs = append(envs, envStr(CeInEnvSubject, e.Subject()))
	envs = append(envs, envStr(CeInEnvTime, e.Time().Unix()))
	envs = append(envs, envStr(CeInEnvDataSchema, e.DataSchema()))
	envs = append(envs, envStr(CeInEnvContentType, e.DataContentType()))
	// Data will be base64 string?
	envs = append(envs, envStr(CeInEnvData, string(e.Data())))
	for k, v := range e.Extensions() {
		// Only support strings.
		if str, ok := v.(string); ok {
			envs = append(envs, envStr(CeInEnvExtPrefix+strings.ToUpper(k), str))
		}
	}
	return envs
}

func envStr(key string, val interface{}) string {
	return fmt.Sprintf("%s=%v", key, val)
}

func ValueFromEnv(preemptive, key string) string {
	if preemptive != "" {
		return preemptive
	}
	return os.Getenv(key)
}

func ExtsFromEnv(preemptive []string) []string {
	if len(preemptive) > 0 {
		return preemptive
	}
	exts := os.Getenv(CeOutEnvExt)
	return strings.Split(exts, ",")
}
