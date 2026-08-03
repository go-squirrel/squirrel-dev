package contract

import "testing"

func TestLegacyContractIsValid(t *testing.T) {
	value, err := LoadLegacy()
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(value); err != nil {
		t.Fatal(err)
	}
}

func TestLegacyRouteInventory(t *testing.T) {
	value, err := LoadLegacy()
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]int{
		"squ-agent":     24,
		"squ-apiserver": 51,
	}
	for serviceName, expectedCount := range expected {
		service, ok := value.Service(serviceName)
		if !ok {
			t.Fatalf("service %q is missing", serviceName)
		}
		if len(service.Routes) != expectedCount {
			t.Fatalf("%s route count = %d, want %d", serviceName, len(service.Routes), expectedCount)
		}
	}
}

func TestLegacyCompatibilitySentinels(t *testing.T) {
	value, err := LoadLegacy()
	if err != nil {
		t.Fatal(err)
	}

	message, ok := value.ResponseMessage("apiserver", 50000)
	if !ok || message != "sql error" {
		t.Fatalf("response code 50000 = %q, %v; want %q, true", message, ok, "sql error")
	}

	agentConfig, ok := value.Config("squ-agent")
	if !ok {
		t.Fatal("squ-agent config contract is missing")
	}
	if agentConfig.DefaultFile != "config/agent.yaml" {
		t.Fatalf("agent default config = %q, want %q", agentConfig.DefaultFile, "config/agent.yaml")
	}

	apiServer, ok := value.Service("squ-apiserver")
	if !ok {
		t.Fatal("squ-apiserver service contract is missing")
	}
	assertRoute(t, apiServer, "GET", "/api/v1/health")
	assertRoute(t, apiServer, "GET", "/api/v1/ws/server/:id")
	assertRoute(t, apiServer, "POST", "/api/v1/deployment/report")

	var healthSample *HTTPSample
	for i := range value.HTTPSamples {
		sample := &value.HTTPSamples[i]
		if sample.Service == "squ-apiserver" && sample.Path == "/api/v1/health" {
			healthSample = sample
			break
		}
	}
	if healthSample == nil {
		t.Fatal("apiserver health HTTP sample is missing")
	}
	if healthSample.Body != `{"code":0,"message":"success","data":"health"}` {
		t.Fatalf("health body = %q", healthSample.Body)
	}
}

func assertRoute(t *testing.T, service Service, method, path string) {
	t.Helper()
	for _, route := range service.Routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("%s does not contain %s %s", service.Name, method, path)
}
