package payments

//go:generate mockgen -destination=mocks/mock_payments_repository.go -package=mock_repository -mock_names=Repository=PaymentsRepository . Repository

type Repository interface {
	// GetUserPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetUserPayments(userID int64)
	// GetCreatorPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetCreatorPayments(creatorID int64)
}
