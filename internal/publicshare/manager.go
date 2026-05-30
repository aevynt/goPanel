package publicshare

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// max concurrent thumbnail generations — e2-micro can't handle more
var thumbSem = make(chan struct{}, 2)

type Share struct {
	ID          string `json:"id"`
	Folder      string `json:"folder"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type FileEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	IsDir     bool   `json:"is_dir"`
	ModTime   string `json:"mod_time"`
	Thumbnail string `json:"thumbnail,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
}

type Manager struct {
	publicDir string
}

func New(publicDir string) *Manager {
	os.MkdirAll(publicDir, 0755)
	return &Manager{publicDir: publicDir}
}

func (m *Manager) SharePath(share *Share) string {
	return filepath.Join(m.publicDir, share.Folder)
}

func (m *Manager) CreateShare(folder, title, description string) (*Share, error) {
	folderPath := filepath.Join(m.publicDir, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, fmt.Errorf("create folder: %w", err)
	}
	share := &Share{
		ID:          folder,
		Folder:      folder,
		Type:        "album",
		Title:       title,
		Description: description,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	return share, nil
}

func (m *Manager) DeleteShare(share *Share) error {
	return os.RemoveAll(m.SharePath(share))
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

func isImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return imageExts[ext]
}

func (m *Manager) safeFolder(folder string) (string, error) {
	full := filepath.Clean(filepath.Join(m.publicDir, folder))
	if full != m.publicDir && !strings.HasPrefix(full, m.publicDir+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal denied: %s", folder)
	}
	return full, nil
}

func (m *Manager) safePath(folder, name string) (string, error) {
	folderPath, err := m.safeFolder(folder)
	if err != nil {
		return "", err
	}
	full := filepath.Clean(filepath.Join(folderPath, name))
	if full != folderPath && !strings.HasPrefix(full, folderPath+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal denied: %s", name)
	}
	return full, nil
}

func (m *Manager) ListFiles(folder string) ([]FileEntry, error) {
	folderPath, err := m.safeFolder(folder)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	var files []FileEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		fe := FileEntry{
			Name:    info.Name(),
			Path:    filepath.ToSlash(filepath.Join("/", folder, info.Name())),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().UTC().Format(time.RFC3339),
		}
		if isImage(info.Name()) {
			fe.Thumbnail = "/p/" + folder + "/" + info.Name() + "?thumb=1"
			fe.MimeType = "image/" + strings.TrimPrefix(filepath.Ext(info.Name()), ".")
		}
		files = append(files, fe)
	}
	return files, nil
}

func (m *Manager) thumbDir(folder string) string {
	folderPath, err := m.safeFolder(folder)
	if err != nil {
		return filepath.Join(m.publicDir, ".invalid_thumb")
	}
	return filepath.Join(folderPath, ".thumb")
}

func (m *Manager) cachedThumbPath(folder, name string) string {
	td := m.thumbDir(folder)
	cleanedName := filepath.Base(name)
	return filepath.Join(td, cleanedName+".jpg")
}

// GenerateThumbnailFile generates a thumbnail and saves it to disk cache.
// Returns the path to the cached thumbnail.
func (m *Manager) GenerateThumbnailFile(folder, name string) (string, error) {
	cachePath := m.cachedThumbPath(folder, name)

	// already cached
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	data, err := m.generateThumbnailBytes(folder, name)
	if err != nil {
		return "", err
	}

	os.MkdirAll(m.thumbDir(folder), 0755)
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return "", err
	}
	return cachePath, nil
}

// GenerateThumbnail returns thumbnail bytes, preferring disk cache.
// Concurrent generation is limited to 2 to avoid CPU spikes on small servers.
func (m *Manager) GenerateThumbnail(folder, name string) ([]byte, error) {
	cachePath := m.cachedThumbPath(folder, name)
	// fast path: already cached
	if data, err := os.ReadFile(cachePath); err == nil {
		return data, nil
	}

	thumbSem <- struct{}{}        // acquire (block if 2 already running)
	defer func() { <-thumbSem }() // release

	// double-check after acquiring semaphore
	if data, err := os.ReadFile(cachePath); err == nil {
		return data, nil
	}

	data, err := m.generateThumbnailBytes(folder, name)
	if err != nil {
		return nil, err
	}

	// save to disk cache for future requests
	os.MkdirAll(m.thumbDir(folder), 0755)
	os.WriteFile(cachePath, data, 0644)
	return data, nil
}

// DeleteCachedThumbnail removes a cached thumbnail if it exists.
func (m *Manager) DeleteCachedThumbnail(folder, name string) {
	os.Remove(m.cachedThumbPath(folder, name))
}

func (m *Manager) generateThumbnailBytes(folder, name string) ([]byte, error) {
	srcPath, err := m.safePath(folder, name)
	if err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(name))

	src, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(src)
	case ".png":
		img, err = png.Decode(src)
	default:
		img, _, err = image.Decode(src)
	}
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	max := 640
	var tw, th int
	if w > h {
		tw = max
		th = int(math.Round(float64(h) * float64(max) / float64(w)))
	} else {
		th = max
		tw = int(math.Round(float64(w) * float64(max) / float64(h)))
	}

	dst := image.NewNRGBA(image.Rect(0, 0, tw, th))
	resizeNRGBA(dst, img)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func resizeNRGBA(dst *image.NRGBA, src image.Image) {
	b := dst.Bounds()
	sb := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx := float64(x-b.Min.X) * float64(sb.Dx()) / float64(b.Dx())
			sy := float64(y-b.Min.Y) * float64(sb.Dy()) / float64(b.Dy())
			r, g, b, a := sampleBilinear(src, sx+float64(sb.Min.X), sy+float64(sb.Min.Y))
			off := dst.PixOffset(x, y)
			dst.Pix[off+0] = r
			dst.Pix[off+1] = g
			dst.Pix[off+2] = b
			dst.Pix[off+3] = a
		}
	}
}

func sampleBilinear(src image.Image, x, y float64) (uint8, uint8, uint8, uint8) {
	fx := int(math.Floor(x))
	fy := int(math.Floor(y))
	cx := fx + 1
	cy := fy + 1

	xf := x - float64(fx)
	yf := y - float64(fy)

	r1, g1, b1, a1 := src.At(fx, fy).RGBA()
	r2, g2, b2, a2 := src.At(cx, fy).RGBA()
	r3, g3, b3, a3 := src.At(fx, cy).RGBA()
	r4, g4, b4, a4 := src.At(cx, cy).RGBA()

	r := uint32(bilinear(float64(r1), float64(r2), float64(r3), float64(r4), xf, yf))
	g := uint32(bilinear(float64(g1), float64(g2), float64(g3), float64(g4), xf, yf))
	b := uint32(bilinear(float64(b1), float64(b2), float64(b3), float64(b4), xf, yf))
	a := uint32(bilinear(float64(a1), float64(a2), float64(a3), float64(a4), xf, yf))

	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)
}

func bilinear(tl, tr, bl, br, xf, yf float64) float64 {
	top := tl + (tr-tl)*xf
	bottom := bl + (br-bl)*xf
	return top + (bottom-top)*yf
}
