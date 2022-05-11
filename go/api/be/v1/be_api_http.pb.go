// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.2.1

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type BEHTTPServer interface {
	Extract(context.Context, *ExtractRequest) (*ExtractReply, error)
}

func RegisterBEHTTPServer(s *http.Server, srv BEHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/extract", _BE_Extract0_HTTP_Handler(srv))
}

func _BE_Extract0_HTTP_Handler(srv BEHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ExtractRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.be.v1.BE/Extract")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Extract(ctx, req.(*ExtractRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ExtractReply)
		return ctx.Result(200, reply)
	}
}

type BEHTTPClient interface {
	Extract(ctx context.Context, req *ExtractRequest, opts ...http.CallOption) (rsp *ExtractReply, err error)
}

type BEHTTPClientImpl struct {
	cc *http.Client
}

func NewBEHTTPClient(client *http.Client) BEHTTPClient {
	return &BEHTTPClientImpl{client}
}

func (c *BEHTTPClientImpl) Extract(ctx context.Context, in *ExtractRequest, opts ...http.CallOption) (*ExtractReply, error) {
	var out ExtractReply
	pattern := "/v1/extract"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.be.v1.BE/Extract"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
