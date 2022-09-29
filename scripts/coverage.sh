go test ./... -coverprofile=/tmp/coverage.out;
grep -v "Timestopeofficial/feechain/core" /tmp/coverage.out > /tmp/coverage1.out
go tool cover -func=/tmp/coverage1.out
