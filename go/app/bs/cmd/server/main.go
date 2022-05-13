package main

import (
	"context"
	"flag"
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
		"../../configs/config.yaml",
		"config path, eg: -conf config.yaml",
	)
}

func NewOpenIdHandler(c *conf.JWT) gin.HandlerFunc {
	var opts []opts.Option = []opts.Option{
		opts.WithIssuer(c.Issuer),
		opts.WithRequiredTokenType("JWT"),
		opts.WithRequiredAudience(c.Audience),
		opts.WithFallbackSignatureAlgorithm(c.FallbackSignatureAlgorithm),
	}

	return oidcgin.New(opts...)
}

func GetOid(c *gin.Context) string {
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

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	h := log.NewHelper(log.NewFilter(logger, log.FilterLevel(log.LevelDebug)))

	c := config.New(
		config.WithSource(
			env.NewSource("BS_"),
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	h.Debugf("dump config: %v", &bc)

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
	r := gin.Default()

	// multi-template
	mt := multitemplate.NewRenderer()
	mt.Add("Captcha", CAPTCHA_TEMPLATE)
	mt.Add("WsDebug", WS_DEBUG_TEMPLATE)
	mt.Add("Expand", EXPAND_TEMPLATE)
	r.HTMLRender = mt
	r.Use(cors.New(corsConfig))
	r.Use(static.Serve("/", static.LocalFile("website", false)))
	r.StaticFile("/", "website/index.html")

	// handle short url redirect
	r.NoRoute(func(ctx *gin.Context) {
		h.Debugf("no route: %s", ctx.Request.RequestURI)
		dir, alias := path.Split(ctx.Request.RequestURI)
		ext := filepath.Ext(alias)
		if dir == "/" && alias != "" && ext == "" {
			// check if it is an expand request
			expand := alias[len(alias)-1:] == "+"
			if expand {
				alias = alias[:len(alias)-1]
			}

			cachedShortUrl, err := shortUrlService.GetCachedShortUrl(alias)

			// return 404 directly
			if err != nil {
				return
			}
			// if short url is disabled, return 404 directly
			if cachedShortUrl.Disabled {
				return
			}

			// use extrace service (be) to reterive title, keywords, summary
			if expand {
				conn, err := transhttp.NewClient(
					context.Background(),
					transhttp.WithMiddleware(
						recovery.Recovery(),
					),
					transhttp.WithEndpoint(bc.Be.HttpAddr),
				)

				if err != nil {
					return
				}

				defer conn.Close()
				client := bev1.NewBEHTTPClient(conn)
				r, err := client.Extract(
					context.Background(),
					&bev1.ExtractRequest{Url: cachedShortUrl.Url},
				)
				if err != nil {
					return
				}
				ctx.HTML(http.StatusOK, "Expand", gin.H{
					"Alias":    alias,
					"Url":      cachedShortUrl.Url,
					"Title":    r.Title,
					"Keywords": strings.Join(r.Keywords, " "),
					"Summary":  r.Summary,
					"Snapshot": bc.GoWitness.Addr + "/?url=" + cachedShortUrl.Url,
				})
				return
			}

			// if fraud detection is enabled, return re-captcha page
			if cachedShortUrl.FraudDetection {
				referrer := "strict-origin-when-cross-origin"
				if cachedShortUrl.NoReferrer {
					referrer = "no-referrer"
				}
				ctx.HTML(http.StatusOK, "Captcha", gin.H{
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
		} else if dir == "/pages/" {
			ctx.File("./website/index.html")
		}
	})

	// web socket handling
	r.GET("/ws/debug", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "WsDebug", gin.H{
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
	openIdHandler := NewOpenIdHandler(bc.Jwt)
	v1 := r.Group("/v1", openIdHandler)

	// ping is used to measure client latency
	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// handle system statistic request
	v1.GET("/system/stats", func(ctx *gin.Context) {
		response, err := system.GetStats()
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
		oid := GetOid(ctx)
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
		query := ListShortUrlQuery{Start: 0, Count: 5}
		if err := ctx.BindQuery(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, bsv1.ErrorBadRequest(err.Error()))
			return
		}
		oid := GetOid(ctx)
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
		oid := GetOid(ctx)
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
		oid := GetOid(ctx)
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
		oid := GetOid(ctx)
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
		response, err := shortUrlService.CountClicks(bc.Bi, &request)
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
