package icon

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Extractor interface {
	Extract(ctx context.Context, source string, dest string) error
}

type Cache struct {
	root      string
	extractor Extractor
}

func NewCache(root string, extractor Extractor) *Cache {
	return &Cache{root: root, extractor: extractor}
}

func (c *Cache) Ensure(ctx context.Context, source string, force bool) (string, error) {
	if strings.TrimSpace(source) == "" {
		return "", nil
	}
	if c.extractor == nil {
		return "", ErrUnsupported
	}

	if err := os.MkdirAll(c.root, 0o755); err != nil {
		return "", err
	}

	dest := filepath.Join(c.root, hashPath(source)+".png")
	if !force {
		if _, err := os.Stat(dest); err == nil {
			return dest, nil
		}
	}
	if force {
		_ = os.Remove(dest)
	}

	if err := c.extractor.Extract(ctx, source, dest); err != nil {
		return "", err
	}

	return dest, nil
}

func hashPath(source string) string {
	hash := sha1.Sum([]byte(strings.ToLower(strings.TrimSpace(source))))
	return hex.EncodeToString(hash[:])
}

func CopyFile(source string, dest string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	if _, err := io.Copy(output, input); err != nil {
		return err
	}

	if err := output.Sync(); err != nil {
		return err
	}

	return nil
}

func ValidateSource(source string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("source is empty")
	}
	if _, err := os.Stat(source); err != nil {
		return err
	}
	return nil
}
