go test -tags "sqlite_fts5" -v -cover -coverprofile=coverage.out ./internal/... && \
go tool cover -html=coverage.out -o coverage.html
