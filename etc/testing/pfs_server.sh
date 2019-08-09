# grep out PFS server logs, as otherwise the test output is too verbose to
# follow and breaks travis
go test  -mod=vendor -v ./src/server/pfs/server -timeout 3600s | grep -v "$(date +^%FT)"
