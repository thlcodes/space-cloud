package pubsub

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Module deals with pub sub related activities
type Module struct {
	lock sync.Mutex

	// Redis client
	client *redis.Client

	// Internal variables
	projectID string
	mapping   map[string]*subscription
}

// New creates a new instance of the client
func New(projectID, conn string) (*Module, error) {
	// Set a default connection string if not provided
	if conn == "" {
		conn = "localhost:6379"
	}
	pw := ""
	var tc *tls.Config
	if url, err := url.Parse(conn); err == nil {
		if url.Host != "" {
			conn = url.Host
		}
		if pass, ok := url.User.Password(); ok {
			pw = pass
		}
		if url.Scheme == "rediss" {
			tc = &tls.Config{}
		}
	}
	c := redis.NewClient(&redis.Options{
		Addr:      conn,
		Password:  pw,
		DB:        0,
		TLSConfig: tc,
	})
	log.Printf("REDIS: host=%s pw=%s***%s tls=%t", conn, pw[:1], pw[len(pw)-1:], tc != nil)

	// Create a temporary context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := c.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Module{client: c, projectID: projectID, mapping: map[string]*subscription{}}, nil
}

// Close closes the redis client along with the active subscriptions on it
func (m *Module) Close() {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Close all active subscriptions first
	for _, sub := range m.mapping {
		_ = sub.pubsub.Close()
	}
	m.mapping = map[string]*subscription{}

	// Close the redis client
	_ = m.client.Close()
}
