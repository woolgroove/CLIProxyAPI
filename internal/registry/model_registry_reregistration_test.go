package registry

import "testing"

func TestRegisterClientPreservesTransientStateForUnchangedBinding(t *testing.T) {
	r := newTestModelRegistry()
	models := []*ModelInfo{{ID: "grok-4.5"}}
	r.RegisterClient("client-1", "xai", models)
	r.SetModelQuotaExceeded("client-1", "grok-4.5")
	r.SuspendClientModel("client-1", "grok-4.5", "quota")

	r.RegisterClient("client-1", "xai", models)

	r.mutex.RLock()
	registration := r.models["grok-4.5"]
	_, quotaExceeded := registration.QuotaExceededClients["client-1"]
	reason, suspended := registration.SuspendedClients["client-1"]
	r.mutex.RUnlock()
	if !quotaExceeded {
		t.Fatal("re-registration cleared quota-exceeded state for unchanged binding")
	}
	if !suspended || reason != "quota" {
		t.Fatalf("re-registration cleared suspension for unchanged binding: suspended=%v reason=%q", suspended, reason)
	}
}
