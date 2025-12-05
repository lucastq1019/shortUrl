package service

import (
	"context"
	"log"
	"time"

	"github.com/username/shorturl/internal/model"
	"github.com/username/shorturl/internal/repository"
	shorturlpb "github.com/username/shorturl/internal/rpc/proto"
	"github.com/username/shorturl/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo repository.URLRepository

	asyncQueue chan *model.ShortURL
}

func (s *Service) CreateShortLink(ctx context.Context, longURL string, expiresIn *time.Duration) (*model.ShortURL, error) {
	// 1. 验证URL
	isValide := utils.ValidateURL(longURL)
	if !isValide {
		return nil, status.Error(codes.InvalidArgument, "不是合法的 LonURL")
	}

	// 2. 生成短码
	shortCode, err := utils.GenerateShortCode(6)
	if err != nil {
		return nil, err
	}
	// 3. 创建模型
	createdAt := time.Now()
	var expiresAt time.Time
	if expiresIn != nil && *expiresIn != 0 {
		expiresAt = createdAt.Add(*expiresIn)
	}

	shortURLModel := &model.ShortURL{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: createdAt,
		ExpiresAt: &expiresAt,
	}
	// 4. 写入缓存（同步，必须成功）
	// 5. 异步写入数据库
	dataSources, err := repository.GetDataSources()
	if err != nil {
		return nil, err
	}
	urlRepository := repository.NewURLRepository(dataSources)
	urlRepository.Save(ctx, shortURLModel)

	log.Println(shortURLModel)

	// 6. 返回结果
	return shortURLModel, nil
}

func (s *Service) GetLongURL(ctx context.Context, req *shorturlpb.GetLongURLRequest) (*shorturlpb.GetLongURLResponse, error) {
	shortKey := req.ShortKey

	// 4. 写入缓存（同步，必须成功）
	// 5. 异步写入数据库
	dataSources, err := repository.GetDataSources()
	if err != nil {
		return nil, err
	}
	urlRepository := repository.NewURLRepository(dataSources)

	shortUrLModel, err := urlRepository.Get(ctx, shortKey)
	if err != nil {
		return nil, err
	}
	resp := &shorturlpb.GetLongURLResponse{LongUrl: shortUrLModel.LongURL}

	// 6. 返回结果
	return resp, nil
}

func (s *Service) GetAllShortLink(ctx context.Context, req *shorturlpb.GetAllShortLinkRequest) (*shorturlpb.GetAllShortLinkResponse, error) {

	dataSources, err := repository.GetDataSources()
	if err != nil {
		return nil, err
	}
	urlRepository := repository.NewURLRepository(dataSources)

	shortURLModels, err := urlRepository.GetAll(ctx, "short")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result []*shorturlpb.ShortLink
	for _, v := range *shortURLModels {
		result = append(result, &shorturlpb.ShortLink{
			ShortLink: v.ShortCode,
			LongLink:  v.LongURL,
		})
	}

	resp := &shorturlpb.GetAllShortLinkResponse{ShortLinks: result}

	// 6. 返回结果
	return resp, nil
}
