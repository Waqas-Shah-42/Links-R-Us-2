package index

import (
	"time"

	"github.com/google/uuid"
)

const (
	QueryTypeMatch QueryType = iota
	QueryTypePhrase
)

type Document struct {
	LinkID    uuid.UUID
	URL       string
	Title     string
	Content   string
	IndexedAt time.Time
	PageRank  float64
}
type QueryType uint8

// The type of query (in any order, exact) will depend on value of Type
type Query struct {
	Type       QueryType
	Expression string
	Offset     uint64
}

type Indexer interface {
	// Inserts document to the index or updates the index entry.
	Index(doc *Document) error

	FindByID(LinkID uuid.UUID) (*Document, error)
	Search(query Query) (Iterator, error)
	UpdateScore(LinkID uuid.UUID, score float64) error
}

type Iterator interface {
	Close() error

	// Loads the next document matching search query
	// Returns false if no more documents available
	Next() bool

	// Return the last error encountered by the Iterator
	Error() error

	// Return the current document in the iterator
	Document() *Document

	// Returns approximate number of results
	TotalCount() uint64
}
