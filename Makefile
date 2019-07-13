build: fmt vet
	GOOS=windows GOARCH=amd64 go build -o bin/ansible-playbook.exe github.com/kaitoy/ansible-playbook-wrapper

fmt:
	go fmt ./pkg/... ./cmd/...

vet:
	go vet ./pkg/... ./cmd/...
