package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/quotedprintable"
	"net/http"
	"strconv"
	"strings"
)

// decodeTransferEncoding decodes a MIME part body from its Content-Transfer-Encoding.
// IMAP returns BODY[part] still encoded; base64 and quoted-printable need decoding,
// the identity encodings (7bit/8bit/binary/absent) pass through unchanged.
func decodeTransferEncoding(data []byte, encoding string) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "base64":
		decoded, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data)))
		if err != nil {
			if stripped, e2 := base64.StdEncoding.DecodeString(stripWhitespace(data)); e2 == nil {
				return stripped, nil
			}
			if len(decoded) > 0 {
				return decoded, nil
			}
			return nil, fmt.Errorf("base64 decode failed: %w", err)
		}
		return decoded, nil
	case "quoted-printable":
		decoded, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(data)))
		if err != nil && len(decoded) == 0 {
			return nil, fmt.Errorf("quoted-printable decode failed: %w", err)
		}
		return decoded, nil
	default:
		return data, nil
	}
}

func stripWhitespace(data []byte) string {
	var b strings.Builder
	b.Grow(len(data))
	for _, c := range data {
		switch c {
		case ' ', '\t', '\r', '\n', '\v', '\f':
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

var inlineSafeImageTypes = map[string]bool{
	"image/png":                true,
	"image/jpeg":               true,
	"image/gif":                true,
	"image/webp":               true,
	"image/bmp":                true,
	"image/tiff":               true,
	"image/x-icon":             true,
	"image/vnd.microsoft.icon": true,
}

// dangerousTypes are content types a browser may execute in the page context.
// They are never reflected back as-is; we re-type them to octet-stream and force download.
var dangerousTypes = map[string]bool{
	"text/html":              true,
	"application/xhtml+xml":  true,
	"image/svg+xml":          true,
	"application/xml":        true,
	"text/xml":               true,
	"application/javascript": true,
	"text/javascript":        true,
}

func normalizeMediaType(declared string) string {
	mt := declared
	if parsed, _, err := mime.ParseMediaType(declared); err == nil {
		mt = parsed
	}
	return strings.ToLower(strings.TrimSpace(mt))
}

// resolveContentType decides the Content-Type and Content-Disposition mode to serve with.
// Inline rendering is only allowed for allowlisted raster image types whose bytes actually
// sniff as an image; everything else is forced to download, with executable types neutralized.
func resolveContentType(declared string, data []byte, inline bool) (contentType, disposition string) {
	mediaType := normalizeMediaType(declared)

	if inline && inlineSafeImageTypes[mediaType] && strings.HasPrefix(http.DetectContentType(data), "image/") {
		return mediaType, "inline"
	}

	if mediaType == "" || dangerousTypes[mediaType] {
		return "application/octet-stream", "attachment"
	}
	return mediaType, "attachment"
}

// writeAttachmentResponse serves decoded attachment bytes with hardened headers
// (nosniff, sandbox CSP, immutable caching, ETag) and an RFC 6266 filename.
func writeAttachmentResponse(w http.ResponseWriter, req *http.Request, data []byte, declaredType, filename string, inline bool, etag string) {
	contentType, disposition := resolveContentType(declaredType, data, inline)

	h := w.Header()
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; sandbox")
	h.Set("Cache-Control", "private, max-age=31536000, immutable")
	if etag != "" {
		h.Set("ETag", etag)
		if req.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	h.Set("Content-Type", contentType)
	h.Set("Content-Disposition", formatContentDisposition(disposition, filename))
	h.Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	if req.Method != http.MethodHead {
		w.Write(data)
	}
}

func attachmentETag(uid uint32, partID string) string {
	return fmt.Sprintf("\"%d-%s\"", uid, partID)
}

// formatContentDisposition builds an RFC 6266 header with an ASCII filename= fallback
// and an RFC 5987 filename*=UTF-8” form for non-ASCII names.
func formatContentDisposition(disposition, filename string) string {
	clean := cleanFilename(filename)
	ascii := asciiFilename(clean)
	if ascii == "" {
		ascii = "download"
	}
	header := fmt.Sprintf("%s; filename=%q", disposition, ascii)
	if encoded := encodeRFC5987(clean); strings.Contains(encoded, "%") {
		header += "; filename*=UTF-8''" + encoded
	}
	return header
}

// cleanFilename drops characters that are never valid in a served filename:
// control bytes (header-injection vectors) and path separators (traversal).
func cleanFilename(filename string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 || r == '/' || r == '\\' {
			return -1
		}
		return r
	}, filename)
}

func asciiFilename(filename string) string {
	return strings.Map(func(r rune) rune {
		if r == '"' || r > 126 {
			return '_'
		}
		return r
	}, filename)
}

func isAttrChar(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		return true
	}
	return strings.IndexByte("!#$&+-.^_`|~", c) >= 0
}

func encodeRFC5987(filename string) string {
	var b strings.Builder
	for i := 0; i < len(filename); i++ {
		c := filename[i]
		if isAttrChar(c) {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
