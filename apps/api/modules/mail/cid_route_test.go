package mail

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCIDParamIsDecoded(t *testing.T) {
	const stored = "image001.png@01DD0333.B53E0920"

	var got string
	r := chi.NewRouter()
	r.Get("/cid/{cid}", func(w http.ResponseWriter, req *http.Request) {
		cid := chi.URLParam(req, "cid")
		if decoded, err := url.PathUnescape(cid); err == nil {
			cid = decoded
		}
		got = cid
	})

	req := httptest.NewRequest(http.MethodGet, "/cid/image001.png%4001DD0333.B53E0920", nil)
	r.ServeHTTP(httptest.NewRecorder(), req)

	if got != stored {
		t.Fatalf("cid = %q, want %q", got, stored)
	}
}

func TestDecodeTransferEncoding(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x01, 0x02, 0x03}
	b64 := base64.StdEncoding.EncodeToString(png)
	wrapped := b64[:4] + "\r\n" + b64[4:]

	cases := []struct {
		name     string
		in       []byte
		encoding string
		want     []byte
	}{
		{"base64 with crlf", []byte(wrapped), "base64", png},
		{"base64 case-insensitive", []byte(b64), "BASE64", png},
		{"quoted-printable", []byte("Caf=C3=A9"), "quoted-printable", []byte("Café")},
		{"7bit passthrough", []byte("plain bytes"), "7bit", []byte("plain bytes")},
		{"empty encoding passthrough", []byte("raw"), "", []byte("raw")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := decodeTransferEncoding(c.in, c.encoding)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != string(c.want) {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestResolveContentType(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}
	html := []byte("<html><script>alert(1)</script></html>")

	cases := []struct {
		name     string
		declared string
		data     []byte
		inline   bool
		wantType string
		wantDisp string
	}{
		{"inline png renders", "image/png", png, true, "image/png", "inline"},
		{"svg never inline", "image/svg+xml", []byte("<svg/>"), true, "application/octet-stream", "attachment"},
		{"html neutralized on download", "text/html", html, false, "application/octet-stream", "attachment"},
		{"pdf downloads as-is", "application/pdf", []byte("%PDF-1.4"), false, "application/pdf", "attachment"},
		{"declared image but non-image bytes not inlined", "image/png", html, true, "image/png", "attachment"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ct, disp := resolveContentType(c.declared, c.data, c.inline)
			if ct != c.wantType || disp != c.wantDisp {
				t.Fatalf("got (%q, %q), want (%q, %q)", ct, disp, c.wantType, c.wantDisp)
			}
		})
	}
}

func TestNormalizeCID(t *testing.T) {
	cases := map[string]string{
		"<image001.png@01DD0333.B53E0920>": "image001.png@01DD0333.B53E0920",
		"cid:logo.png@host":                "logo.png@host",
		"  < spaced@id >  ":                "spaced@id",
		"image001.png%4001DD0333":          "image001.png@01DD0333",
		"plain@id":                         "plain@id",
	}
	for in, want := range cases {
		if got := normalizeCID(in); got != want {
			t.Errorf("normalizeCID(%q) = %q, want %q", in, got, want)
		}
	}
}

func buildRelatedMessage(t *testing.T, contentID, filename string, png []byte) []byte {
	t.Helper()
	b64 := base64.StdEncoding.EncodeToString(png)
	parts := []string{
		"From: a@b.com",
		"To: c@d.com",
		"Subject: test",
		"MIME-Version: 1.0",
		`Content-Type: multipart/related; boundary="BOUND"`,
		"",
		"--BOUND",
		`Content-Type: text/html; charset="utf-8"`,
		"",
		`<html><body><img src="cid:image001.png@01DD0333.B53E0920"></body></html>`,
		"--BOUND",
		"Content-Type: image/png; name=\"" + filename + "\"",
		"Content-Transfer-Encoding: base64",
		"Content-ID: " + contentID,
		`Content-Disposition: inline; filename="` + filename + `"`,
		"",
		b64,
		"--BOUND--",
		"",
	}
	return []byte(strings.Join(parts, "\r\n"))
}

func TestResolveInlineImage(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x01, 0x02, 0x03}
	raw := buildRelatedMessage(t, "<image001.png@01DD0333.B53E0920>", "image001.png", png)

	t.Run("matches by content-id with brackets", func(t *testing.T) {
		data, ct, fn, ok := resolveInlineImage(raw, "image001.png@01DD0333.B53E0920")
		if !ok {
			t.Fatal("expected to resolve inline image")
		}
		if string(data) != string(png) {
			t.Fatalf("decoded bytes mismatch: got %v", data)
		}
		if ct != "image/png" {
			t.Errorf("content type = %q, want image/png", ct)
		}
		if fn != "image001.png" {
			t.Errorf("filename = %q, want image001.png", fn)
		}
	})

	t.Run("matches case-insensitively", func(t *testing.T) {
		if _, _, _, ok := resolveInlineImage(raw, "IMAGE001.PNG@01dd0333.b53e0920"); !ok {
			t.Fatal("expected case-insensitive match")
		}
	})

	t.Run("falls back to filename", func(t *testing.T) {
		byFilename := buildRelatedMessage(t, "<unrelated-guid@outlook>", "image001.png", png)
		if _, _, _, ok := resolveInlineImage(byFilename, "image001.png"); !ok {
			t.Fatal("expected filename fallback match")
		}
	})

	t.Run("misses unknown cid", func(t *testing.T) {
		if _, _, _, ok := resolveInlineImage(raw, "nope@nowhere"); ok {
			t.Fatal("expected no match for unknown cid")
		}
	})
}

func TestFormatContentDisposition(t *testing.T) {
	cases := []struct {
		name     string
		disp     string
		filename string
		want     string
	}{
		{"ascii only", "attachment", "report.pdf", `attachment; filename="report.pdf"`},
		{"non-ascii adds rfc5987", "attachment", "café.png", `attachment; filename="caf_.png"; filename*=UTF-8''caf%C3%A9.png`},
		{"control chars stripped", "inline", "a\r\nb.png", `inline; filename="ab.png"`},
		{"quote escaped to fallback", "inline", "a\"b.png", `inline; filename="a_b.png"; filename*=UTF-8''a%22b.png`},
		{"empty falls back", "attachment", "", `attachment; filename="download"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := formatContentDisposition(c.disp, c.filename); got != c.want {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}
