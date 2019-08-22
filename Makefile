# note: call scripts from /scripts
.PHONY: proto
proto:
	./scripts/proto.sh

.PHONY: test
test:
	go test ./... -coverprofile .testCoverage.txt