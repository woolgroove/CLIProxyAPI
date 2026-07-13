package auth

import (
	"context"
	"net/http"
	"testing"

	internalconfig "github.com/router-for-me/CLIProxyAPI/v7/internal/config"
)

func TestManagerMarkResult_DisablesXAIAuthForPermissionDenied(t *testing.T) {
	autoDisable := true
	manager := NewManager(nil, nil, nil)
	manager.SetConfig(&internalconfig.Config{XAI: internalconfig.XAIConfig{
		AutoDisablePermissionDenied: &autoDisable,
	}})

	auth := &Auth{
		ID:       "xai-permission-denied",
		Provider: "xai",
		Metadata: map[string]any{"type": "xai"},
	}
	if _, err := manager.Register(WithSkipPersist(context.Background()), auth); err != nil {
		t.Fatalf("register auth: %v", err)
	}

	reason := `{"code":"permission-denied","error":"Access to the chat endpoint is denied."}`
	manager.MarkResult(context.Background(), Result{
		AuthID:   auth.ID,
		Provider: "xai",
		Model:    "grok-4.5",
		Success:  false,
		Error: &Error{
			HTTPStatus: http.StatusForbidden,
			Message:    reason,
		},
	})

	updated, ok := manager.GetByID(auth.ID)
	if !ok || updated == nil {
		t.Fatal("updated auth not found")
	}
	if !updated.Disabled || updated.Status != StatusDisabled {
		t.Fatalf("auth disabled state = disabled:%t status:%s, want disabled", updated.Disabled, updated.Status)
	}
	if got, _ := updated.Metadata["disabled"].(bool); !got {
		t.Fatalf("metadata disabled = %#v, want true", updated.Metadata["disabled"])
	}
	if got, _ := updated.Metadata["disabled_reason"].(string); got != reason {
		t.Fatalf("disabled_reason = %q, want %q", got, reason)
	}
}

func TestManagerMarkResult_KeepsUnknownXAI403Enabled(t *testing.T) {
	autoDisable := true
	manager := NewManager(nil, nil, nil)
	manager.SetConfig(&internalconfig.Config{XAI: internalconfig.XAIConfig{
		AutoDisablePermissionDenied: &autoDisable,
	}})

	auth := &Auth{ID: "xai-unknown-forbidden", Provider: "xai", Metadata: map[string]any{"type": "xai"}}
	if _, err := manager.Register(WithSkipPersist(context.Background()), auth); err != nil {
		t.Fatalf("register auth: %v", err)
	}
	manager.MarkResult(context.Background(), Result{
		AuthID:   auth.ID,
		Provider: "xai",
		Model:    "grok-4.5",
		Success:  false,
		Error:    &Error{HTTPStatus: http.StatusForbidden, Message: `{"error":"temporary rejection"}`},
	})

	updated, ok := manager.GetByID(auth.ID)
	if !ok || updated == nil {
		t.Fatal("updated auth not found")
	}
	if updated.Disabled {
		t.Fatal("unknown xAI 403 must not disable the auth")
	}
}
