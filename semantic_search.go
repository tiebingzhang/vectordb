package vectordb

import (
	"github.com/tiebingzhang/vectordb/rank"
	"github.com/tiebingzhang/vectordb/typings"
	"github.com/tiebingzhang/vectordb/vectors"
)

func SemanticSearch(query []string, corpus []string, results int, sorted bool) ([][]typings.SearchResult, error) {
	encodedCorpus, err := vectors.EncodeMulti(corpus)
	if err != nil {
		return [][]typings.SearchResult{}, err
	}
	encodedQuery, err := vectors.EncodeMulti(query)
	if err != nil {
		return [][]typings.SearchResult{}, err
	}

	// Semantic search
	searchResult := rank.Rank(encodedQuery, encodedCorpus, results, sorted)
	return searchResult, nil
}
