package events_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBusinessModulesDoNotImportEachOther(t *testing.T) {
	root := filepath.Join("..", "modules")
	forbidden := map[string][]string{
		"tires":     {"maintenance"},
		"fuel":      {"financial"},
		"checklist": {"maintenance"},
		"financial": {"fuel"},
		"ciot":      {"fleet"},
	}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) < 3 {
			return nil
		}
		module := parts[2]
		targets := forbidden[module]
		if len(targets) == 0 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		for _, target := range targets {
			if strings.Contains(content, "internal/modules/"+target) {
				t.Fatalf("%s imports forbidden module %s", path, target)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk modules: %v", err)
	}
}
