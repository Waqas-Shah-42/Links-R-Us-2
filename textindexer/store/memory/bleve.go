package memory

import (
	"sync"
	"time"

	"github.com/Waqas-Shah-42/Links-R-Us-2/textindexer/index"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

// Size of results cached by the iterator
const batchSize = 10

// Compile-time check to ensure InMemoryBleveIndexer implements Indexer.
var _ index.Indexer = (*InMemoryBleveIndexer)(nil)

type InMemoryBleveIndexer struct {
	mu sync.RWMutex

	//The Keys will be document link IDs while the values will be immutable copies of the documents
	//that are processed by the indexer
	docs map[string]*index.Document

	idx bleve.Index
}

type bleveDoc struct {
	Title    string
	Content  string
	PageRank float64
}

func copyDoc(d *index.Document) *index.Document {
	dcopy := new(index.Document)
	*dcopy = *d
	return dcopy
}

func makeBleveDoc(d *index.Document) bleveDoc {
	return bleveDoc{
		Title:    d.Title,
		Content:  d.Content,
		PageRank: d.PageRank,
	}
}

func (i *InMemoryBleveIndexer) Index(doc *index.Document) error {
	if doc.LinkID == uuid.Nil {
		return xerrors.Errorf("index: %w", index.ErrMissingLinkID)
	}
	doc.IndexedAt = time.Now()
	dcopy := copyDoc(doc)
	key := dcopy.LinkID.String()
	i.mu.Lock()
	if orig, exists := i.docs[key]; exists {
		dcopy.PageRank = orig.PageRank
	}
	if err := i.idx.Index(key, makeBleveDoc(dcopy)); err != nil {
		return xerrors.Errorf("idex: %w", err)
	}

	i.docs[key] = dcopy
	i.mu.Unlock()
	return nil
}

func (i *InMemoryBleveIndexer) FindByID(linkID uuid.UUID) (*index.Document, error) {
	//Method created for findByID() so that document lookup code can be reused when iterating through search results
	return i.findByID(linkID.String())

}
func (i *InMemoryBleveIndexer) findByID(linkID string) (*index.Document, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if d, found := i.docs[linkID]; found {
		return copyDoc(d), nil

	}
	return nil, xerrors.Errorf("find by ID: %w", index.ErrNotFound)
}

func (i *InMemoryBleveIndexer) UpdateScore(linkedID uuid.UUID, score float64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	key := linkedID.String()
	doc, found := i.docs[key]
	if !found {
		doc = &index.Document{LinkID: linkedID}
		i.docs[key] = doc
	}

	doc.PageRank = score
	if err := i.idx.Index(key, makeBleveDoc(doc)); err != nil {
		return xerrors.Errorf("update score: %w", err)
	}
	return nil

}


func (i *InMemoryBleveIndexer) Search(q index.Query) (index.Iterator, error) {
	var bq query.Query
	switch q.Type {
	case index.QueryTypeMatch:
		bq = bleve.NewMatchPhraseQuery(q.Expression)
	default:
		bq = bleve.NewMatchQuery(q.Expression)
	}

	searchReq := bleve.NewSearchRequest(bq)
	searchReq.SortBy([]string{"-PageRank", "-_score"})
	searchReq.Size = batchSize
	searchReq.From = int(q.Offset)

	rs, err := i.idx.Search(searchReq)
	if err != nil {
		return nil, xerrors.Errorf("search : %w", err)
	}
	return &bleveIterator{idx: i, searchReq: searchReq, rs: rs, cumIdx:q.Offset}, nil
}
