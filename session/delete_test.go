// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package session

import (
	"context"
	"testing"
)

func TestStateDelete(t *testing.T) {
	s := InMemoryService()
	ctx := context.Background()
	sessResp, err := s.Create(ctx, &CreateRequest{
		AppName:   "app",
		UserID:    "user",
		SessionID: "s1",
		State: map[string]any{
			"foo": "bar",
		},
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	sess := sessResp.Session

	state := sess.State()
	val, err := state.Get("foo")
	if err != nil || val != "bar" {
		t.Errorf("expected foo=bar, got %v, %v", val, err)
	}

	// Try to "delete" by setting to nil
	err = state.Set("foo", nil)
	if err != nil {
		t.Errorf("Set(foo, nil) failed: %v", err)
	}

	val, err = state.Get("foo")
	if err != ErrStateKeyNotExist {
		t.Errorf("expected ErrStateKeyNotExist after Set(foo, nil), got %v, %v", val, err)
	}
}

func TestAppendEventDelete(t *testing.T) {
	s := InMemoryService()
	ctx := context.Background()
	sessResp, err := s.Create(ctx, &CreateRequest{
		AppName:   "app",
		UserID:    "user",
		SessionID: "s1",
		State: map[string]any{
			"foo": "bar",
		},
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	sess := sessResp.Session

	// Simulate Turn 1: Tool deletes foo
	event := &Event{
		Actions: EventActions{
			StateDelta: map[string]any{
				"foo": nil,
			},
		},
	}
	err = s.AppendEvent(ctx, sess, event)
	if err != nil {
		t.Fatalf("AppendEvent failed: %v", err)
	}

	// Check if foo is deleted from the session state
	val, err := sess.State().Get("foo")
	if err != ErrStateKeyNotExist {
		t.Errorf("expected ErrStateKeyNotExist after AppendEvent with foo=nil, got %v, %v", val, err)
	}
	
	// Check if it's also deleted when we Get the session again
	getResp, err := s.Get(ctx, &GetRequest{AppName: "app", UserID: "user", SessionID: "s1"})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	val, err = getResp.Session.State().Get("foo")
	if err != ErrStateKeyNotExist {
		t.Errorf("expected ErrStateKeyNotExist after fresh Get, got %v, %v", val, err)
	}
}
