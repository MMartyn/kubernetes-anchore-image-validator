package anchore

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupEnv(url string) {
	os.Setenv("ANCHORE_ENGINE_URL", url)
	os.Setenv("ANCHORE_ENGINE_USERNAME", "admin")
	os.Setenv("ANCHORE_ENGINE_PASSWORD", "foobar")
}

func TestGetStatus(t *testing.T) {
	digest := "sha256:1d8f14b6d4e01369e1df18cfae17eb0894a39a21c28c6f8dbf6e2fe895b36522"
	tag := "latest"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[{"%s":{"docker.io/%s:latest": [{"status":"pass"}]}}]`, digest, tag)
	}))
	defer ts.Close()

	setupEnv(ts.URL)

	if !getStatus(digest, tag) {
		t.Log("Status was fail when it should have been pass")
		t.Fail()
	}

}

func TestGetDigest(t *testing.T) {
	imageRef := "viglesiasce/sample-app"
	digest := "sha256:1d8f14b6d4e01369e1df18cfae17eb0894a39a21c28c6f8dbf6e2fe895b36522"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[{"imageDigest": "%s"}]`, digest)
	}))
	defer ts.Close()

	setupEnv(ts.URL)

	digestResult, err := getImageDigest(imageRef)
	if err != nil {
		t.Logf("getImageDigest returned an error: %v", err)
		t.Fail()
	}
	if digestResult != digest {
		t.Logf("getImageDigest returned wrong digest: %s", digestResult)
		t.Fail()
	}
}
