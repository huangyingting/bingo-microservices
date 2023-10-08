package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"bingo/app/bs/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/go-kratos/kratos/v2/log"
)

type ElasticSearch struct {
	es        *elasticsearch.TypedClient
	indexName string
	h         *log.Helper
}

func NewElasticSearch(c *conf.Search, h *log.Helper) (*ElasticSearch, error) {

	cfg := elasticsearch.Config{
		Addresses:             c.Addrs,
		Username:              c.Username,
		Password:              c.Password,
		DiscoverNodesOnStart:  true,
		DiscoverNodesInterval: time.Duration(5 * time.Second),
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		h.Errorf("connect to elasticsearch error: %v", err)
		return nil, err
	}

	// Dynamically create the index with the specified number of shards/replicas
	// Settings
	settings := types.NewIndexSettings()
	settingsAnalysis := types.NewIndexSettingsAnalysis()
	settingsAnalysis.Normalizer = map[string]types.Normalizer{"lowercase": types.CustomNormalizer{Filter: []string{"lowercase"}}}
	settings.Analysis = settingsAnalysis
	settings.NumberOfShards = strconv.FormatInt(int64(c.NumberOfShards), 10)
	settings.NumberOfReplicas = strconv.FormatInt(int64(c.NumberOfReplicas), 10)

	// Mappings
	mappings := types.NewTypeMapping()
	mappings.Properties["alias"] = types.NewKeywordProperty()
	mappings.Properties["oid"] = types.NewKeywordProperty()
	tags := types.NewKeywordProperty()
	normalizer := "lowercase"
	tags.Normalizer = &normalizer
	tags.Fields["suggest"] = types.NewCompletionProperty()
	mappings.Properties["tags"] = tags

	// Create index if not created before
	ctx := context.Background()
	exists, err := es.Indices.Exists(c.IndexName).Do(ctx)
	if err != nil {
		h.Errorf("check index exist error: %v", err)
		return nil, err
	}
	if !exists {
		_, err = es.Indices.Create(c.IndexName).Settings(settings).Mappings(mappings).Do(ctx)
		if err != nil {
			h.Errorf("create index error: %v", err)
			return nil, err
		}

	}
	return &ElasticSearch{
		es:        es,
		indexName: c.IndexName,
		h:         h,
	}, nil
}

func (es *ElasticSearch) Index(alias Alias) error {
	es.h.Debugf("index alias: %v", alias)
	_, err := es.es.Index(es.indexName).Id(alias.Alias).Document(alias).Refresh(refresh.True).Do(context.Background())
	if err != nil {
		es.h.Errorf("index error: %v", alias)
	}
	return err
}

func (es *ElasticSearch) SearchOr(oid string, tags []string) ([]Alias, error) {
	es.h.Debugf("search or oid: %s tags: %v", oid, tags)
	var queries []types.Query
	tagsFieldValue := make([]types.FieldValue, len(tags))
	for i, s := range tags {
		tagsFieldValue[i] = s
	}

	queries = append(queries, types.Query{Terms: &types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{"tags": tagsFieldValue}}})

	queries = append(queries, types.Query{Term: map[string]types.TermQuery{"oid": {Value: oid}}})

	b := types.NewBoolQuery()
	b.Filter = queries

	resp, err := es.es.Search().Index(es.indexName).Query(&types.Query{
		Bool: b,
	}).Do(context.Background())

	if err != nil {
		es.h.Errorf("search or error: %v", err)
		return nil, err
	}

	var aliases []Alias
	for _, v := range resp.Hits.Hits {
		var alias Alias
		if err = json.Unmarshal(v.Source_, &alias); err != nil {
			es.h.Errorf("search or unmarshal error: %v", err)
			return nil, err
		}
		aliases = append(aliases, alias)
	}
	return aliases, nil
}

func (es *ElasticSearch) SearchAnd(oid string, tags []string) ([]Alias, error) {
	es.h.Debugf("search and oid: %s tags: %v", oid, tags)

	var queries []types.Query

	for _, v := range tags {
		queries = append(queries, types.Query{Term: map[string]types.TermQuery{"tags": {Value: v}}})
	}

	queries = append(queries, types.Query{Term: map[string]types.TermQuery{"oid": {Value: oid}}})

	b := types.NewBoolQuery()
	b.Filter = queries

	resp, err := es.es.Search().Index(es.indexName).Query(&types.Query{
		Bool: b,
	}).Do(context.Background())

	if err != nil {
		es.h.Errorf("search and error: %v", err)
		return nil, err
	}

	var aliases []Alias
	for _, v := range resp.Hits.Hits {
		var alias Alias
		if err = json.Unmarshal(v.Source_, &alias); err != nil {
			es.h.Errorf("search and unmarshal error: %v", err)
			return nil, err
		}
		aliases = append(aliases, alias)
	}
	return aliases, nil
}

func (es *ElasticSearch) Suggest(text string) ([]string, error) {
	es.h.Debugf("suggest: %s", text)
	prefix := "se"
	size := 5
	skipDuplicates := true
	filedSuggester := types.NewFieldSuggester()
	filedSuggester.Prefix = &prefix
	filedSuggester.Completion = types.NewCompletionSuggester()
	filedSuggester.Completion.Field = "tags.suggest"
	filedSuggester.Completion.Size = &size
	filedSuggester.Completion.SkipDuplicates = &skipDuplicates
	suggester := types.NewSuggester()
	suggester.Suggesters["tags-suggest"] = *filedSuggester

	result, err := es.es.Search().Index(es.indexName).Suggest(suggester).Source_(false).
		Do(context.Background())

	if err != nil {
		es.h.Errorf("suggest error: %v", err)
		return nil, err
	}

	results := make([]string, 0)
	for _, tags_suggest := range result.Suggest["tags-suggest"] {
		completionSuggest := tags_suggest.(*types.CompletionSuggest)
		for _, option := range completionSuggest.Options {
			results = append(results, option.Text)
		}
	}

	es.h.Debugf("suggest results: %v", results)
	return results, nil
}

func (es *ElasticSearch) Delete(alias string, oid string) error {
	es.h.Debugf("delete alias: %s oid: %s", alias, oid)
	var queries []types.Query

	queries = append(queries, types.Query{Term: map[string]types.TermQuery{"alias": {Value: alias}}})

	queries = append(queries, types.Query{Term: map[string]types.TermQuery{"oid": {Value: oid}}})

	b := types.NewBoolQuery()
	b.Filter = queries

	resp, err := es.es.Search().Index(es.indexName).Query(&types.Query{
		Bool: b,
	}).Source_(false).Do(context.Background())

	if err != nil {
		es.h.Errorf("search error: %v", err)
		return err
	}

	if len(resp.Hits.Hits) > 0 {
		_, err := es.es.Delete(es.indexName, alias).Refresh(refresh.True).Do(context.Background())
		if err != nil {
			es.h.Errorf("delete error: %v", err)
			return err
		}
	} else {
		es.h.Errorf("delete - alias: %s with oid: %s not found", alias, oid)
		return fmt.Errorf("delete - alias: %s with oid: %s not found", alias, oid)
	}
	return nil
}
