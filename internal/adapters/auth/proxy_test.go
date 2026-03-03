package auth

import (
	"testing"
)

func TestExtractAuthCookies_HasAuth(t *testing.T) {
	headers := []string{
		".ASPXFORMSAUTHUACLOUD=abc123; Path=/; HttpOnly",
		"ASP.NET_SessionId=sess456; Path=/; HttpOnly",
		"_ga=GA1.1.1234; Path=/",
	}
	got := extractAuthCookies(headers)
	if got == "" {
		t.Fatal("expected auth cookies, got empty string")
	}
	if !containsSubstr(got, ".ASPXFORMSAUTHUACLOUD=abc123") {
		t.Errorf("missing .ASPXFORMSAUTHUACLOUD in %q", got)
	}
	if !containsSubstr(got, "ASP.NET_SessionId=sess456") {
		t.Errorf("missing ASP.NET_SessionId in %q", got)
	}
}

func TestExtractAuthCookies_NoAuth(t *testing.T) {
	headers := []string{
		"_ga=GA1.1.1234; Path=/",
		"Politica=Aceptada; Path=/",
	}
	got := extractAuthCookies(headers)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestExtractAuthCookies_Empty(t *testing.T) {
	got := extractAuthCookies(nil)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRewriteURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		port int
		want string
	}{
		{
			name: "https cvnet",
			url:  "https://cvnet.cpd.ua.es/uaCloud/home",
			port: 18923,
			want: "http://localhost:18923/uaCloud/home",
		},
		{
			name: "http cvnet",
			url:  "http://cvnet.cpd.ua.es/uaCloud",
			port: 9999,
			want: "http://localhost:9999/uaCloud",
		},
		{
			name: "cas url rewritten",
			url:  "https://autentica.cpd.ua.es/cas/login?service=x",
			port: 18923,
			want: "http://localhost:18923/cas/login?service=x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteURL(tt.url, tt.port)
			if got != tt.want {
				t.Errorf("rewriteURL(%q, %d) = %q, want %q", tt.url, tt.port, got, tt.want)
			}
		})
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && contains(s, sub))
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
