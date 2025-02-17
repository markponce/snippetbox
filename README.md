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