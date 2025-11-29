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

package artifact_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genai"
	"google.golang.org/adk/artifact"
)

func TestFilesystemService_Save_Load_Delete(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Create service and save artifact
	svc, err := artifact.FilesystemService(tempDir)
	if err != nil {
		t.Fatalf("FilesystemService() error = %v", err)
	}

	part := &genai.Part{Text: "hello world"}
	req := &artifact.SaveRequest{
		AppName:   "app1",
		UserID:    "user1",
		SessionID: "sess1",
		FileName:  "foo.txt",
		Part:      part,
	}

	resp, err := svc.Save(t.Context(), req)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if resp.Version != 1 {
		t.Errorf("Expected version 1, got %d", resp.Version)
	}

	// 2. Verify file exists
	expectedPath := filepath.Join(tempDir, "app1", "user1", "sess1", "artifacts", "foo.txt.v1.json")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("Expected file at %s, got error: %v", expectedPath, err)
	}

	// 3. Create new service and load
	svc2, err := artifact.FilesystemService(tempDir)
	if err != nil {
		t.Fatalf("FilesystemService() error = %v", err)
	}
	fsSvc2, ok := svc2.(interface{ LoadFromDisk(context.Context) error })
	if !ok {
		t.Fatalf("Service does not implement LoadFromDisk method")
	}
	if err := fsSvc2.LoadFromDisk(t.Context()); err != nil {
		t.Fatalf("LoadFromDisk() error = %v", err)
	}

	// 4. Verify artifact exists in new service
	loadResp, err := svc2.Load(t.Context(), &artifact.LoadRequest{
		AppName:   "app1",
		UserID:    "user1",
		SessionID: "sess1",
		FileName:  "foo.txt",
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if diff := cmp.Diff(part, loadResp.Part); diff != "" {
		t.Errorf("Artifact mismatch (-want +got):\n%s", diff)
	}

	// 5. Delete artifact
	if err := svc2.Delete(t.Context(), &artifact.DeleteRequest{
		AppName:   "app1",
		UserID:    "user1",
		SessionID: "sess1",
		FileName:  "foo.txt",
	}); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(expectedPath); !os.IsNotExist(err) {
		t.Errorf("Expected file to be deleted, got error: %v", err)
	}
}
