package es

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Waqas-Shah-42/Links-R-Us-2/textindexer/index"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"golang.org/x/xerrors"
)

const (
	indexName = "textindexer"
	batchSize = 10
)

var esMappings = `
{
  "mappings" : {
    "properties": {
      "LinkID": {"type": "keyword"},
      "URL": {"type": "keyword"},
      "Content": {"type": "text"},
      "Title": {"type": "text"},
      "IndexedAt": {"type": "date"},
      "PageRank": {"type": "double"}
    }
  }
}`

type esSearchRes struct {
	Hits esSearchResHits `json: "hits"`
}

type esSearchResHits struct {
	Total   esTotal        `json: "total"`
	HitList []esHitWrapper `json: "hits"`
}

type esTotal struct {
	Count uint64 `json: "value"`
}

type esHitWrapper struct {
	DocSource esDoc `json: "_source"`
}

type esDoc struct {
	LinkID    string    `json:"LinkID"`
	URL       string    `json:"URL"`
	Title     string    `json:"Title"`
	Content   string    `json:"Content"`
	IndexedAt time.Time `json:"IndexedAt"`
	PageRank  float64   `json:"PageRank,omitempty"`
}

type esUpdateRes struct {
	Result string `json:"result"`
}

type esErrorRes struct {
	Error esError `json:"error"`
}

type esError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type ElasticSearchIndexer struct {
	es         *elasticsearch.Client
	refreshOpt func(*esapi.UpdateRequest)
}

func (e esError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Reason)
}

// Compile-time check to ensure ElasticSearchIndexer implements Indexer.
var _ index.Indexer = (*ElasticSearchIndexer)(nil)

func unmarshalResponse(res *esapi.Response, to interface{}) error {
	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		var errRes esErrorRes
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return err
		}

		return errRes.Error
	}

	return json.NewDecoder(res.Body).Decode(to)
}

func unmarshalError(res *esapi.Response) error {
	return unmarshalResponse(res, nil)
}

func ensureIndex(es *elasticsearch.Client) error {
	mappingsReader := strings.NewReader(esMappings)
	res, err := es.Indices.Create(indexName, es.Indices.Create.WithBody(mappingsReader))
	if err != nil {
		return xerrors.Errorf("cannot create ES index: %w", err)
	} else if res.IsError() {
		err := unmarshalError(res)
		if esErr, valid := err.(esError); valid && esErr.Type == "resource_already_exists_exception" {
			return nil
		}
		return xerrors.Errorf("cannot create ES index: %w", err)
	}

	return nil
}

func NewElasticSearchIndexer(esNodes []string, syncUpdates bool) (*ElasticSearchIndexer, error) {
	cfg := elasticsearch.Config{
		Addresses: esNodes,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if err = ensureIndex(es); err != nil {
		return nil, err
	}
}
