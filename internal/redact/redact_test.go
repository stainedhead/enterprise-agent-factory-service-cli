package redact

import (
	"net/http"
	"strings"
	"testing"
)

// NFR observability: Authorization and cookies are redacted.
func TestRedactsHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer afsk_secret")
	h.Set("Cookie", "a=b")
	h.Add("Set-Cookie", "c=d")
	h.Set("X-Request-Id", "req_1")
	got := strings.Join(Headers(h), "\n")
	if strings.Contains(got, "afsk_secret") || strings.Contains(got, "a=b") || strings.Contains(got, "c=d") {
		t.Fatalf("leak: %s", got)
	}
	if !strings.Contains(got, "Authorization: [redacted]") || !strings.Contains(got, "X-Request-Id: req_1") {
		t.Errorf("headers = %s", got)
	}
}

// NFR observability: body fields of data-dictionary section 8, nested and in arrays.
func TestRedactsBodyFields(t *testing.T) {
	in := `{"name":"gh","value":"s3cret","nested":{"secretKey":"afsk_x","ok":1},"list":[{"token":"t"},{"password":"p"}],
	"refreshToken":"r","accessToken":"a","clientSecret":"c","idToken":"i","deviceCode":"d","n":1.50}`
	got := JSON([]byte(in))
	for _, leak := range []string{"s3cret", "afsk_x", `"t"`, `"p"`, `"r"`, `"a"`, `"c"`, `"i"`, `"d"`} {
		if strings.Contains(got, leak) {
			t.Errorf("leaked %s in %s", leak, got)
		}
	}
	if !strings.Contains(got, `"name":"gh"`) || !strings.Contains(got, `"ok":1`) || !strings.Contains(got, `1.50`) {
		t.Errorf("lost data: %s", got)
	}
	if strings.Count(got, Mask) != 9 {
		t.Errorf("masks = %d in %s", strings.Count(got, Mask), got)
	}
}

func TestNonJSONBodyIsNeverLogged(t *testing.T) {
	if got := JSON([]byte("password=hunter2")); strings.Contains(got, "hunter2") || !strings.Contains(got, "not JSON") {
		t.Errorf("got %q", got)
	}
	if JSON(nil) != "" || JSON([]byte("  ")) != "" {
		t.Error("empty body")
	}
	if len(Fields()) != 9 {
		t.Errorf("fields = %v", Fields())
	}
}

// CLI-US-33.1: the edge credential header and the edge cookie are masked in every debug line.
func TestRedactsEdgeHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("X-Afs-Authorization", "Bearer afs-token-sentinel")
	h.Set("Authorization", "Bearer paperclip-key-sentinel")
	h.Set("Proxy-Authorization", "Basic proxy-sentinel")
	h.Set("Cookie", "__Host-afs-collab=cookie-sentinel")
	h.Add("Set-Cookie", "__Host-afs-collab=cookie-sentinel; Secure")
	got := strings.Join(Headers(h), "\n")
	for _, leak := range []string{"afs-token-sentinel", "paperclip-key-sentinel", "proxy-sentinel", "cookie-sentinel"} {
		if strings.Contains(got, leak) {
			t.Errorf("leaked %s in %s", leak, got)
		}
	}
	if strings.Count(got, Mask) != 5 || HeaderValue("x-afs-authorization", "v") != Mask {
		t.Errorf("headers = %s", got)
	}
}

// CLI-US-31.2: an ARN reference is logged, a clientSecret and any ...Secret value is not.
func TestCollabClientSecretArnNotRedactedButSecretIs(t *testing.T) {
	in := `{"collabClientId":"cid","collabClientSecretArn":"arn:aws:secretsmanager:eu-west-1:000000000000:secret:x","clientSecret":"s1","otherSecret":"s2","refSecret":"arn:aws:secretsmanager:eu-west-1:000000000000:secret:y"}`
	got := JSON([]byte(in))
	for _, leak := range []string{`"s1"`, `"s2"`} {
		if strings.Contains(got, leak) {
			t.Errorf("leaked %s in %s", leak, got)
		}
	}
	for _, keep := range []string{"arn:aws:secretsmanager:eu-west-1:000000000000:secret:x", "arn:aws:secretsmanager:eu-west-1:000000000000:secret:y", `"collabClientId":"cid"`} {
		if !strings.Contains(got, keep) {
			t.Errorf("lost %s in %s", keep, got)
		}
	}
}
