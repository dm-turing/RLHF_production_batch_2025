
package main  
import (  
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// dbConn represents a single database connection with a deferred close.
type dbConn struct {
    *sql.DB
    connID int
}

// Close implements the io.Closer interface for dbConn.
func (c *dbConn) Close() error {
    fmt.Printf("Closing connection ID: %d\n", c.connID)
    return c.DB.Close()
}

// dbPoolManager is a connection pool manager for SQLite databases.
type dbPoolManager struct {
    mu       sync.Mutex
    dbConns  []*dbConn
    nextConnID int
    maxConn   int
}

// NewDBPoolManager creates a new dbPoolManager with the specified maximum connections.
func NewDBPoolManager(maxConn int) *dbPoolManager {
    return &dbPoolManager{
        dbConns:  make([]*dbConn, 0),
        maxConn:  maxConn,
        nextConnID: 1,
    }
}

// getConnection returns an open database connection from the pool or creates a new one if needed.
func (pm *dbPoolManager) getConnection() (*dbConn, error) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Return a reusable connection if available.
    if len(pm.dbConns) > 0 {
        conn := pm.dbConns[len(pm.dbConns)-1]
        pm.dbConns = pm.dbConns[:len(pm.dbConns)-1]
        fmt.Printf("Using connection ID: %d\n", conn.connID)
        return conn, nil
    }

    // Create a new connection if the pool is empty and the maximum connections is not reached.
    if len(pm.dbConns) < pm.maxConn {
        db, err := sql.Open("sqlite3", "./example.db")
        if err != nil {
            return nil, fmt.Errorf("error opening database: %w", err)
        }

        connID := pm.nextConnID
        pm.nextConnID++
        conn := &dbConn{db, connID}
        fmt.Printf("Creating new connection ID: %d\n", connID)
        return conn, nil
    }

    return nil, fmt.Errorf("maximum connections reached")
}

// ReleaseConnection releases the database connection back to the pool.
func (pm *dbPoolManager) ReleaseConnection(conn *dbConn) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    pm.dbConns = append(pm.dbConns, conn)
}

// PerformRandomOperations performs random database operations on a given connection.
func PerformRandomOperations(conn *dbConn) {
    rand.Seed(time.Now().UnixNano())
    // Perform some random operations on the connection (simulated using simple SQL queries).
    for i := 0; i < rand.Intn(5)+1; i++ {
        _, err := conn.Exec("INSERT INTO test_table (id, value) VALUES (?, ?)", rand.Int(), rand.Int())
        if err != nil {
            log.Printf("Error executing query: %v\n", err)
            break
        }
    }
}

func main() {  
    // Create a database connection pool manager with a maximum of 5 connections.
    pm := NewDBPoolManager(5)
    defer pm.shutdown()
    
    // Create random number of database instances and perform random operations on them.
    numDatabases := rand.Intn(10) + 1
    var wg sync.WaitGroup
    wg.Add(numDatabases)
    