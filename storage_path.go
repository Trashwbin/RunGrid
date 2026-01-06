package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const storageConfigFileName = "storage.json"

type storageConfig struct {
	DataRoot    string `json:"data_root"`
	MigrateFrom string `json:"migrate_from,omitempty"`
}

func dataRootPath(appName string, migrate bool) (string, error) {
	cfg, configRoot, err := loadStorageConfig(appName)
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(cfg.DataRoot)
	if root == "" {
		root = configRoot
	}
	root = filepath.Clean(root)

	if migrate {
		migrateFrom := strings.TrimSpace(cfg.MigrateFrom)
		if migrateFrom != "" && filepath.Clean(migrateFrom) != root {
			if err := migrateDataRoot(migrateFrom, root); err != nil {
				return filepath.Clean(migrateFrom), nil
			}
			cfg.MigrateFrom = ""
			cfg.DataRoot = root
			if err := saveStorageConfig(appName, cfg); err != nil {
				return "", err
			}
		}
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}
	return root, nil
}

func appDataPath(appName string) (string, error) {
	root, err := dataRootPath(appName, true)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "rungrid.db"), nil
}

func prepareDataRootChange(appName, nextRoot string) error {
	info, err := os.Stat(nextRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("target path is not a directory")
	}
	conflict, err := dataRootHasConflicts(nextRoot)
	if err != nil {
		return err
	}
	if conflict {
		return fmt.Errorf("target path already contains RunGrid data")
	}
	return nil
}

func configRootPath(appName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		configDir = "."
	}
	root := filepath.Join(configDir, appName)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}
	return root, nil
}

func loadStorageConfig(appName string) (storageConfig, string, error) {
	configRoot, err := configRootPath(appName)
	if err != nil {
		return storageConfig{}, "", err
	}
	configPath := filepath.Join(configRoot, storageConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return storageConfig{}, configRoot, nil
		}
		return storageConfig{}, configRoot, err
	}
	if len(data) == 0 {
		return storageConfig{}, configRoot, nil
	}
	var cfg storageConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return storageConfig{}, configRoot, err
	}
	return cfg, configRoot, nil
}

func saveStorageConfig(appName string, cfg storageConfig) error {
	configRoot, err := configRootPath(appName)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configRoot, storageConfigFileName), data, 0o644)
}

func migrateDataRoot(oldRoot, newRoot string) error {
	oldRoot = filepath.Clean(oldRoot)
	newRoot = filepath.Clean(newRoot)
	if oldRoot == newRoot {
		return nil
	}
	if err := os.MkdirAll(newRoot, 0o755); err != nil {
		return err
	}
	conflict, err := dataRootHasConflicts(newRoot)
	if err != nil {
		return err
	}
	if conflict {
		return fmt.Errorf("target path already contains RunGrid data")
	}

	for _, entry := range []string{"rungrid.db", "icons", "rules"} {
		source := filepath.Join(oldRoot, entry)
		if _, err := os.Stat(source); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		dest := filepath.Join(newRoot, entry)
		if err := movePath(source, dest); err != nil {
			return err
		}
	}
	return nil
}

func dataRootHasConflicts(root string) (bool, error) {
	for _, entry := range []string{"rungrid.db", "icons", "rules"} {
		path := filepath.Join(root, entry)
		if _, err := os.Stat(path); err == nil {
			return true, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
	}
	return false, nil
}

func movePath(source, dest string) error {
	if err := os.Rename(source, dest); err == nil {
		return nil
	}
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if err := copyDir(source, dest); err != nil {
			return err
		}
		return os.RemoveAll(source)
	}
	if err := copyFile(source, dest, info.Mode()); err != nil {
		return err
	}
	return os.Remove(source)
}

func copyDir(source, dest string) error {
	return filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(source, dest string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
