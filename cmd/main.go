package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	argocd "github.com/argoproj-labs/argocd-ephemeral-access/api/argoproj/v1alpha1"
	api "github.com/argoproj-labs/argocd-ephemeral-access/api/ephemeral-access/v1alpha1"
	"github.com/hashicorp/go-hclog"

	"github.com/argoproj-labs/argocd-ephemeral-access/pkg/log"
	"github.com/argoproj-labs/argocd-ephemeral-access/pkg/plugin"
	goPlugin "github.com/hashicorp/go-plugin"
)

type SomePlugin struct {
	Logger               hclog.Logger
	sleepByAccessRequest map[string]time.Time
}

func (p *SomePlugin) Init() error {
	p.Logger.Info("This is a call to the Init method")
	sleepByAccessRequest := make(map[string]time.Time)
	p.sleepByAccessRequest = sleepByAccessRequest
	return nil
}

func (p *SomePlugin) GrantAccess(ar *api.AccessRequest, app *argocd.Application) (*plugin.GrantResponse, error) {
	p.Logger.Info("This is a call to the GrantAccess method")
	addedAt, ok := p.sleepByAccessRequest[ar.GetName()]
	if !ok {
		p.sleepByAccessRequest[ar.GetName()] = time.Now()
		return pending()
	}
	if time.Now().Before(addedAt.Add(time.Second * 30)) {
		return pending()
	}
	if isMinuteEven() {
		return granted()
	}
	return denied()
}

func isMinuteEven() bool {
	currentTime := time.Now().Format(time.TimeOnly)
	parts := strings.Split(currentTime, ":")
	min, _ := strconv.Atoi(parts[1])
	return min%2 == 0
}

func pending() (*plugin.GrantResponse, error) {
	return &plugin.GrantResponse{
		Status:  plugin.GrantStatusPending,
		Message: "Pending access by the example plugin",
	}, nil
}

func granted() (*plugin.GrantResponse, error) {
	return &plugin.GrantResponse{
		Status:  plugin.GrantStatusGranted,
		Message: "Granted access by the example plugin",
	}, nil
}

func denied() (*plugin.GrantResponse, error) {
	return &plugin.GrantResponse{
		Status:  plugin.GrantStatusDenied,
		Message: "Denied access by the example plugin",
	}, nil
}

func (p *SomePlugin) RevokeAccess(ar *api.AccessRequest, app *argocd.Application) (*plugin.RevokeResponse, error) {
	p.Logger.Info("This is a call to the RevokeAccess method")
	delete(p.sleepByAccessRequest, ar.GetName())
	return &plugin.RevokeResponse{
		Status:  plugin.RevokeStatusRevoked,
		Message: "Revoked access by the example plugin",
	}, nil
}

func main() {
	logger, err := log.NewPluginLogger()
	if err != nil {
		panic(fmt.Sprintf("Error creating plugin logger: %s", err))
	}

	p := &SomePlugin{
		Logger: logger,
	}
	srvConfig := plugin.NewServerConfig(p, logger)
	goPlugin.Serve(srvConfig)
}
