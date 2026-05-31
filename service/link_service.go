package service

import (
	"link-storage-service/cache"
	"link-storage-service/model"
	"link-storage-service/repository"
	"math/rand/v2"
)

type LinkService struct {
	repo       *repository.LinkRepository
	cache      *cache.MemoryCache
	codeLength int
	locker     *KeyLocker
}

func New(repo *repository.LinkRepository, cache *cache.MemoryCache, codeLength int) *LinkService {
	return &LinkService{repo: repo, cache: cache, codeLength: codeLength, locker: NewKeyLocker()}
}

func (s *LinkService) List(limit int, offset int) ([]model.Link, error) {
	return s.repo.List(limit, offset)
}

func (s *LinkService) CreateLink(originalURL string) (string, error) {
	code := GenerateRandomCode(s.codeLength)
	link, err := s.repo.Insert(model.Link{ShortCode: code, OriginalURL: originalURL})
	return link.ShortCode, err
}

func (s *LinkService) GetLinkAndIncrementVisits(shortCode string) (model.Link, error) {
	link, err := s.GetLinkByShortCode(shortCode)
	if err != nil {
		return model.Link{}, err
	}

	unlock := s.locker.Lock(shortCode)
	defer unlock()

	visits, err := s.repo.IncrementVisits(shortCode)
	if err != nil {
		return model.Link{}, err
	}
	link.Visits = visits
	s.cache.Set(link)

	return link, nil
}

func (s *LinkService) GetLinkByShortCode(shortCode string) (model.Link, error) {
	cachedLink, ok := s.cache.Get(shortCode)
	if ok {
		return cachedLink, nil
	}
	return s.repo.GetByShortCode(shortCode)
}

func (s *LinkService) DeleteLink(shortCode string) error {
	err := s.repo.Delete(shortCode)
	if err != nil {
		return err
	}
	s.cache.Delete(shortCode)
	return nil
}

func GenerateRandomCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}

	return string(b)
}
