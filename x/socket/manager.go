package socket

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/totegamma/concurrent/x/stream"
)

// ソケット管理のルール
// 1. Subscribeのリクエストのうち外部のリクエストは増える方向にしか更新をかけない
// 2. 外から受信したイベントのうち、すでにキャッシュのキーがあるものに対してappendをかける
// 3. chunk更新タイミングで、外部リクエストの棚卸しを行う
// 4. その際、継続している外部リクエストのキャッシュを新しく空で作っておく(なければ)

var ctx = context.Background()

type subscriptionManager struct {
	mc *memcache.Client
	rdb         *redis.Client
	clientSubs map[*websocket.Conn][]string
	remoteSubs map[string][]string

	remoteConns map[string]*websocket.Conn
}

func NewSubscriptionManager(mc *memcache.Client, rdb *redis.Client) *subscriptionManager {
	manager := &subscriptionManager{
		mc: mc,
		rdb: rdb,
		clientSubs: make(map[*websocket.Conn][]string),
		remoteSubs: make(map[string][]string),
		remoteConns: make(map[string]*websocket.Conn),
	}
	go manager.chunkUpdaterRoutine()
	return manager
}


func NewSubscriptionManagerForTest(mc *memcache.Client, rdb *redis.Client) *subscriptionManager {
	manager := &subscriptionManager{
		mc: mc,
		rdb: rdb,
		clientSubs: make(map[*websocket.Conn][]string),
		remoteSubs: make(map[string][]string),
		remoteConns: make(map[string]*websocket.Conn),
	}
	return manager
}

// Subscribe subscribes a client to a stream
func (m *subscriptionManager) Subscribe(conn *websocket.Conn, streams []string) {
	m.clientSubs[conn] = streams
	m.createInsufficientSubs() // TODO: this should be done in a goroutine
}

// Unsubscribe unsubscribes a client from a stream
func (m *subscriptionManager) Unsubscribe(conn *websocket.Conn) {
	delete(m.clientSubs, conn)
}

func (m *subscriptionManager) createInsufficientSubs() {
	currentSubs := make(map[string][]string)
	for _, streams := range m.clientSubs {
		for _, stream := range streams {
			split := strings.Split(stream, "@")
			if len(split) != 2 {
				continue
			}
			domain := split[1]
			if _, ok := currentSubs[domain]; !ok {
				currentSubs[domain] = append(currentSubs[domain], stream)
			}
		}
	}

	// on this func, update only if there is a new subscription
	for domain, streams := range currentSubs {
		if _, ok := m.remoteSubs[domain]; !ok {
			m.remoteSubs[domain] = streams
		} else {
			for _, stream := range streams {
				if !slices.Contains(m.remoteSubs[domain], stream) {
					m.remoteSubs[domain] = append(m.remoteSubs[domain], stream)
				}
			}
		}
	}

	// FIXME: should be call if the subsucription is updated
	for domain, streams := range m.remoteSubs {
		if _, ok := m.remoteConns[domain]; !ok {
			m.RemoteSubRoutine(domain, streams)
		}
	}
}

// DeleteExcessiveSubs deletes subscriptions that are not needed anymore
func (m *subscriptionManager) deleteExcessiveSubs() {
	var currentSubs []string
	for _, streams := range m.clientSubs {
		for _, stream := range streams {
			if !slices.Contains(currentSubs, stream) {
				currentSubs = append(currentSubs, stream)
			}
		}
	}

	var closeList []string

	for domain, streams := range m.remoteSubs {
		for _, stream := range streams {
			var newSubs []string
			for _, currentSub := range currentSubs {
				if currentSub == stream {
					newSubs = append(newSubs, currentSub)
				}
			}
			m.remoteSubs[domain] = newSubs

			if len(m.remoteSubs[domain]) == 0 {
				closeList = append(closeList, domain)
			}
		}
	}

	for _, domain := range closeList {

		// close connection
		if conn, ok := m.remoteConns[domain]; ok {
			conn.Close()
		}

		delete(m.remoteSubs, domain)
		delete(m.remoteConns, domain)
	}
}

// RemoteSubRoutine subscribes to a remote server
func (m *subscriptionManager) RemoteSubRoutine(domain string, streams []string) {
	if _, ok := m.remoteConns[domain]; !ok {
		// new server, create new connection
		u := url.URL{Scheme: "wss", Host: domain, Path: "/api/v1/socket"}
		dialer := websocket.DefaultDialer
		dialer.HandshakeTimeout = 10 * time.Second

		// TODO: add config for TLS
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		c, _, err := dialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("fail to dial: %v", err)
		}

		m.remoteConns[domain] = c

		// launch a new goroutine for handling incoming messages
		go func(c *websocket.Conn) {
			defer c.Close()
			for {
				// check if the connection is still alive
				if c == nil {
					log.Printf("connection is nil (domain: %s)", domain)
					return
				}
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Printf("fail to read message: %v", err)
					return
				}

				var event stream.Event
				err = json.Unmarshal(message, &event)
				if err != nil {
					log.Printf("fail to Unmarshall redis message: %v", err)
				}

				// publish message to Redis
				err = m.rdb.Publish(ctx, event.Stream, string(message)).Err()
				if err != nil {
					log.Printf("fail to publish message to Redis: %v", err)
				}

				// update cache
				json, err := json.Marshal(event.Item)
				if err != nil {
					log.Printf("fail to Marshall item: %v", err)
				}
				json = append(json, ',')

				// update cache
				cacheKey := "stream:body:all:" + event.Item.StreamID + ":" + stream.Time2Chunk(event.Item.CDate)
				m.mc.Append(&memcache.Item{Key: cacheKey, Value: json})
			}
		}(c)
	}
	request := channelRequest{
		Channels: streams,
	}
	err := websocket.WriteJSON(m.remoteConns[domain], request)
	if err != nil {
		log.Printf("fail to send subscribe request to remote server %v: %v", domain, err)
		delete(m.remoteConns, domain)
	}
}

func (m *subscriptionManager) updateChunks(newchunk string) {
	// update cache
	for _, streams := range m.remoteSubs {
		for _, stream := range streams {
			m.mc.Add(&memcache.Item{Key: "stream:body:all:" + stream + ":" + newchunk, Value: []byte("")})
		}
	}
}

// ChunkUpdaterRoutine
func (m *subscriptionManager) chunkUpdaterRoutine() {
	currentChunk := stream.Time2Chunk(time.Now())
	for {
		// 次の実行時刻を計算
		nextRun := time.Now().Truncate(time.Hour).Add(time.Minute * 10)
		if time.Now().After(nextRun) {
			// 現在時刻がnextRunを過ぎている場合、次の10分単位の時刻を計算
			elapsed := time.Now().Sub(nextRun)
			nextRun = nextRun.Add(time.Minute * 10 * ((elapsed / (time.Minute * 10)) + 1))
		}

		// 次の実行時刻まで待機
		time.Sleep(time.Until(nextRun))

		// まだだったら待ちなおす
		newChunk := stream.Time2Chunk(time.Now())
		if newChunk == currentChunk {
			continue
		}

		m.deleteExcessiveSubs()
		m.updateChunks(newChunk)

		currentChunk = newChunk
	}
}

type channelRequest struct {
	Channels []string `json:"channels"`
}

