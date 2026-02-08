package cron

const (
	uriScriptResults = "/scripts/receive-result"
	uriAppReport     = "/deployment/report"
)

var composeProjects []struct {
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	ConfigFiles string `json:"ConfigFiles"`
}
