package cdb

import (
	"database/sql"

	"github.com/Waqas-Shah-42/Links-R-Us-2/linkgraph/graph"
	"golang.org/x/xerrors"
)

var (
	upsertLinkQuery = `
	INSERT INTO links (url, retrieved_at) VALUES ($1, $2)
ON CONFLICT (url) DO UPDATE SET retrieved_at=GREATEST(links.retrieved_at, $2)
RETURNING id, retrieved_at
	`
	findLinkQuery         = "SELECT url, retrieved_at FROM links WHERE id=$1"
	linksInPartitionQuery = "SELECT id, url, retrieved_at FROM links WHERE id >= $1 AND id < $2 AND retrieved_at < $3"

	upsertEdgeQuery = `
INSERT INTO edges (src, dst, updated_at) VALUES ($1, $2, NOW())
ON CONFLICT (src,dst) DO UPDATE SET updated_at=NOW()
RETURNING id, updated_at
`
	edgesInPartitionQuery = "SELECT id, src, dst, updated_at FROM edges WHERE src >= $1 AND src < $2 AND updated_at < $3"
	removeStaleEdgesQuery = "DELETE FROM edges WHERE src=$1 AND updated_at < $2"

	// Compile-time check for ensuring CockroachDbGraph implements Graph.
	_ graph.Graph = (*CockroachDBGraph)(nil)
)

// Stores the connection to the db
type CockroachDBGraph struct {
	db *sql.DB
}

// Creates the connection to the database
func NewCockroachDbGraph(dsn string) (*CockroachDBGraph, error){
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}
	return &CockroachDBGraph{db:db}, nil
}

// Terminates the database connection
func (c *CockroachDBGraph) Close() error {
	return c.db.Close()
}

// Creates or Updates link
func (c *CockroachDBGraph) UpsertLink(link *graph.Link) error {
	row := c.db.QueryRow(upsertLinkQuery, link.URL, link.RetrievedAt.UTC())
	if err := row.Scan(&link.ID, &link.RetrievedAt); err != nil {
		return xerrors.Errorf("upsert link:%w", err)
	}

	link.RetrievedAt = link.RetrievedAt.UTC()
	return nil
}


