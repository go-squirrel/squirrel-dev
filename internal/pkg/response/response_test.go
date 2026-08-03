package response

import "testing"

func TestLegacyCommonResponseContract(t *testing.T) {
	Init()

	success := Success("health")
	if success.Code != 0 || success.Message != "success" || success.Data != "health" {
		t.Fatalf("unexpected success response: %#v", success)
	}

	sqlError := Error(ErrSQL)
	if sqlError.Code != 50000 || sqlError.Message != "sql error" {
		t.Fatalf("unexpected SQL error response: %#v", sqlError)
	}
}
