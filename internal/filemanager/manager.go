package filemanager

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
}

type Manager struct {
	root string
}

func New(root string) *Manager {
	return &Manager{root: filepath.Clean(root)}
}

func (m *Manager) safePath(path string) (string, error) {
	full := filepath.Clean(filepath.Join(m.root, path))
	if full != m.root && !strings.HasPrefix(full, m.root+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal denied: %s", path)
	}
	return full, nil
}

func (m *Manager) List(dir string) ([]FileInfo, error) {
	safePath, err := m.safePath(dir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	files := make([]FileInfo, 0)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, m.toFileInfo(info, dir))
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})
	return files, nil
}

func (m *Manager) Stat(path string) (*FileInfo, error) {
	safePath, err := m.safePath(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(safePath)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	dir := filepath.Dir(path)
	fi := m.toFileInfo(info, dir)
	fi.Path = path
	return &fi, nil
}

func (m *Manager) Read(path string) ([]byte, error) {
	safePath, err := m.safePath(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(safePath)
}

func (m *Manager) Write(path string, data []byte) error {
	safePath, err := m.safePath(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(safePath), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	return os.WriteFile(safePath, data, 0644)
}

func (m *Manager) Mkdir(path string) error {
	safePath, err := m.safePath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(safePath, 0755)
}

func (m *Manager) Remove(path string) error {
	safePath, err := m.safePath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(safePath)
}

func (m *Manager) Rename(oldPath, newPath string) error {
	safeOld, err := m.safePath(oldPath)
	if err != nil {
		return err
	}
	safeNew, err := m.safePath(newPath)
	if err != nil {
		return err
	}
	return os.Rename(safeOld, safeNew)
}

func (m *Manager) Upload(path string, reader io.Reader) error {
	safePath, err := m.safePath(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(safePath), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	dst, err := os.Create(safePath)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, reader)
	return err
}

func (m *Manager) toFileInfo(info fs.FileInfo, baseDir string) FileInfo {
	return FileInfo{
		Name:    info.Name(),
		Path:    path.Join(baseDir, info.Name()),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		Mode:    info.Mode().String(),
		ModTime: info.ModTime().Format("2006-01-02T15:04:05Z"),
	}
}
