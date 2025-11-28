package service

import (
	"context"
	"log"
	"time"

	shorturlpb "github.com/username/shorturl/internal/rpc/proto"
	shortener "github.com/username/shorturl/internal/service/shortener"
)

type Server struct {
	shorturlpb.UnimplementedShortenerServiceServer
	service *shortener.Service
}

func (s *Server) CreateShortLink(ctx context.Context, req *shorturlpb.CreateShortLinkRequest) (*shorturlpb.CreateShortLinkResponse, error) {
	var time time.Duration = time.Second * 100
	shortURLModel, err := s.service.CreateShortLink(ctx, req.GetLongUrl(), &time)

	log.Println(shortURLModel)
	if err != nil {
		return nil, err
	}
	response := &shorturlpb.CreateShortLinkResponse{}

	if shortURLModel != nil {
		response.ShortKey = shortURLModel.ShortCode
	}

	return response, err
}

func (s *Server) GetLongURL(ctx context.Context, req *shorturlpb.GetLongURLRequest) (*shorturlpb.GetLongURLResponse, error) {
	return s.service.GetLongURL(ctx, req)
}
