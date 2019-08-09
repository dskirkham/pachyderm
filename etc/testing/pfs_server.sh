# grep out PFS server logs, as otherwise the test output is too verbose to
# follow and breaks travis
go test -v ./src/server/pfs/server -timeout $TIMEOUT | grep -v "$(date +^%FT)"
