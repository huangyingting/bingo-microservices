// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.2.1

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type ShortUrlHTTPServer interface {
	CreateShortUrl(context.Context, *CreateShortUrlRequest) (*ShortUrlResponse, error)
	DeleteShortUrl(context.Context, *DeleteShortUrlRequest) (*emptypb.Empty, error)
	GetShortUrl(context.Context, *GetShortUrlRequest) (*ShortUrlResponse, error)
	ListShortUrl(context.Context, *ListShortUrlRequest) (*ListShortUrlResponse, error)
	UpdateShortUrl(context.Context, *UpdateShortUrlRequest) (*ShortUrlResponse, error)
}

func RegisterShortUrlHTTPServer(s *http.Server, srv ShortUrlHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/shorturl", _ShortUrl_CreateShortUrl0_HTTP_Handler(srv))
	r.PUT("/v1/shorturl", _ShortUrl_UpdateShortUrl0_HTTP_Handler(srv))
	r.GET("/v1/shorturl", _ShortUrl_ListShortUrl0_HTTP_Handler(srv))
	r.GET("/v1/shorturl/{alias}", _ShortUrl_GetShortUrl0_HTTP_Handler(srv))
	r.DELETE("/v1/shorturl/{alias}", _ShortUrl_DeleteShortUrl0_HTTP_Handler(srv))
}

func _ShortUrl_CreateShortUrl0_HTTP_Handler(srv ShortUrlHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateShortUrlRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrl/CreateShortUrl")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateShortUrl(ctx, req.(*CreateShortUrlRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ShortUrlResponse)
		return ctx.Result(200, reply)
	}
}

func _ShortUrl_UpdateShortUrl0_HTTP_Handler(srv ShortUrlHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateShortUrlRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrl/UpdateShortUrl")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateShortUrl(ctx, req.(*UpdateShortUrlRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ShortUrlResponse)
		return ctx.Result(200, reply)
	}
}

func _ShortUrl_ListShortUrl0_HTTP_Handler(srv ShortUrlHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListShortUrlRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrl/ListShortUrl")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListShortUrl(ctx, req.(*ListShortUrlRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListShortUrlResponse)
		return ctx.Result(200, reply)
	}
}

func _ShortUrl_GetShortUrl0_HTTP_Handler(srv ShortUrlHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetShortUrlRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrl/GetShortUrl")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetShortUrl(ctx, req.(*GetShortUrlRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ShortUrlResponse)
		return ctx.Result(200, reply)
	}
}

func _ShortUrl_DeleteShortUrl0_HTTP_Handler(srv ShortUrlHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteShortUrlRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrl/DeleteShortUrl")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteShortUrl(ctx, req.(*DeleteShortUrlRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type ShortUrlHTTPClient interface {
	CreateShortUrl(ctx context.Context, req *CreateShortUrlRequest, opts ...http.CallOption) (rsp *ShortUrlResponse, err error)
	DeleteShortUrl(ctx context.Context, req *DeleteShortUrlRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetShortUrl(ctx context.Context, req *GetShortUrlRequest, opts ...http.CallOption) (rsp *ShortUrlResponse, err error)
	ListShortUrl(ctx context.Context, req *ListShortUrlRequest, opts ...http.CallOption) (rsp *ListShortUrlResponse, err error)
	UpdateShortUrl(ctx context.Context, req *UpdateShortUrlRequest, opts ...http.CallOption) (rsp *ShortUrlResponse, err error)
}

type ShortUrlHTTPClientImpl struct {
	cc *http.Client
}

func NewShortUrlHTTPClient(client *http.Client) ShortUrlHTTPClient {
	return &ShortUrlHTTPClientImpl{client}
}

func (c *ShortUrlHTTPClientImpl) CreateShortUrl(ctx context.Context, in *CreateShortUrlRequest, opts ...http.CallOption) (*ShortUrlResponse, error) {
	var out ShortUrlResponse
	pattern := "/v1/shorturl"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrl/CreateShortUrl"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ShortUrlHTTPClientImpl) DeleteShortUrl(ctx context.Context, in *DeleteShortUrlRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/shorturl/{alias}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrl/DeleteShortUrl"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ShortUrlHTTPClientImpl) GetShortUrl(ctx context.Context, in *GetShortUrlRequest, opts ...http.CallOption) (*ShortUrlResponse, error) {
	var out ShortUrlResponse
	pattern := "/v1/shorturl/{alias}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrl/GetShortUrl"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ShortUrlHTTPClientImpl) ListShortUrl(ctx context.Context, in *ListShortUrlRequest, opts ...http.CallOption) (*ListShortUrlResponse, error) {
	var out ListShortUrlResponse
	pattern := "/v1/shorturl"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrl/ListShortUrl"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ShortUrlHTTPClientImpl) UpdateShortUrl(ctx context.Context, in *UpdateShortUrlRequest, opts ...http.CallOption) (*ShortUrlResponse, error) {
	var out ShortUrlResponse
	pattern := "/v1/shorturl"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrl/UpdateShortUrl"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

type CaptchaHTTPServer interface {
	Verify(context.Context, *VerifyRequest) (*VerifyResponse, error)
}

func RegisterCaptchaHTTPServer(s *http.Server, srv CaptchaHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/captcha/verify", _Captcha_Verify0_HTTP_Handler(srv))
}

func _Captcha_Verify0_HTTP_Handler(srv CaptchaHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in VerifyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.Captcha/Verify")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Verify(ctx, req.(*VerifyRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*VerifyResponse)
		return ctx.Result(200, reply)
	}
}

type CaptchaHTTPClient interface {
	Verify(ctx context.Context, req *VerifyRequest, opts ...http.CallOption) (rsp *VerifyResponse, err error)
}

type CaptchaHTTPClientImpl struct {
	cc *http.Client
}

func NewCaptchaHTTPClient(client *http.Client) CaptchaHTTPClient {
	return &CaptchaHTTPClientImpl{client}
}

func (c *CaptchaHTTPClientImpl) Verify(ctx context.Context, in *VerifyRequest, opts ...http.CallOption) (*VerifyResponse, error) {
	var out VerifyResponse
	pattern := "/v1/captcha/verify"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.shorturl.v1.Captcha/Verify"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

type TagSuggestHTTPServer interface {
	Verify(context.Context, *SuggestRequest) (*SuggestResponse, error)
}

func RegisterTagSuggestHTTPServer(s *http.Server, srv TagSuggestHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/tag-suggest/{query}", _TagSuggest_Verify1_HTTP_Handler(srv))
}

func _TagSuggest_Verify1_HTTP_Handler(srv TagSuggestHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SuggestRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.TagSuggest/Verify")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Verify(ctx, req.(*SuggestRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SuggestResponse)
		return ctx.Result(200, reply)
	}
}

type TagSuggestHTTPClient interface {
	Verify(ctx context.Context, req *SuggestRequest, opts ...http.CallOption) (rsp *SuggestResponse, err error)
}

type TagSuggestHTTPClientImpl struct {
	cc *http.Client
}

func NewTagSuggestHTTPClient(client *http.Client) TagSuggestHTTPClient {
	return &TagSuggestHTTPClientImpl{client}
}

func (c *TagSuggestHTTPClientImpl) Verify(ctx context.Context, in *SuggestRequest, opts ...http.CallOption) (*SuggestResponse, error) {
	var out SuggestResponse
	pattern := "/v1/tag-suggest/{query}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.TagSuggest/Verify"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

type SystemHTTPServer interface {
	Stats(context.Context, *emptypb.Empty) (*StatsResponse, error)
	UpdateCpuLoad(context.Context, *CpuLoadRequest) (*emptypb.Empty, error)
	UpdateMemLoad(context.Context, *MemLoadRequest) (*emptypb.Empty, error)
}

func RegisterSystemHTTPServer(s *http.Server, srv SystemHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/system/mem", _System_UpdateMemLoad0_HTTP_Handler(srv))
	r.PUT("/v1/system/cpu", _System_UpdateCpuLoad0_HTTP_Handler(srv))
	r.GET("/v1/system/stats", _System_Stats0_HTTP_Handler(srv))
}

func _System_UpdateMemLoad0_HTTP_Handler(srv SystemHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in MemLoadRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.System/UpdateMemLoad")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateMemLoad(ctx, req.(*MemLoadRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _System_UpdateCpuLoad0_HTTP_Handler(srv SystemHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CpuLoadRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.System/UpdateCpuLoad")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateCpuLoad(ctx, req.(*CpuLoadRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _System_Stats0_HTTP_Handler(srv SystemHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.System/Stats")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Stats(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*StatsResponse)
		return ctx.Result(200, reply)
	}
}

type SystemHTTPClient interface {
	Stats(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *StatsResponse, err error)
	UpdateCpuLoad(ctx context.Context, req *CpuLoadRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateMemLoad(ctx context.Context, req *MemLoadRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type SystemHTTPClientImpl struct {
	cc *http.Client
}

func NewSystemHTTPClient(client *http.Client) SystemHTTPClient {
	return &SystemHTTPClientImpl{client}
}

func (c *SystemHTTPClientImpl) Stats(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*StatsResponse, error) {
	var out StatsResponse
	pattern := "/v1/system/stats"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.System/Stats"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *SystemHTTPClientImpl) UpdateCpuLoad(ctx context.Context, in *CpuLoadRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/system/cpu"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.shorturl.v1.System/UpdateCpuLoad"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *SystemHTTPClientImpl) UpdateMemLoad(ctx context.Context, in *MemLoadRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/system/mem"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.shorturl.v1.System/UpdateMemLoad"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

type ShortUrlBIHTTPServer interface {
	Clicks(context.Context, *ClicksRequest) (*ClicksReply, error)
}

func RegisterShortUrlBIHTTPServer(s *http.Server, srv ShortUrlBIHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/shorturl-bi/clicks/{alias}", _ShortUrlBI_Clicks0_HTTP_Handler(srv))
}

func _ShortUrlBI_Clicks0_HTTP_Handler(srv ShortUrlBIHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ClicksRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.shorturl.v1.ShortUrlBI/Clicks")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Clicks(ctx, req.(*ClicksRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ClicksReply)
		return ctx.Result(200, reply)
	}
}

type ShortUrlBIHTTPClient interface {
	Clicks(ctx context.Context, req *ClicksRequest, opts ...http.CallOption) (rsp *ClicksReply, err error)
}

type ShortUrlBIHTTPClientImpl struct {
	cc *http.Client
}

func NewShortUrlBIHTTPClient(client *http.Client) ShortUrlBIHTTPClient {
	return &ShortUrlBIHTTPClientImpl{client}
}

func (c *ShortUrlBIHTTPClientImpl) Clicks(ctx context.Context, in *ClicksRequest, opts ...http.CallOption) (*ClicksReply, error) {
	var out ClicksReply
	pattern := "/v1/shorturl-bi/clicks/{alias}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.shorturl.v1.ShortUrlBI/Clicks"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
