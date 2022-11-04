package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	bev1 "bingo/api/be/v1"
	bsv1 "bingo/api/bs/v1"

	"bingo/app/bs/internal/conf"
	"bingo/app/bs/internal/service"

	"bingo/app/bs/internal/ws"

	"bingo/app/bs/internal/system"

	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/xenitab/go-oidc-middleware/oidcgin"
	opts "github.com/xenitab/go-oidc-middleware/options"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	corsConfig = cors.Config{
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowOrigins:     []string{"*"},
	}
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "bs"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(
		&flagconf,
		"conf",
		"config.yaml",
		"config path, eg: -conf config.yaml",
	)
}

func openIDHandler(c *conf.JWT) gin.HandlerFunc {
	var opts []opts.Option = []opts.Option{
		opts.WithIssuer(c.Issuer),
		opts.WithRequiredTokenType("JWT"),
		opts.WithRequiredAudience(c.Audience),
		opts.WithFallbackSignatureAlgorithm(c.FallbackSignatureAlgorithm),
	}

	return oidcgin.New(opts...)
}

func getOID(c *gin.Context) string {
	var oid interface{} = nil
	claimsValue, found := c.Get("claims")
	if found {
		claims, ok := claimsValue.(map[string]interface{})
		if ok {
			oid = claims["oid"]
			if oid != nil {
				return oid.(string)
			}
		}
	}
	return ""
}

// initTracer init jaeger tracer provider
func initTracer(url string) error {
	// create the jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	//  	propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)
	h := log.NewHelper(log.NewFilter(logger, log.FilterLevel(log.LevelDebug)))

	c := config.New(
		config.WithSource(
			env.NewSource("BS_"),
			file.NewSource(flagconf),
		),
		config.WithResolver(conf.BsResolver),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var sentinel_addrs []string
	for _, e := range bc.Cache.SentinelAddrs {
		sentinel_addrs = append(sentinel_addrs, strings.Split(e, ",")...)
	}
	bc.Cache.SentinelAddrs = sentinel_addrs

	var es_addrs []string
	for _, e := range bc.Search.Addrs {
		es_addrs = append(es_addrs, strings.Split(e, ",")...)
	}
	bc.Search.Addrs = es_addrs

	var etcd_addrs []string
	for _, e := range bc.Alias.EtcdAddrs {
		etcd_addrs = append(etcd_addrs, strings.Split(e, ",")...)
	}
	bc.Alias.EtcdAddrs = etcd_addrs

	h.Debugf("dump config: %v", &bc)

	if err := initTracer(bc.Jaeger.Addr); err != nil {
		panic(err)
	}

	shortUrlService, err := service.NewShortUrlService(&bc, h)
	if err != nil {
		panic(err)
	}

	// start cpu, memory load and system monitor, monitor uses web socket to return
	// cpu and memory load to client
	cpu := system.NewCpuLoad(runtime.NumCPU(), h)
	cpu.Start()
	memory := system.NewMemLoad(h)
	memory.Start()
	monitor := ws.NewMonitor(h)
	monitor.Start()

	// gin debug mode
	if bc.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if bc.Server.Debug {
		r.Use(gin.Logger(), gin.Recovery())
	} else {
		r.Use(gin.Recovery())
	}

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Namespace("bingo"),
		ginprom.Subsystem("bs"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())

	//r.Use(otelgin.Middleware("bs", otelgin.WithPropagators(p)))
	r.Use(otelgin.Middleware("bs"))

	// multi-template
	mt := multitemplate.NewRenderer()
	mt.Add("Captcha", CAPTCHA_TEMPLATE)
	mt.Add("WsDebug", WS_DEBUG_TEMPLATE)
	mt.Add("Expand", EXPAND_TEMPLATE)
	r.HTMLRender = mt
	r.Use(cors.New(corsConfig))
	r.Use(static.Serve("/", static.LocalFile("website", false)))
	r.StaticFile("/", "website/index.html")

	// to support react router
	// redirect to index html so react router gets chance to handle the browser routing
	r.NoRoute(func(ctx *gin.Context) {
		h.Debugf("no route: %s", ctx.Request.RequestURI)
		dir, _ := path.Split(ctx.Request.RequestURI)
		if dir == "/pages/" {
			ctx.File("./website/index.html")
		}
	})

	r.GET("/:alias", func(ctx *gin.Context) {
		alias := ctx.Param("alias")
		ext := filepath.Ext(alias)
		if alias != "" && ext == "" {
			// check if it is an expand request
			expand := alias[len(alias)-1:] == "+"
			if expand {
				alias = alias[:len(alias)-1]
			}

			cachedShortUrl, err := shortUrlService.GetCachedShortUrl(alias)

			// return 404 directly instead of detailed error to avoid attack
			if err != nil {
				ctx.String(http.StatusNotFound, "404 page not found")
				h.Error(err)
				// ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			// if short url is disabled, return 404 directly
			if cachedShortUrl.Disabled {
				ctx.String(http.StatusNotFound, "404 page not found")
				return
			}

			// use extrace service (be) to reterive title, keywords, summary
			if expand {
				otelCtx := ctx.Request.Context()
				conn, err := transhttp.NewClient(
					otelCtx,
					transhttp.WithTimeout(30*time.Second),
					transhttp.WithMiddleware(
						recovery.Recovery(),
						tracing.Client(),
					),
					transhttp.WithEndpoint(bc.Be.HttpAddr),
				)
				// return 404 directly instead of detailed error to avoid attack
				if err != nil {
					ctx.String(http.StatusNotFound, "404 page not found")
					// ctx.String(http.StatusInternalServerError, err.Error())
					h.Error(err)
					return
				}
				defer conn.Close()
				client := bev1.NewBEHTTPClient(conn)
				r, err := client.Extract(
					otelCtx,
					&bev1.ExtractRequest{Url: cachedShortUrl.Url},
				)
				// return 404 directly instead of detailed error to avoid attack
				if err != nil {
					ctx.String(http.StatusNotFound, "404 page not found")
					h.Error(err)
					// ctx.String(http.StatusInternalServerError, err.Error())
					return
				}

				var base64Encoding string = NO_IMAGE_AVAILABLE

				// load url snapshot from gowitness
				var data = []byte(fmt.Sprintf("{\"url\": \"%s\",\"oneshot\": \"true\"}", cachedShortUrl.Url))
				resp, err := otelhttp.Post(
					otelCtx,
					"http://localhost:7171/api/screenshot",
					"application/json",
					bytes.NewBuffer(data),
				)
				if err == nil {
					defer resp.Body.Close()
					bytes, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						base64Encoding = base64.StdEncoding.EncodeToString(bytes)
					} else {
						h.Error(err)
					}
				} else {
					h.Error(err)
				}

				otelgin.HTML(ctx, http.StatusOK, "Expand", gin.H{
					"Alias":    alias,
					"Url":      cachedShortUrl.Url,
					"Title":    r.Title,
					"Keywords": strings.Join(r.Keywords, " "),
					"Summary":  r.Summary,
					"Snapshot": base64Encoding,
				})
				return
			}

			// if fraud detection is enabled, return re-captcha page
			if cachedShortUrl.FraudDetection {
				referrer := "strict-origin-when-cross-origin"
				if cachedShortUrl.NoReferrer {
					referrer = "no-referrer"
				}
				otelgin.HTML(ctx, http.StatusOK, "Captcha", gin.H{
					"Alias":            alias,
					"RecaptchaSiteKey": bc.Recaptcha.SiteKey,
					"Referrer":         referrer})
				return
			}
			// redirect
			go shortUrlService.Click(
				alias,
				ctx.ClientIP(),
				ctx.Request.UserAgent(),
				ctx.Request.Referer(),
			)
			ctx.Redirect(http.StatusFound, cachedShortUrl.GeneratedUrl())
		} else {
			ctx.String(http.StatusNotFound, "404 page not found")
		}
	})

	// web socket handling
	r.GET("/ws/debug", func(ctx *gin.Context) {
		otelgin.HTML(ctx, http.StatusOK, "WsDebug", gin.H{
			"Host": ctx.Request.Host})
	})
	r.GET("/ws", func(ctx *gin.Context) {
		monitor.Serve(ctx.Writer, ctx.Request)
	})

	// recaptcha
	r.POST("/v1/captcha/verify", func(ctx *gin.Context) {
		// get re-captcha token
		request := bsv1.VerifyRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		response, err := shortUrlService.SiteVerify(&request, bc.Recaptcha.SecretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, bsv1.ErrorSiteVerifyFailed(err.Error()))
			return
		}
		// non-blocking, as publish to rabbitmq is a blocked operation
		go shortUrlService.Click(
			request.Alias,
			ctx.ClientIP(),
			ctx.Request.UserAgent(),
			ctx.Request.Referer(),
		)

		ctx.JSON(http.StatusOK, response)
	})

	// use oidc to verify jwt token in following apis
	openIdHandler := openIDHandler(bc.Jwt)
	v1 := r.Group("/v1", openIdHandler)

	// ping is used to measure client latency
	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// handle system statistic request
	v1.GET("/system/stats", func(ctx *gin.Context) {
		response, err := system.GetStats(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// handle cpu load request
	v1.POST("/system/cpu", func(ctx *gin.Context) {
		request := bsv1.CpuLoadRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
		} else {
			cpu.Update(&request)
			ctx.Status(http.StatusOK)
		}
	})

	// handle memory load request
	v1.POST("/system/memory", func(ctx *gin.Context) {
		request := bsv1.MemLoadRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
		} else {
			memory.Update(&request)
			ctx.Status(http.StatusOK)
		}
	})

	// handle create short url request
	v1.POST("/shorturl", func(ctx *gin.Context) {
		request := bsv1.CreateShortUrlRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		oid := getOID(ctx)
		var response *bsv1.ShortUrlResponse
		if response, err = shortUrlService.CreateShortUrl(oid, &request); err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// handle list short url request
	v1.GET("/shorturl", func(ctx *gin.Context) {
		type ListShortUrlQuery struct {
			Start int64 `form:"start"`
			Count int64 `form:"count"`
		}
		query := ListShortUrlQuery{Start: 0, Count: bc.Server.PageSize}
		if err := ctx.BindQuery(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		// restrict requested page size to maximum page size configured at server
		if query.Count > bc.Server.PageSize {
			query.Count = bc.Server.PageSize
		}
		oid := getOID(ctx)
		response, err := shortUrlService.ListShortUrl(oid, &bsv1.ListShortUrlRequest{
			Start: query.Start,
			Count: query.Count,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// handle get short url request
	v1.GET("/shorturl/:alias", func(ctx *gin.Context) {
		request := bsv1.GetShortUrlRequest{Alias: ctx.Param("alias")}
		oid := getOID(ctx)
		response, err := shortUrlService.GetShortUrl(oid, &request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// handle update short url request
	v1.PUT("/shorturl", func(ctx *gin.Context) {
		request := bsv1.UpdateShortUrlRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		oid := getOID(ctx)
		response, err := shortUrlService.UpdateShortUrl(oid, &request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// handle delete short url request
	v1.DELETE("/shorturl/:alias", func(ctx *gin.Context) {
		request := bsv1.DeleteShortUrlRequest{
			Alias: ctx.Param("alias"),
		}
		oid := getOID(ctx)
		err := shortUrlService.DeleteShortUrl(oid, &request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, bsv1.ErrorInternalServerError(err.Error()))
			return
		}
		ctx.Status(http.StatusOK)
	})

	// handle count clicks request
	v1.GET("/shorturl-bi/clicks/:alias", func(ctx *gin.Context) {
		request := bsv1.ClicksRequest{Alias: ctx.Param("alias")}
		response, err := shortUrlService.CountClicks(bc.Bi, &request, ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusBadGateway, bsv1.ErrorBadGateway(err.Error()))
		} else {
			ctx.JSON(http.StatusOK, response)
		}
	})

	// handle suggested tag request
	v1.GET("/tag-suggest/:query", func(ctx *gin.Context) {
		request := bsv1.SuggestRequest{
			Query: ctx.Param("query"),
		}
		response, err := shortUrlService.SuggestedTags(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	// start http server
	srv := &http.Server{
		Addr:    bc.Server.Http.Addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// gracefully shutdown server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	h.Info("shutdown server ...")
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	shortUrlService.Close()
	h.Info("good-bye!")
}
