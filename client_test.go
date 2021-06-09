package x

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBcjClient(t *testing.T) {
	go Server()
	defer removeCerts()

	waitForCerts()

	t.Run("empty query", func(t *testing.T) {
		content := []string{}
		result, err := BcjClient(content)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, []bool{}) {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
	t.Run("Two times query one hit", func(t *testing.T) {
		content1 := []string{"a", "b", "c"}
		result, err := BcjClient(content1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, []bool{false, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}

		content2 := []string{"a", "c", "d"}
		result2, err := BcjClient(content2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result2, []bool{true, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}

		if !reflect.DeepEqual(cache.list, []Set{
			{"a": struct{}{}},
			{
				"b": struct{}{},
				"c": struct{}{},
			},
			{
				"c": struct{}{},
				"d": struct{}{},
			},
		}) {
			t.Fatalf("Unexpected cache result: %+v\n", cache.list)
		}
	})
	t.Run("shorter query content all hits", func(t *testing.T) {
		content := []string{"a", "b"}
		result, err := BcjClient(content)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, []bool{true, true}) {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
	t.Run("longger query content no hit", func(t *testing.T) {
		content := []string{"e", "f", "g", "h", "i"}
		result, err := BcjClient(content)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, []bool{false, false, false, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}
		fmt.Printf("cache: %+v\n", cache.list)
	})
}

func TestHTTPClient(t *testing.T) {
	if err := GenerateCerts(); err != nil {
		t.Fatal(err)
	}
	defer removeCerts()

	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	serverTLSConfig := &tls.Config{Certificates: []tls.Certificate{cer}}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))
	server.TLS = serverTLSConfig

	server.StartTLS()
	defer server.Close()

	client, err := HTTPClient()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "success!") {
		t.Fatalf("response body is unexpected: %s", string(b))
	}
}

func waitForCerts() {
	// Wait until certs are prepared.
	for {
		n := 0
		for _, file := range []string{
			caFile,
			certFile,
			keyFile,
		} {
			if _, err := os.Stat(file); !os.IsNotExist(err) {
				n++
			}
		}
		if n == 3 {
			break
		}
	}
}
