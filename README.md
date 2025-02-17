# snippetbox

# run command 
go run ./cmd/web -addr=":80"
go run ./cmd/web -help
go run ./cmd/web >>/tmp/web.log

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
