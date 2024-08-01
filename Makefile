mockgen:
	mockgen -source=internal/core/ports/inbound/user_service.go -destination=internal/core/ports/inbound/mock/mock_user_service.go -package=mock_user