package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	biv1 "bingo/api/bi/v1"
	bsv1 "bingo/api/bs/v1"
	"bingo/app/bs/internal/alias"
	"bingo/app/bs/internal/conf"
	"bingo/app/bs/internal/data"
	"bingo/app/bs/internal/search"
	"bingo/pkg/rabbitmq"

	"bingo/app/bs/internal/cache"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShortUrlService struct {
	alias     alias.IAlias
	store     data.IShortUrlStore
	cache     cache.ICache
	search    search.ISearch
	publisher *rabbitmq.Publisher
	h         *log.Helper
}

func NewShortUrlService(c *conf.Bootstrap, h *log.Helper) (*ShortUrlService, error) {
	var _alias alias.IAlias = nil
	var _store data.IShortUrlStore = nil
	var _cache cache.ICache = nil
	var _search search.ISearch = nil
	var _publisher *rabbitmq.Publisher = nil

	var err error = nil

	// create an alias generator which will generate unique alias for short url
	if _alias, err = alias.NewMsAlias(c.Alias, h); err != nil {
		return nil, err
	}

	// database
	switch c.Store.Driver {
	case "sqlite":
		dsn := "file::memory:?cache=shared"
		_store, err = data.NewSqliteBSStore(dsn)
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			c.Store.Username,
			c.Store.Password,
			c.Store.Host,
			c.Store.Port,
			c.Store.Database)
		_store, err = data.NewMysqlBSStore(dsn)
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Store.Host,
			c.Store.Port,
			c.Store.Username,
			c.Store.Password,
			c.Store.Database)
		_store, err = data.NewPostgresBSStore(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?data=%s",
			c.Store.Username,
			c.Store.Password,
			c.Store.Host,
			c.Store.Port,
			c.Store.Database)
		_store, err = data.NewSQLServerBSStore(dsn)
	case "mongo":
		dsn := fmt.Sprintf(
			"mongodb://%s:%d/%s",
			c.Store.Host,
			c.Store.Port,
			c.Store.Database,
		)
		_store, err = data.NewMongoBSStore(
			dsn,
			c.Store.Database,
			c.Store.Username,
			c.Store.Password,
		)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", c.Store.Driver)
	}
	if err != nil {
		return nil, err
	}

	// create database, table and index
	if err = _store.Open(); err != nil {
		return nil, err
	}

	// use redis for cache
	if c.Cache.Addr != "" {
		if _cache, err = cache.NewRedis(c.Cache, h); err != nil {
			h.Debugf("no redis cache avaiable: %v", err)
			_cache = nil
		}
	}

	// use elasticsearch for inverted search of tags
	if c.Search.Addr != nil {
		if _search, err = search.NewElasticSearch(
			c.Search, h); err != nil {
			h.Debugf("no elasticsearch available: %v", err)
			_search = nil
		}
	}

	// use rabbitmq for processing click log
	if c.Publisher.AmqpUri != "" {
		if _publisher, err = rabbitmq.NewPublisher(c.Publisher, h); err != nil {
			h.Debugf("no rabbitmq available: %v", err)
			_publisher = nil
		}
	}
	shortUrlService := &ShortUrlService{
		alias:     _alias,
		store:     _store,
		cache:     _cache,
		search:    _search,
		publisher: _publisher,
		h:         h,
	}

	return shortUrlService, nil
}

func (s *ShortUrlService) CreateShortUrl(
	oid string,
	request *bsv1.CreateShortUrlRequest,
) (*bsv1.ShortUrlResponse, error) {
	alias := request.Alias
	customized := request.Alias != ""

	// validate url
	_, err := url.ParseRequestURI(request.Url)
	if err != nil {
		s.h.Errorf("invalid url: %v", err)
		return nil, err
	}

	// generate a new alias if customized alias is not requested
	if alias == "" {
		_alias, err := s.alias.Next()
		if err != nil {
			s.h.Errorf("get next alias failed: %v", err)
			return nil, err
		}
		alias = _alias
	} else if _, err := s.store.GetShortUrl(alias); err == nil {
		s.h.Errorf("alias %s already exists", alias)
		return nil, fmt.Errorf("alias %s already exists", alias)
	}

	// save short url to database
	err = s.store.CreateShortUrl(alias, customized, request.Url, oid)
	if err != nil {
		s.h.Errorf("create short url failed: %v", err)
		return nil, err
	}

	// get the short url from database and return it back
	shortUrl, err := s.store.GetShortUrl(alias)
	if err != nil {
		s.h.Errorf("get short url after creating failed: %v", err)
		return nil, err
	}
	// add alias to bloom filter if redis cache is enabled
	if s.cache != nil {
		added, err := s.cache.BFAdd("alias", alias)
		if err == nil && added {
			s.h.Debugf("added alias %s to bloom filter", alias)
		}
	}

	// return created short url
	v := &bsv1.ShortUrlResponse{
		Alias:          shortUrl.Alias,
		Url:            shortUrl.Url,
		Title:          shortUrl.Title,
		Tags:           shortUrl.Tags,
		FraudDetection: shortUrl.FraudDetection(),
		Disabled:       shortUrl.Disabled(),
		NoReferrer:     shortUrl.NoReferrer(),
		UtmSource:      shortUrl.UtmSource,
		UtmMedium:      shortUrl.UtmMedium,
		UtmCampaign:    shortUrl.UtmCampaign,
		UtmTerm:        shortUrl.UtmTerm,
		UtmContent:     shortUrl.UtmContent,
		CreatedAt:      timestamppb.New(shortUrl.CreatedAt),
	}

	return v, nil
}

func (s *ShortUrlService) ListShortUrl(
	oid string,
	request *bsv1.ListShortUrlRequest,
) (*bsv1.ListShortUrlResponse, error) {
	r := &bsv1.ListShortUrlResponse{
		Value: make([]*bsv1.ShortUrlResponse, 0),
		Start: request.Start,
		Count: 0,
	}

	s.h.Debugf("list short url request: %v", request)

	shortUrls, err := s.store.ListShortUrl(oid, request.Start, request.Count)
	if err != nil {
		s.h.Errorf("list short url failed: %v", err)
		return nil, err
	}

	for _, it := range shortUrls {
		r.Value = append(
			r.Value,
			&bsv1.ShortUrlResponse{
				Alias:          it.Alias,
				Url:            it.Url,
				Title:          it.Title,
				Tags:           it.Tags,
				FraudDetection: it.FraudDetection(),
				Disabled:       it.Disabled(),
				NoReferrer:     it.NoReferrer(),
				UtmSource:      it.UtmSource,
				UtmMedium:      it.UtmMedium,
				UtmCampaign:    it.UtmCampaign,
				UtmTerm:        it.UtmTerm,
				UtmContent:     it.UtmContent,
				CreatedAt:      timestamppb.New(it.CreatedAt),
			},
		)
	}

	r.Count = int64(len(r.Value))

	// return short urls
	return r, nil
}

func (s *ShortUrlService) GetShortUrl(
	oid string,
	request *bsv1.GetShortUrlRequest,
) (*bsv1.ShortUrlResponse, error) {

	shortUrl, err := s.store.GetShortUrlByOid(oid, request.Alias)
	if err != nil {
		s.h.Errorf("get short url failed: %v", err)
		return nil, err
	}

	// return short url
	return &bsv1.ShortUrlResponse{
		Alias:          shortUrl.Alias,
		Url:            shortUrl.Url,
		Title:          shortUrl.Title,
		Tags:           shortUrl.Tags,
		FraudDetection: shortUrl.FraudDetection(),
		Disabled:       shortUrl.Disabled(),
		NoReferrer:     shortUrl.NoReferrer(),
		UtmSource:      shortUrl.UtmSource,
		UtmMedium:      shortUrl.UtmMedium,
		UtmCampaign:    shortUrl.UtmCampaign,
		UtmTerm:        shortUrl.UtmTerm,
		UtmContent:     shortUrl.UtmContent,
		CreatedAt:      timestamppb.New(shortUrl.CreatedAt),
	}, err
}

func (s *ShortUrlService) UpdateShortUrl(
	oid string,
	request *bsv1.UpdateShortUrlRequest,
) (*bsv1.ShortUrlResponse, error) {
	err := s.store.UpdateShortUrl(
		request.Alias,
		oid,
		data.UpdateShortUrl{
			Url:   request.Url,
			Title: request.Title,
			Tags:  request.Tags,
			Flags: data.Bits(0).
				Set(data.FLAG_FRAUD_DETECTION, request.FraudDetection).
				Set(data.FLAG_DISABLED, request.Disabled).
				Set(data.FLAG_NO_REFERRER, request.NoReferrer),
			UtmSource:   request.UtmSource,
			UtmMedium:   request.UtmMedium,
			UtmCampaign: request.UtmCampaign,
			UtmTerm:     request.UtmTerm,
			UtmContent:  request.UtmContent},
	)

	if err != nil {
		s.h.Errorf("update short url failed: %v", err)
		return nil, err
	}

	// get updated short url from database and return it back
	shortUrl, err := s.store.GetShortUrl(request.Alias)
	if err != nil {
		s.h.Errorf("get short url after updating failed: %v", err)
		return nil, err
	}

	// cache aside algorithm, invalidate cache
	if s.cache != nil {
		s.h.Debug("invalidate cache")
		s.cache.Delete(request.Alias)
	}

	// index it to elasticsearch
	if s.search != nil {
		if len(shortUrl.Tags) > 0 {
			s.search.Index(search.Alias{Alias: shortUrl.Alias, Oid: oid, Tags: shortUrl.Tags})
		}
	}

	// return short url
	return &bsv1.ShortUrlResponse{
		Alias:          shortUrl.Alias,
		Url:            shortUrl.Url,
		Title:          shortUrl.Title,
		Tags:           shortUrl.Tags,
		FraudDetection: shortUrl.FraudDetection(),
		Disabled:       shortUrl.Disabled(),
		NoReferrer:     shortUrl.NoReferrer(),
		UtmSource:      shortUrl.UtmSource,
		UtmMedium:      shortUrl.UtmMedium,
		UtmCampaign:    shortUrl.UtmCampaign,
		UtmTerm:        shortUrl.UtmTerm,
		UtmContent:     shortUrl.UtmContent,
		CreatedAt:      timestamppb.New(shortUrl.CreatedAt),
	}, nil

}

func (s *ShortUrlService) DeleteShortUrl(oid string, request *bsv1.DeleteShortUrlRequest) error {
	s.h.Debugf("delete short url: %v", request)
	err := s.store.DeleteShortUrl(request.Alias, oid)
	if err != nil {
		s.h.Errorf("delete short url failed: %v", err)
		return err
	}
	// cache aside algorithm, invalidate cache
	if s.cache != nil {
		s.cache.Delete(request.Alias)
	}

	// delete from elasticsearh as well
	if s.search != nil {
		s.search.Delete(request.Alias, oid)
	}
	return nil
}

func (s *ShortUrlService) GetCachedShortUrl(alias string) (*CachedShortUrl, error) {
	// validate alias
	valid := s.alias.Validate(alias)
	if !valid {
		s.h.Errorf("invalid alias: %s", alias)
		return nil, fmt.Errorf("invalid alias: %s", alias)
	}

	// use bloom filter to check if alias exists
	// error could be reported if there is any network glitch, ignore it to continue
	if s.cache != nil {
		exists, err := s.cache.BFExists("alias", alias)
		if !exists && err == nil {
			s.h.Errorf("alias %s doesn't exist in bloom filter", alias)
			return nil, fmt.Errorf("alias %s doesn't exist in bloom filter", alias)
		}
	}

	// get alias from cache first, ingore error as database is our last resort
	if s.cache != nil {
		var cachedShortUrl CachedShortUrl
		if err := s.cache.Get(alias, &cachedShortUrl); err == nil {
			return &cachedShortUrl, nil
		}
	}

	// last resort, get short url from database and update cache
	shortUrl, err := s.store.GetShortUrl(alias)
	if err == nil {
		cachedShortUrl := CachedShortUrl{
			Url:            shortUrl.Url,
			FraudDetection: shortUrl.FraudDetection(),
			Disabled:       shortUrl.Disabled(),
			NoReferrer:     shortUrl.NoReferrer(),
			UtmSource:      shortUrl.UtmSource,
			UtmMedium:      shortUrl.UtmMedium,
			UtmCampaign:    shortUrl.UtmCampaign,
			UtmTerm:        shortUrl.UtmTerm,
			UtmContent:     shortUrl.UtmContent}
		if s.cache != nil {
			s.cache.Set(
				alias,
				&cachedShortUrl)
		}
		return &cachedShortUrl, nil
	}
	return nil, err
}

func (ss *ShortUrlService) Click(alias string, ip string, ua string, referer string) error {
	clickEvent := biv1.ClickEvent{
		Alias:     alias,
		Ip:        ip,
		Ua:        ua,
		Referer:   referer,
		CreatedAt: timestamppb.New(time.Now()),
	}
	data, err := proto.Marshal(&clickEvent)
	if err != nil {
		return err
	}
	return ss.publisher.Publish(data)
}

func (ss *ShortUrlService) SuggestedTags(
	request *bsv1.SuggestRequest,
) (*bsv1.SuggestResponse, error) {
	response := bsv1.SuggestResponse{Value: make([]string, 0)}
	var err error = nil
	if ss.search != nil {
		response.Value, err = ss.search.Suggest(request.Query)
	}
	return &response, err
}

func (ss *ShortUrlService) SiteVerify(
	request *bsv1.VerifyRequest,
	secretKey string,
) (*bsv1.VerifyResponse, error) {
	// google re-captcha requires form based post, send for verification
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {secretKey},
		"response": {request.Token},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// decode response
	var siteVerifyResponse SiteVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&siteVerifyResponse); err != nil {
		return nil, err
	}
	// check re-captcha verification result
	if !siteVerifyResponse.Success {
		return nil, errors.New("recaptcha verify failed")
	}

	cachedShortUrl, err := ss.GetCachedShortUrl(request.Alias)
	if err != nil {
		return nil, err
	}

	return &bsv1.VerifyResponse{
		Score: siteVerifyResponse.Score,
		Url:   cachedShortUrl.GeneratedUrl(),
	}, nil

}

func (ss *ShortUrlService) CountClicks(
	c *conf.BI,
	request *bsv1.ClicksRequest,
) (*bsv1.ClicksReply, error) {
	if c.Protocol == conf.BI_HTTP {
		conn, err := transhttp.NewClient(
			context.Background(),
			transhttp.WithMiddleware(
				recovery.Recovery(),
			),
			transhttp.WithEndpoint(c.HttpAddr),
		)

		if err != nil {
			return nil, err
		}

		defer conn.Close()
		client := biv1.NewBIHTTPClient(conn)
		r, err := client.Clicks(context.Background(), &biv1.ClicksRequest{Alias: request.Alias})
		if err != nil {
			return nil, err
		}
		return &bsv1.ClicksReply{Clicks: r.Clicks}, nil
	}

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint(c.HttpAddr),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := biv1.NewBIClient(conn)
	r, err := client.Clicks(context.Background(), &biv1.ClicksRequest{Alias: request.Alias})

	if err != nil {
		return nil, err
	}
	return &bsv1.ClicksReply{Clicks: r.Clicks}, nil
}

func (ss *ShortUrlService) Close() {
	ss.alias.Close()
	ss.store.Close()
}
