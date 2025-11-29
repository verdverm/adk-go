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

package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
	"google.golang.org/adk/session"
)

// FilesystemService returns a new filesystem-backed implementation of the memory service.
// It persists sessions to disk as JSON files in the specified base directory.
func FilesystemService(baseDir string) (Service, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	ims, ok := InMemoryService().(*inMemoryService)
	if !ok {
		return nil, fmt.Errorf("failed to cast in-memory service")
	}

	return &filesystemService{
		inMemoryService: ims,
		baseDir:         baseDir,
	}, nil
}

type filesystemService struct {
	*inMemoryService
	baseDir string
}

type persistentEntry struct {
	Content   *genai.Content `json:"content"`
	Author    string         `json:"author"`
	Timestamp time.Time      `json:"timestamp"`
}

func (s *filesystemService) AddSession(ctx context.Context, curSession session.Session) error {
	// First, update the in-memory store.
	if err := s.inMemoryService.AddSession(ctx, curSession); err != nil {
		return err
	}

	// Retrieve the updated values from the in-memory store.
	k := key{
		appName: curSession.AppName(),
		userID:  curSession.UserID(),
	}
	sid := sessionID(curSession.ID())

	s.mu.RLock()
	sessionMap, ok := s.store[k]
	if !ok {
		s.mu.RUnlock()
		return fmt.Errorf("session not found in memory after adding")
	}
	values, ok := sessionMap[sid]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("session values not found in memory after adding")
	}

	// Prepare data for persistence.
	entries := make([]persistentEntry, 0, len(values))
	for _, v := range values {
		entries = append(entries, persistentEntry{
			Content:   v.content,
			Author:    v.author,
			Timestamp: v.timestamp,
		})
	}

	// Determine file path.
	// Structure: <baseDir>/<appName>/<userID>/<sessionID>/memory/session.json
	// Note: The prompt mentioned <path-arg>.md but also "simple json". We use json.
	dirPath := filepath.Join(s.baseDir, curSession.AppName(), curSession.UserID(), curSession.ID(), "memory")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, "session.json")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("failed to encode session memory: %w", err)
	}

	return nil
}

// Load walks the base directory and loads all session memories into the in-memory store.
// This is not part of the Service interface but is required for initialization.
func (s *filesystemService) Load(ctx context.Context) error {
	return filepath.WalkDir(s.baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) != "session.json" {
			return nil
		}

		// Parse path to extract metadata if needed, but we rely on file content + path structure.
		// Path: <baseDir>/<appName>/<userID>/<sessionID>/memory/session.json
		relPath, err := filepath.Rel(s.baseDir, path)
		if err != nil {
			return err
		}
		parts := strings.Split(relPath, string(os.PathSeparator))
		if len(parts) < 4 {
			// unexpected structure, skip
			return nil
		}
		appName := parts[0]
		userID := parts[1]
		sessionIDStr := parts[2]

		// Read and decode the file.
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		var entries []persistentEntry
		if err := json.NewDecoder(file).Decode(&entries); err != nil {
			return fmt.Errorf("failed to decode file %s: %w", path, err)
		}

		// Convert back to internal values.
		values := make([]value, 0, len(entries))
		for _, e := range entries {
			// Re-compute words.
			words := make(map[string]struct{})
			if e.Content != nil {
				for _, part := range e.Content.Parts {
					if part.Text == "" {
						continue
					}
					// extractWords is defined in inmemory.go (same package)
					for k, v := range extractWords(part.Text) {
						words[k] = v
					}
				}
			}

			values = append(values, value{
				content:   e.Content,
				author:    e.Author,
				timestamp: e.Timestamp,
				words:     words,
			})
		}

		// Update store.
		k := key{
			appName: appName,
			userID:  userID,
		}
		sid := sessionID(sessionIDStr)

		s.mu.Lock()
		v, ok := s.store[k]
		if !ok {
			v = make(map[sessionID][]value)
			s.store[k] = v
		}
		v[sid] = values
		s.mu.Unlock()

		return nil
	})
}
