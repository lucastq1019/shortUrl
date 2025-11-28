package handler

import "github.com/username/shorturl/internal/manager"

type ShortenerService struct {
	CM *manager.ClientGetter
}

func (s *ShortenerService) CreateShortLink() {
	s.CM.GetClient
}
