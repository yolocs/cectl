package cmd

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/yolocs/cectl/pkg/utils"
)

var (
	ceID          string
	ceSource      string
	ceType        string
	ceSubject     string
	ceDataSchema  string
	ceContentType string
	ceData        string
	ceExts        []string
	target        string
	// add retry

	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send a CloudEvent to target",
		Long:  "Send a CloudEvent to target",
		Run: func(cmd *cobra.Command, args []string) {
			if err := defaultAttrs(); err != nil {
				utils.Errorln("Missing CloudEvent attribtues: %v", err)
				return
			}

			e := event.New()
			e.SetID(ceID)
			e.SetTime(time.Now())
			e.SetSource(ceSource)
			e.SetType(ceType)
			e.SetSubject(ceSubject)
			e.SetDataSchema(ceDataSchema)
			e.SetDataContentType(ceContentType)
			e.SetData(ceContentType, ceData)
			for _, ext := range ceExts {
				p := strings.Split(ext, "=")
				e.SetExtension(strings.ToLower(p[0]), strings.ToLower(p[1]))
			}

			ctx := context.TODO()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, nil)
			if err != nil {
				utils.Errorln("Failed to construct http request: %v", err)
				return
			}
			if err := cehttp.WriteRequest(ctx, (*binding.EventMessage)(&e), req); err != nil {
				utils.Errorln("Failed to construct http request: %v", err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			// TODO: retry
			if err != nil {
				utils.Errorln("Failed to send CloudEvent: %v", err)
			} else if resp.StatusCode/100 != 2 {
				utils.Errorln("Failed to send CloudEvent: resp code = %d", resp.StatusCode)
			}

			utils.Println("Event sent.")
		},
	}
)

func init() {
	sendCmd.Flags().StringVar(&ceID, "id", "", "CE ID")
	sendCmd.Flags().StringVar(&ceSource, "source", "", "CE source")
	sendCmd.Flags().StringVar(&ceType, "type", "", "CE type")
	sendCmd.Flags().StringVar(&ceSubject, "subject", "", "CE subject")
	sendCmd.Flags().StringVar(&ceDataSchema, "dataschema", "", "CE data schema")
	sendCmd.Flags().StringVar(&ceContentType, "contenttype", "", "CE content")
	sendCmd.Flags().StringArrayVar(&ceExts, "extensions", nil, "CE extensions")
	sendCmd.Flags().StringVar(&ceData, "data", "", "CE type")
	sendCmd.Flags().StringVarP(&target, "target", "t", "", "Target to send the event")
	sendCmd.MarkFlagRequired("target")
	rootCmd.AddCommand(sendCmd)
}

func defaultAttrs() error {
	// Required fields.
	if ceID = utils.ValueFromEnv(ceID, utils.CeOutEnvID); ceID == "" {
		ceID = uuid.New().String()
	}
	if ceSource = utils.ValueFromEnv(ceSource, utils.CeOutEnvSource); ceSource == "" {
		return errors.New("Event source missing")
	}
	if ceType = utils.ValueFromEnv(ceType, utils.CeOutEnvType); ceType == "" {
		return errors.New("Event source missing")
	}

	// Optional
	ceSubject = utils.ValueFromEnv(ceSubject, utils.CeOutEnvSubject)
	ceDataSchema = utils.ValueFromEnv(ceDataSchema, utils.CeOutEnvDataSchema)
	ceContentType = utils.ValueFromEnv(ceContentType, utils.CeOutEnvContentType)
	ceData = utils.ValueFromEnv(ceData, utils.CeOutEnvData)
	ceExts = utils.ExtsFromEnv(ceExts)

	return nil
}
