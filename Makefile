build: format
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags="-s -w"

format:
	gofmt -s -w .
