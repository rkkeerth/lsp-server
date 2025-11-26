package jsonrpc

import (
	"encoding/json"
	"testing"
)

func TestNewRequest(t *testing.T) {
	params := map[string]string{"key": "value"}
	req, err := NewRequest(1, "test/method", params)

	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	if req.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC version %s, got %s", JSONRPCVersion, req.JSONRPC)
	}

	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %v", req.ID)
	}

	if req.Method != "test/method" {
		t.Errorf("Expected method 'test/method', got %s", req.Method)
	}

	var decodedParams map[string]string
	if err := json.Unmarshal(req.Params, &decodedParams); err != nil {
		t.Fatalf("Failed to unmarshal params: %v", err)
	}

	if decodedParams["key"] != "value" {
		t.Errorf("Expected params key=value, got %v", decodedParams)
	}
}

func TestNewResponse(t *testing.T) {
	result := map[string]string{"result": "success"}
	resp, err := NewResponse(1, result)

	if err != nil {
		t.Fatalf("NewResponse failed: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}

	var decodedResult map[string]string
	if err := json.Unmarshal(resp.Result, &decodedResult); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if decodedResult["result"] != "success" {
		t.Errorf("Expected result=success, got %v", decodedResult)
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse(1, InternalError, "test error")

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}

	if resp.Error == nil {
		t.Fatal("Expected error to be set")
	}

	if resp.Error.Code != InternalError {
		t.Errorf("Expected error code %d, got %d", InternalError, resp.Error.Code)
	}

	if resp.Error.Message != "test error" {
		t.Errorf("Expected error message 'test error', got %s", resp.Error.Message)
	}
}

func TestNewNotification(t *testing.T) {
	params := map[string]string{"key": "value"}
	notif, err := NewNotification("test/notification", params)

	if err != nil {
		t.Fatalf("NewNotification failed: %v", err)
	}

	if notif.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC version %s, got %s", JSONRPCVersion, notif.JSONRPC)
	}

	if notif.Method != "test/notification" {
		t.Errorf("Expected method 'test/notification', got %s", notif.Method)
	}

	var decodedParams map[string]string
	if err := json.Unmarshal(notif.Params, &decodedParams); err != nil {
		t.Fatalf("Failed to unmarshal params: %v", err)
	}

	if decodedParams["key"] != "value" {
		t.Errorf("Expected params key=value, got %v", decodedParams)
	}
}
