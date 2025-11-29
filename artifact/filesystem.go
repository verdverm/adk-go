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

package artifact

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/genai"
)

// FilesystemService returns a new filesystem-backed implementation of the artifact service.
// It persists artifacts to disk as JSON files in the specified base directory.
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

var versionFileRegex = regexp.MustCompile(`^(.*)\.v(\d+)\.json$`)

// Save persists the artifact to disk after saving it to memory.
func (s *filesystemService) Save(ctx context.Context, req *SaveRequest) (*SaveResponse, error) {
	// Delegate to in-memory service first to handle validation, versioning, and memory update.
	resp, err := s.inMemoryService.Save(ctx, req)
	if err != nil {
		return nil, err
	}

	appName, userID, sessionID, fileName := req.AppName, req.UserID, req.SessionID, req.FileName
	// Replicate the logic for user-scoped artifacts to determine the storage path.
	if fileHasUserNamespace(fileName) {
		sessionID = userScopedArtifactKey
	}

	// Structure: <baseDir>/<appName>/<userID>/<sessionID>/artifacts/<fileName>.v<version>.json
	dirPath := filepath.Join(s.baseDir, appName, userID, sessionID, "artifacts")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// We use the version returned by the in-memory service.
	filenameOnDisk := fmt.Sprintf("%s.v%d.json", fileName, resp.Version)
	filePath := filepath.Join(dirPath, filenameOnDisk)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(req.Part); err != nil {
		return nil, fmt.Errorf("failed to encode artifact: %w", err)
	}

	return resp, nil
}

// Delete removes the artifact from disk after removing it from memory.
func (s *filesystemService) Delete(ctx context.Context, req *DeleteRequest) error {
	// Delegate to in-memory service first.
	if err := s.inMemoryService.Delete(ctx, req); err != nil {
		return err
	}

	appName, userID, sessionID, fileName := req.AppName, req.UserID, req.SessionID, req.FileName
	if fileHasUserNamespace(fileName) {
		sessionID = userScopedArtifactKey
	}

	dirPath := filepath.Join(s.baseDir, appName, userID, sessionID, "artifacts")

	// If version is specified, delete that specific file.
	if req.Version != 0 {
		filenameOnDisk := fmt.Sprintf("%s.v%d.json", fileName, req.Version)
		filePath := filepath.Join(dirPath, filenameOnDisk)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", filePath, err)
		}
		return nil
	}

	// If version is 0, delete all versions of the file.
	// We walk the directory and delete matches.
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		matches := versionFileRegex.FindStringSubmatch(name)
		if matches != nil && matches[1] == fileName {
			filePath := filepath.Join(dirPath, name)
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to delete file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// LoadFromDisk walks the base directory and loads all artifacts into the in-memory store.
func (s *filesystemService) LoadFromDisk(ctx context.Context) error {
	return filepath.WalkDir(s.baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Check if it looks like an artifact file
		if !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		// Parse path: <baseDir>/<appName>/<userID>/<sessionID>/artifacts/<fileName>.v<ver>.json
		relPath, err := filepath.Rel(s.baseDir, path)
		if err != nil {
			return err
		}
		parts := strings.Split(relPath, string(os.PathSeparator))
		// Minimum parts: app, user, session, "artifacts", file
		if len(parts) < 5 {
			return nil
		}

		// Check if "artifacts" segment exists
		if parts[3] != "artifacts" {
			return nil
		}

		appName := parts[0]
		userID := parts[1]
		sessionID := parts[2]
		
		// Reconstruct filename from the rest (in case filename had slashes, though flatten strategy handles it via Join?)
		// Wait, filepath.Join in Save uses path separators.
		// If fileName had slashes "a/b.txt", Save joined `.../artifacts` + `a/b.txt.v1.json`.
		// So it created `.../artifacts/a/b.txt.v1.json`.
		// So `parts` will look like: ..., "artifacts", "a", "b.txt.v1.json".
		// We need to join everything after "artifacts".
		
		rawFileName := filepath.Join(parts[4:]...)
		
		// Parse version and actual filename
		matches := versionFileRegex.FindStringSubmatch(rawFileName)
		if matches == nil {
			return nil
		}
		
		fileName := matches[1]
		versionStr := matches[2]
		version, err := strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			return nil // Skip invalid version
		}

		// Read content
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		var part genai.Part
		if err := json.NewDecoder(file).Decode(&part); err != nil {
			return fmt.Errorf("failed to decode artifact %s: %w", path, err)
		}

		// Store in memory
		s.inMemoryService.set(appName, userID, sessionID, fileName, version, &part)

		return nil
	})
}
