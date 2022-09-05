package memory

import (
	"sync"

	"github.com/Waqas-Shah-42/Links-R-Us-2/textindexer/index"
	"github.com/blevesearch/bleve"
)

type InMemoryBleveIndexer struct {
	mu sync.RWMutex

	docs map[string]*index.Document

	idx bleve.Index

}
