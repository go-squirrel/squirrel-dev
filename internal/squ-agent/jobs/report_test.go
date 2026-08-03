package jobs

import (
	"context"
	"encoding/json"
	"testing"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	scriptDomain "squirrel-dev/internal/squ-agent/module/script/domain"
	"squirrel-dev/pkg/httpclient"
)

type scriptRepositoryStub struct {
	tasks  []scriptDomain.Task
	marked []uint
}

func (s *scriptRepositoryStub) Add(context.Context, *scriptDomain.Task) error { return nil }
func (s *scriptRepositoryStub) Get(context.Context, uint) (scriptDomain.Task, error) {
	return scriptDomain.Task{}, nil
}
func (s *scriptRepositoryStub) GetRunningTask(context.Context) (scriptDomain.Task, error) {
	return scriptDomain.Task{}, nil
}
func (s *scriptRepositoryStub) Update(context.Context, *scriptDomain.Task) error { return nil }
func (s *scriptRepositoryStub) GetUnreportedTasks(context.Context) ([]scriptDomain.Task, error) {
	return s.tasks, nil
}
func (s *scriptRepositoryStub) MarkAsReported(_ context.Context, id uint) error {
	s.marked = append(s.marked, id)
	return nil
}

type posterStub struct {
	url      string
	requests []scriptResultReport
	body     []byte
}

func (p *posterStub) Post(url string, value any, _ httpclient.Header) ([]byte, error) {
	p.url = url
	p.requests = append(p.requests, value.(scriptResultReport))
	return p.body, nil
}

func TestReportScriptResultsPreservesLegacyMarkingRule(t *testing.T) {
	response.Init()
	body, err := json.Marshal(response.Success("ok"))
	if err != nil {
		t.Fatal(err)
	}
	repository := &scriptRepositoryStub{tasks: []scriptDomain.Task{
		{ID: 1, TaskID: 10, ScriptID: 20, Status: "success", Output: "done"},
		{ID: 2, TaskID: 11, ScriptID: 21, Status: "failed", ErrorMsg: "boom"},
	}}
	poster := &posterStub{body: body}
	instance := &Jobs{
		config: &config.Config{Apiserver: config.Apiserver{
			Http: config.Http{Scheme: "http", Server: "api.example", BaseUri: "/api/v1"},
		}},
		scriptTasks: repository,
		http:        poster,
	}
	instance.reportScriptResults()

	if poster.url != "http://api.example/api/v1/scripts/receive-result" {
		t.Fatalf("url = %q", poster.url)
	}
	if len(poster.requests) != 2 || poster.requests[0].ScriptsID != 20 || poster.requests[1].ErrorMessage != "boom" {
		t.Fatalf("requests = %#v", poster.requests)
	}
	if len(repository.marked) != 1 || repository.marked[0] != 1 {
		t.Fatalf("only successful task should be marked: %#v", repository.marked)
	}
}

func TestLegacyMonitorFilters(t *testing.T) {
	for _, device := range []string{"loop0", "zram0", "dm-0"} {
		if !skipDisk(device) {
			t.Fatalf("%q should be skipped", device)
		}
	}
	if skipDisk("sda") {
		t.Fatal("sda should not be skipped")
	}
	for _, name := range []string{"lo", "Docker0", "veth123", "br-test"} {
		if !skipInterface(name) {
			t.Fatalf("%q should be skipped", name)
		}
	}
	if skipInterface("eth0") {
		t.Fatal("eth0 should not be skipped")
	}
}
