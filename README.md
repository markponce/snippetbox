# snippetbox

# run command 
go run ./cmd/web -addr=":80"
go run ./cmd/web -help
go run ./cmd/web >>/tmp/web.log

# go mod 

go mod init <package name or repo link>
go mod download
go mod verify
# automatically remove any unused packages from your go.mod and go.sum files
go mod tidy

# latest version of package
go get -u github.com/foo/bar
# Specific version
go get -u github.com/foo/bar@v2.0.0
# Remove package
go get github.com/foo/bar@none

# curl command
# get 
curl -i localhost:4000/
# head only
curl --head localhost:4000/
# post
curl -i -d "" localhost:4000/snippet/create
# delete
curl -i -X DELETE localhost:4000/snippet/create

# mysql

# show charset and collation
SELECT SCHEMA_NAME, DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME 
FROM information_schema.SCHEMATA 
WHERE SCHEMA_NAME = 'snippetbox';

# update database charset and collation
ALTER DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


# Create table sessions 
USE snippetbox;

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

# generate cert
mkdir tls
cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost


# users table 
USE snippetbox;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

# Build 
$ go build -o /tmp/web ./cmd/web/
$ cp -r ./tls /tmp/
$ cd /tmp/
$ ./web 

# Test 
go test -v ./cmd/web

# Run test current project per package?
go test ./...

# Run specific test with regex
go test -v -run="^TestPing$" ./cmd/web/

# Run test Format {test regexp}/{sub-test regexp}
go test -v -run="^TestHumanDate$/^UTC$" ./cmd/web

# Run test but skip regex
go test -v -skip="^TestHumanDate$" ./cmd/web/

# Run test not cached
go test -count=1 ./cmd/web 

# Clean cache test result
go clean -testcache

# Stop test when test fails
go test -failfast ./cmd/web

# Test concurrent
Must put t.Parallel()
```go
func TestPing(t *testing.T) {
    t.Parallel()

    ...
}
```
go test -parallel=4 ./...

# Test with race condition detector
go test -race ./cmd/web/

# go run specific test
go test -v -run="<function name>" ./cmd/web
eg.:
go test -v -run="TestUserSignup" ./cmd/web

# short
go test -v -short ./...

go test -cover ./...

# Profiling test coverage
go test -coverprofile=/tmp/profile.out ./...
go tool cover -func=/tmp/profile.out

go test -covermode=count -coverprofile=/tmp/profile.out ./...
go tool cover -func=/tmp/profile.out