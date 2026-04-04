package service

type SetsStorage interface {
}

type SetsService struct {
	Store SetsStorage
}

func NewSetsService(s SetsStorage) *SetsService {
	return &SetsService{
		Store: s,
	}
}

