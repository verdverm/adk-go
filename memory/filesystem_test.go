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

package memory_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genai"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
)

func TestFilesystemService_AddSession_And_Load(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Create service and add a session
	svc, err := memory.FilesystemService(tempDir)
	if err != nil {
		t.Fatalf("FilesystemService() error = %v", err)
	}

	sess := makeSession(t, "app1", "user1", "sess1", []*session.Event{
		{
			Author: "user1",
			LLMResponse: model.LLMResponse{
				Content: genai.NewContentFromText("The Quick brown fox", genai.RoleUser),
			},
			Timestamp: must(time.Parse(time.RFC3339, "2023-10-01T10:00:00Z")),
		},
	})

	if err := svc.AddSession(t.Context(), sess); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	// 2. Verify file exists
	expectedPath := filepath.Join(tempDir, "app1", "user1", "sess1", "memory", "session.json")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("Expected file at %s, got error: %v", expectedPath, err)
	}

	// 3. Create a new service and load from disk
	svc2, err := memory.FilesystemService(tempDir)
	if err != nil {
		t.Fatalf("FilesystemService() error = %v", err)
	}

	// We need to cast to call Load since it's not in Service interface
	fsSvc2, ok := svc2.(interface {
		Load(context.Context) error
		Search(context.Context, *memory.SearchRequest) (*memory.SearchResponse, error)
	})
	if !ok {
		t.Fatalf("Service does not implement Load method")
	}

	if err := fsSvc2.Load(t.Context()); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// 4. Verify search works on loaded service
	req := &memory.SearchRequest{
		AppName: "app1",
		UserID:  "user1",
		Query:   "quick",
	}
	got, err := fsSvc2.Search(t.Context(), req)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got.Memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(got.Memories))
	} else {
		want := memory.Entry{
			Content:   genai.NewContentFromText("The Quick brown fox", genai.RoleUser),
			Author:    "user1",
			Timestamp: must(time.Parse(time.RFC3339, "2023-10-01T10:00:00Z")),
		}
		if diff := cmp.Diff(want, got.Memories[0]); diff != "" {
			t.Errorf("Memory mismatch (-want +got):\n%s", diff)
		}
	}
}
