package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/greboid/irc-bot/v4/plugins"
	"github.com/greboid/irc-bot/v4/rpc"
	"github.com/kouhin/envflag"
)

var (
	rpcHost     = flag.String("rpc-host", "localhost", "gRPC server to connect to")
	rpcPort     = flag.Int("rpc-port", 8001, "gRPC server port")
	rpcToken    = flag.String("rpc-token", "", "gRPC authentication token")
	channel     = flag.String("channel", "", "Channel to send messages to")
	bearerToken = flag.String("bearer-token", "", "Token to require in webhook payloads")
)

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %v", err)
		return
	}

	helper, err := plugins.NewHelper(fmt.Sprintf("%s:%d", *rpcHost, uint16(*rpcPort)), *rpcToken)
	if err != nil {
		log.Fatalf("Unable to create plugin helper: %v", err)
		return
	}

	if err := helper.RegisterWebhook("distribution", webhookHandler(helper)); err != nil {
		log.Fatalf("Error registering webhook: %v", err)
		return
	}
}

func webhookHandler(helper *plugins.PluginHelper) func(request *rpc.HttpRequest) *rpc.HttpResponse {
	return func(request *rpc.HttpRequest) *rpc.HttpResponse {
		if !hasCorrectToken(request) {
			return &rpc.HttpResponse{Status: http.StatusForbidden}
		}

		e := envelope{}
		if err := json.Unmarshal(request.Body, &e); err != nil {
			log.Printf("Failed to unmarshal webhook body: %v", err)
			return &rpc.HttpResponse{Status: http.StatusBadRequest}
		}

		for i := range e.Events {
			if e.Events[i].Action == actionPush {
				if err := handlePush(helper, e.Events[i]); err != nil {
					log.Printf("Failed to handle push event: %v", err)
				}
			}
		}

		return &rpc.HttpResponse{Status: http.StatusNoContent}
	}
}

func handlePush(helper *plugins.PluginHelper, e event) error {
	if e.Target.Tag == "" {
		return nil
	}

	return helper.SendChannelMessage(
		*channel,
		fmt.Sprintf(
			"[%s] Push to tag %s by %s: %s (%d bytes)",
			e.Target.Repository,
			e.Target.Tag,
			e.Actor.Name,
			e.Target.Digest,
			e.Target.Size,
		),
	)
}

func hasCorrectToken(request *rpc.HttpRequest) bool {
	for i := range request.Header {
		if request.Header[i].Key == "Authorization" {
			return request.Header[i].Value == fmt.Sprintf("Bearer %s", *bearerToken)
		}
	}
	return false
}
