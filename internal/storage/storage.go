package storage

import (
	"bytes"
	"fmt"
	"errors"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	// secondaryIndexPrefix is the prefix for secondary index keys.
	secondaryIndexPrefix = "_idx:"
)

// Storage is the interface for the storage engine.
type Storage interface {
	Set(key, value []byte) error
	SetWithIndex(key, value []byte, indexName string, indexValue []byte) error
	Get(key []byte) ([]byte, error)
	GetByIndex(indexName string, indexValue []byte) ([][]byte, error)
	Delete(key []byte) error
}

var ErrKeyNotFound = errors.New("key not found")

// BadgerStorage is the implementation of the storage engine using BadgerDB.
type BadgerStorage struct {
	db *badger.DB
}

// NewBadgerStorage creates a new BadgerStorage.
func NewBadgerStorage(path string) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerStorage{db: db}, nil
}

// Set sets a key-value pair.
func (s *BadgerStorage) Set(key, value []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	return err
}

// SetWithIndex sets a key-value pair and creates a secondary index entry.
// The secondary index key format will be: _idx:indexName:indexValue:primaryKey
func (s *BadgerStorage) SetWithIndex(key, value []byte, indexName string, indexValue []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		// Set the primary key-value pair
		if err := txn.Set(key, value); err != nil {
			return err
		}

		// Create the secondary index key
		indexKey := makeIndexKey(indexName, indexValue, key)
		// Store an empty value for the secondary index entry, as the primary key is already in the index key
		if err := txn.Set(indexKey, []byte{}); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get gets a value by key.
func (s *BadgerStorage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return value, err
}

// GetByIndex retrieves primary keys associated with a secondary index value.
func (s *BadgerStorage) GetByIndex(indexName string, indexValue []byte) ([][]byte, error) {
	var primaryKeys [][]byte
	prefix := makeIndexPrefix(indexName, indexValue)

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()

			// Extract the primary key from the index key
			pk, err := extractPrimaryKeyFromIndexKey(k, indexName, indexValue)
			if err != nil {
				return fmt.Errorf("failed to extract primary key from index key %s: %w", k, err)
			}
			primaryKeys = append(primaryKeys, pk)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return primaryKeys, nil
}

// Delete deletes a key-value pair.
func (s *BadgerStorage) Delete(key []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return err
}

// Close closes the storage engine.
func (s *BadgerStorage) Close() error {
	return s.db.Close()
}

// makeIndexKey constructs a secondary index key: _idx:indexName:indexValue:primaryKey
func makeIndexKey(indexName string, indexValue []byte, primaryKey []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString(secondaryIndexPrefix)
	buf.WriteString(indexName)
	buf.WriteString(":")
	buf.Write(indexValue)
	buf.WriteString(":")
	buf.Write(primaryKey)
	return buf.Bytes()
}

// makeIndexPrefix constructs a prefix for iterating over secondary index entries for a given indexName and indexValue.
func makeIndexPrefix(indexName string, indexValue []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString(secondaryIndexPrefix)
	buf.WriteString(indexName)
	buf.WriteString(":")
	buf.Write(indexValue)
	buf.WriteString(":") // Include the trailing colon to match only exact index values
	return buf.Bytes()
}

// extractPrimaryKeyFromIndexKey extracts the primary key from a secondary index key.
func extractPrimaryKeyFromIndexKey(indexKey []byte, indexName string, indexValue []byte) ([]byte, error) {
	expectedPrefix := makeIndexPrefix(indexName, indexValue)
	if !bytes.HasPrefix(indexKey, expectedPrefix) {
		return nil, fmt.Errorf("index key %s does not start with expected prefix %s", indexKey, expectedPrefix)
	}
	return indexKey[len(expectedPrefix):], nil
}
