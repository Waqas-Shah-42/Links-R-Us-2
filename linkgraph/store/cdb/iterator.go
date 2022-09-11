package cdb

import (
	"database/sql"

	"github.com/Waqas-Shah-42/Links-R-Us-2/linkgraph/graph"
	"golang.org/x/xerrors"
)

//implements graph.LinkIterator
type linkIterator struct {
	rows *sql.Rows
	lastErr	error
	latchedLink	*graph.Link
}

func (i *linkIterator) Next() bool {
	if i.lastErr != nil || !i.rows.Next() {
		return false
	}

	l := new(graph.Link)
	i.lastErr = i.rows.Scan(&l.ID,&l.URL, &l.RetrievedAt)
	if i.lastErr != nil {
		return false
	}
	l.RetrievedAt = l.RetrievedAt.UTC()

	i.latchedLink = l
	return true
}


func (i *linkIterator) Error() error {
	return i.lastErr
}

func (i *linkIterator) Close() error {
	err := i.rows.Close()
	if err != nil {
		return xerrors.Errorf("Link iterator: %w", err)
	}

	return nil
}

func (i *linkIterator) Link() *graph.Link {
	return i.latchedLink
}


//edgeIterator
type edgeIterator struct {
	rows 	*sql.Rows
	lastErr	error
	latchedEdge *graph.Edge
}
