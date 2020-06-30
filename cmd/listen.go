package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
	"github.com/yolocs/cectl/pkg/utils"
)

var (
	port    int
	command string
	server  *http.Server

	listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listen CloudEvents and trigger action",
		Long:  "Listen CloudEvents and trigger action",
		Run: func(cmd *cobra.Command, args []string) {
			parts, err := shlex.Split(command)
			if err != nil {
				utils.Errorln("Failed to parse command: %v", err)
				return
			}
			if len(parts) <= 0 {
				utils.Errorln("Command is empty")
				return
			}

			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				utils.Errorln("Failed to create TCP listener: %v", err)
				return
			}

			server = &http.Server{
				Addr: listener.Addr().String(),
				Handler: &handler{
					path: parts[0],
					args: parts[1:],
				},
			}
			setupCloseHandler()

			if err := server.Serve(listener); err != nil {
				utils.Warnln("Server closed: %v", err)
			}
		},
	}
)

func init() {
	listenCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port to listening to CloudEvents")
	listenCmd.Flags().StringVar(&command, "cmd", "", "The command to run on receiving CloudEvents")
	listenCmd.MarkFlagRequired("cmd")
	rootCmd.AddCommand(listenCmd)
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nTerminating...")
		if server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		}
		os.Exit(0)
	}()
}

type handler struct {
	path string
	args []string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e, err := toEvent(req)
	if err != nil {
		utils.Errorln("Invalid request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	c := exec.CommandContext(req.Context(), h.path, h.args...)
	c.Env = append(os.Environ(), utils.EvnsFromEvent(e)...)
	output, err := c.CombinedOutput()
	utils.PrintCmdOutput(e.ID(), output)
	if err != nil {
		utils.Errorln("Executing command returned error: %v", err)
	}
}

func toEvent(request *http.Request) (*event.Event, error) {
	message := cehttp.NewMessageFromHttpRequest(request)
	defer func() {
		if err := message.Finish(nil); err != nil {
			utils.Warnln("Failed to close message: %v", err)
		}
	}()
	// If encoding is unknown, the message is not an event.
	if message.ReadEncoding() == binding.EncodingUnknown {
		return nil, fmt.Errorf("Encoding is unknown; not a CloudEvent request? %+v", request)
	}
	event, err := binding.ToEvent(request.Context(), message)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert request to event: %w", err)
	}
	return event, nil
}
