package connections

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Aman123at/cdc-go/utils"
	_ "github.com/lib/pq"
)

var PgPool *connPool

type Connection struct {
	dbInstance *sql.DB
}

type connPool struct {
	mute           *sync.Mutex
	connections    []*Connection
	maxConnections int
	channel        chan interface{}
}

func NewConnectionPool(maxConns int) (*connPool, error) {
	var m = sync.Mutex{}
	pool := &connPool{
		mute:           &m,
		connections:    make([]*Connection, 0, maxConns),
		maxConnections: maxConns,
		channel:        make(chan interface{}, maxConns),
	}
	for i := 0; i < maxConns; i++ {
		pool.connections = append(pool.connections, &Connection{NewConn()})
		pool.channel <- nil
	}
	return pool, nil
}

func NewConn() *sql.DB {
	conn, connerr := sql.Open("postgres", utils.GetPostgresConnString())
	if connerr != nil {
		log.Println("Unable to connect postgres")
		log.Fatal(connerr.Error())
	}

	pingerr := conn.Ping()

	if pingerr != nil {
		log.Println("Unable to connect postgres")
		log.Fatal(pingerr.Error())
	}

	return conn
}

func (pool *connPool) Close() {
	close(pool.channel)
	for i := range pool.connections {
		pool.connections[i].dbInstance.Close()
	}
}

func (pool *connPool) Get() (*Connection, error) {
	<-pool.channel
	pool.mute.Lock()
	c := pool.connections[0]
	pool.connections = pool.connections[1:]
	pool.mute.Unlock()
	return c, nil
}

func (pool *connPool) Put(c *Connection) {
	pool.mute.Lock()
	pool.connections = append(pool.connections, c)
	pool.mute.Unlock()
	pool.channel <- nil
}
