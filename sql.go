package sqlxcache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type DB struct {
	*sqlx.DB
	cache *Cache
}

// default time to life for cached sql results
var ttl = 30 * time.Second

func NewDbx(db *sql.DB, driverName string) *DB {
	dbx := sqlx.NewDb(db, driverName)
	c := NewCache(nil)
	return &DB{DB: dbx, cache: c}
}

// generates a sha256 hashcode (as string) for the query and its params
func hashSQLQuery(query string, params ...any) string {
	// concat the query and params to create a unique string
	concat := query
	for _, param := range params {
		concat += fmt.Sprintf("%v", param)
	}
	h := sha256.New()
	h.Write([]byte(concat))
	// return the string representation of the hash bytes
	return hex.EncodeToString(h.Sum(nil))
}

// Select checks whether the query has a value in the cache
// if there is a value in the cache for the given combination of query and params - the value is returned
// otherwise the database is queried and the result is stored in the cahce
func (db *DB) Select(dest any, query string, args ...any) error {
	hash := hashSQLQuery(query, args)
	cachedValue := db.cache.Get(hash)
	if cachedValue == nil {
		err := db.DB.Select(dest, query, args...)
		if err != nil {
			return err
		}
		db.cache.Put(hash, dest, time.Now().Add(ttl))
		return nil
	}
	dest = cachedValue
	return nil
}
