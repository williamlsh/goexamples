package x // import github.com/williamlsh/x

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiURL = "https://127.0.0.1:8080/api"
)

func BcjClient(content []string) ([]bool, error) {
	client, err := HTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	query := Query{
		Content: content,
	}
	b, err := json.Marshal(&query)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("http post failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode result: %v", err)
	}

	return result.Content, nil
}

func HTTPClient() (*http.Client, error) {
	ca, err := os.Open(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not open ca file: %w", err)
	}
	defer ca.Close()

	b, err := ioutil.ReadAll(ca)
	if err != nil {
		return nil, fmt.Errorf("could not read ca file: %w", err)
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(b)
	clientTLSConfig := &tls.Config{RootCAs: certpool}

	client := http.DefaultClient
	client.Transport = &http.Transport{
		TLSClientConfig: clientTLSConfig,
	}

	return client, nil
}
