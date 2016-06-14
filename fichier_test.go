package fichier

import (
	"net/url"
	"testing"
)

func TestGetUploadHost(t *testing.T) {
	h, err := GetUploadHost()
	if err != nil {
		t.Fatal(err)
	}

	if h == "" {
		t.Fatal("should not return empty string")
	}

	if _, err := url.Parse(h); err != nil {
		t.Fatal("should not return a valid url:", err)
	}

	t.Log("host: ", h)
}
