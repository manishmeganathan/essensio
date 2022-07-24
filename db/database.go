package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

type Database struct {
	client *badger.DB
}

// Open opens a Badger client to the database at Dir()
func Open() (*Database, error) {
	// Setup Badger Options
	opts := badger.DefaultOptions(Dir())
	opts.Logger = nil

	// Open Badger Client
	client, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("db open fail: %w", err)
	}

	// Wrap client inside Database and return
	return &Database{client}, nil
}

// Close closes the Badger client to the database at Dir()
func (db *Database) Close() {
	if err := db.client.Close(); err != nil {
		panic(fmt.Errorf("db close fail: %w", err))
	}
}

func (db *Database) GetEntry(key []byte) (value []byte, err error) {
	// Define a view transaction on the database
	err = db.client.View(func(txn *badger.Txn) error {
		// Attempt to get the Item for the given key
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("db get on key '%x' fail: %w", key, err)
		}

		// Retrieve the value from the Item
		if err = item.Value(func(val []byte) error {
			value = val
			return nil

		}); err != nil {
			return fmt.Errorf("db value get on key '%x' fail: %w", key, err)
		}

		return nil
	})

	return
}

func (db *Database) SetEntry(key, value []byte) error {
	// Define an update transaction the database
	return db.client.Update(func(txn *badger.Txn) error {
		// Attempt to set the key-value pair to the database
		if err := txn.Set(key, value); err != nil {
			return fmt.Errorf("db set for key '%x' failed: %w", key, err)
		}

		return nil
	})
}
