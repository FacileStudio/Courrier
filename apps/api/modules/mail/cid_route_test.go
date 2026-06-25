package mail

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	// base64 with MIME-style CRLF line breaks, as servers actually return it
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
