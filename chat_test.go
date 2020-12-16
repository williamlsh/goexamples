package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

func Test_chatServer(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		url, closeFn := setupTest(t)
		defer closeFn()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		cl, err := newClient(ctx, url)
		assertSuccess(t, err)
		defer cl.close()

		expMsg := randString(512)
		err = cl.publish(ctx, expMsg)
		assertSuccess(t, err)

		msg, err := cl.nextMessage()
		assertSuccess(t, err)

		if expMsg != msg {
			t.Fatalf("expected %v but got %v", expMsg, msg)
		}
	})

	t.Run("concurrency", func(t *testing.T) {
		t.Parallel()

		const (
			nmessages      = 128
			maxMessageSize = 128
			nclients       = 16
		)

		url, closeFn := setupTest(t)
		defer closeFn()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var (
			clients    []*client
			clientMsgs []map[string]struct{}
		)
		for i := 0; i < nclients; i++ {
			cl, err := newClient(ctx, url)
			assertSuccess(t, err)
			defer cl.close()

			clients = append(clients, cl)
			clientMsgs = append(clientMsgs, randMessages(nmessages, maxMessageSize))
		}

		allMessages := make(map[string]struct{})
		for _, msgs := range clientMsgs {
			for m := range msgs {
				allMessages[m] = struct{}{}
			}
		}

		var wg sync.WaitGroup
		for i, cl := range clients {
			i := i
			cl := cl

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := cl.publishMsgs(ctx, clientMsgs[i]); err != nil {
					t.Errorf("client %d failed to publish all messages: %v", i, err)
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := testAllMessagesReceived(cl, nclients*nmessages, allMessages); err != nil {
					t.Errorf("client %d failed to receive all messages: %v", i, err)
				}
			}()
		}

		wg.Wait()
	})
}

func setupTest(t *testing.T) (url string, closeFn func()) {
	cs := newChatServer()
	cs.logf = t.Logf

	cs.subscriberMessageBuffer = 4096
	cs.publishLimiter.SetLimit(rate.Inf)

	s := httptest.NewServer(cs)
	return s.URL, func() {
		s.Close()
	}
}

type client struct {
	url string
	c   *websocket.Conn
}

func newClient(ctx context.Context, url string) (*client, error) {
	c, _, err := websocket.Dial(ctx, url+"/subscribe", nil)
	if err != nil {
		return nil, err
	}

	return &client{url, c}, nil
}

func (cl *client) publish(ctx context.Context, msg string) (err error) {
	defer func() {
		if err != nil {
			cl.c.Close(websocket.StatusInternalError, "publish failed")
		}
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, cl.url+"/publish", strings.NewReader(msg))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("publish request failed: %v", resp.StatusCode)
	}

	return nil
}

func (cl *client) publishMsgs(ctx context.Context, msgs map[string]struct{}) error {
	for m := range msgs {
		if err := cl.publish(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (cl *client) nextMessage() (string, error) {
	typ, b, err := cl.c.Read(context.Background())
	if err != nil {
		return "", err
	}

	if typ != websocket.MessageText {
		cl.c.Close(websocket.StatusUnsupportedData, "expected text message")
		return "", fmt.Errorf("expected text message but got %v", typ)
	}

	return string(b), nil
}

func (cl *client) close() error {
	return cl.c.Close(websocket.StatusNormalClosure, "")
}

func assertSuccess(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func randString(n int) string {
	b := make([]byte, n)
	_, err := rand.Reader.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate rand bytes: %v", err))
	}

	s := strings.ToValidUTF8(string(b), "_")
	s = strings.ReplaceAll(s, "\x00", "_")
	if len(s) > n {
		return s[:n]
	}
	if len(s) < n {
		extra := n - len(s)
		return s + strings.Repeat("=", extra)
	}

	return s
}

func randMessages(n, maxMessageLength int) map[string]struct{} {
	msgs := make(map[string]struct{})
	for i := 0; i < n; i++ {
		m := randString(randInt(maxMessageLength))
		if _, ok := msgs[m]; ok {
			i--
			continue
		}
		msgs[m] = struct{}{}
	}
	return msgs
}

func randInt(max int) int {
	x, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(fmt.Sprintf("failed to get random int: %v", err))
	}
	return int(x.Int64())
}

func testAllMessagesReceived(cl *client, n int, msgs map[string]struct{}) error {
	msgs = cloneMessages(msgs)

	for i := 0; i < n; i++ {
		msg, err := cl.nextMessage()
		if err != nil {
			return err
		}
		delete(msgs, msg)
	}

	if len(msgs) != 0 {
		return fmt.Errorf("did not receive all expected messages: %q", msgs)
	}
	return nil
}

func cloneMessages(msgs map[string]struct{}) map[string]struct{} {
	msgs2 := make(map[string]struct{}, len(msgs))
	for m := range msgs {
		msgs2[m] = struct{}{}
	}
	return msgs2
}
