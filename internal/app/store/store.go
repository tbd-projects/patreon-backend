package store

type Store interface {
	User() UserRepository
	Creator() CreatorRepository
}
