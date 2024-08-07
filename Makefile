mockgen:
	mockgen -source=internal/core/ports/inbound/user_service.go -destination=internal/core/ports/inbound/mock/mock_user_service.go -package=mock_user
	mockgen -source=internal/core/ports/outbound/user_repository.go -destination=internal/core/ports/outbound/mock/mock_user_repository.go -package=mock_repo