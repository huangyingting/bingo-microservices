package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"bingo/app/bs/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
)

const (
	MAPPING = `
	{
		"settings":{
			"number_of_shards": {{.NumberOfShards}},
			"number_of_replicas": {{.NumberOfReplicas}},
			"analysis": {
				"normalizer": {
					"lowercase": {
						"type": "custom",
						"filter": [
							"lowercase"
						]
					}
				}
			}	
		},
		"mappings":{
			"_source": {
				"enabled": true
			},	
			"properties":{
				"alias": {
					"type": "keyword"
				},
				"oid": {
					"type": "keyword"
				},
				"tags":{
					"type":"keyword",
					"normalizer": "lowercase",
					"fields": {
						"suggest": {
							"type": "completion"
						}
					}
				}
			}
		}
	}
	`
)

type ElasticSearch struct {
	es        *elastic.Client
	indexName string
	h         *log.Helper
}

func NewElasticSearch(c *conf.Search, h *log.Helper) (*ElasticSearch, error) {
	es, err := elastic.NewClient(elastic.SetURL(c.Addr...), elastic.SetSniff(c.Sniff))
	if err != nil {
		h.Errorf("connect to elasticsearch error: %v", err)
		return nil, err
	}

	// Dynamically create the index with the specified number of shards/replicas
	tmpl, err := template.New("T").Parse(MAPPING)
	if err != nil {
		h.Errorf("parse mapping template error: %v", err)
		return nil, err
	}
	var body bytes.Buffer
	err = tmpl.ExecuteTemplate(&body, "T", struct {
		NumberOfShards   uint32
		NumberOfReplicas uint32
	}{
		NumberOfShards:   c.NumberOfShards,
		NumberOfReplicas: c.NumberOfReplicas,
	})
	if err != nil {
		h.Errorf("execute template error: %v", err)
		return nil, err
	}

	// Create index if not created before
	ctx := context.Background()
	exists, err := es.IndexExists(c.IndexName).Do(ctx)
	if err != nil {
		h.Errorf("check index exist error: %v", err)
		return nil, err
	}
	if !exists {
		_, err = es.CreateIndex(c.IndexName).BodyString(body.String()).Do(ctx)
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
	_, err := es.es.Index().
		Index(es.indexName).Id(alias.Alias).
		BodyJson(alias).
		Refresh("true").
		Do(context.TODO())
	if err != nil {
		es.h.Errorf("index error: %v", alias)
	}
	return err
}

func (es *ElasticSearch) SearchOr(oid string, tags []string) ([]Alias, error) {
	es.h.Debugf("search or oid: %s tags: %v", oid, tags)
	tagsInterface := make([]interface{}, len(tags))
	for i, s := range tags {
		tagsInterface[i] = s
	}
	resp, err := es.es.Search().
		Index(es.indexName).Query(
		elastic.NewBoolQuery().
			Filter(elastic.NewTermsQuery("tags", tagsInterface...),
				elastic.NewTermsQuery("oid", oid))).
		Do(context.TODO())
	if err != nil {
		es.h.Errorf("search or error: %v", err)
		return nil, err
	}
	var aliases []Alias
	for _, v := range resp.Hits.Hits {
		var alias Alias
		if err = json.Unmarshal(v.Source, &alias); err != nil {
			es.h.Errorf("search or unmarshal error: %v", err)
			return nil, err
		}
		aliases = append(aliases, alias)
	}
	return aliases, nil
}

func (es *ElasticSearch) SearchAnd(oid string, tags []string) ([]Alias, error) {
	es.h.Debugf("search and oid: %s tags: %v", oid, tags)
	var queries []elastic.Query
	for _, v := range tags {
		t := elastic.NewTermQuery("tags", v)
		queries = append(queries, t)
	}
	queries = append(queries, elastic.NewTermsQuery("oid", oid))
	resp, err := es.es.Search().
		Index(es.indexName).Query(
		elastic.NewBoolQuery().
			Filter(queries...)).
		Do(context.TODO())
	if err != nil {
		es.h.Errorf("search and error: %v", err)
		return nil, err
	}
	var aliases []Alias
	for _, v := range resp.Hits.Hits {
		var alias Alias
		if err = json.Unmarshal(v.Source, &alias); err != nil {
			es.h.Errorf("search and unmarshal error: %v", err)
			return nil, err
		}
		aliases = append(aliases, alias)
	}
	return aliases, nil
}

func (es *ElasticSearch) Suggest(text string) ([]string, error) {
	es.h.Debug("suggest: %s", text)
	suggester := elastic.NewCompletionSuggester("tags-suggest").
		Text(text).Field("tags.suggest").SkipDuplicates(true).Size(5)
	searchSource := elastic.NewSearchSource().
		Suggester(suggester).FetchSource(false)

	result, err := es.es.Search().
		Index(es.indexName).
		SearchSource(searchSource).
		Do(context.TODO())

	if err != nil {
		es.h.Errorf("suggest error: %v", err)
		return nil, err
	}

	results := make([]string, 0)
	suggest := result.Suggest["tags-suggest"]
	for _, options := range suggest {
		for _, option := range options.Options {
			results = append(results, option.Text)
		}
	}
	es.h.Debug("suggest results: %v", results)
	return results, nil
}

func (es *ElasticSearch) Delete(alias string, oid string) error {
	es.h.Debugf("delete alias: %s oid: %s", alias, oid)
	resp, err := es.es.Search().
		Index(es.indexName).Query(
		elastic.NewBoolQuery().
			Filter(elastic.NewTermsQuery("alias", alias),
				elastic.NewTermsQuery("oid", oid),
			),
	).FetchSource(false).
		Do(context.TODO())

	if err != nil {
		es.h.Errorf("search error: %v", err)
		return err
	}

	if len(resp.Hits.Hits) > 0 {
		_, err := es.es.Delete().Index(es.indexName).Id(alias).Refresh("true").Do(context.TODO())
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
