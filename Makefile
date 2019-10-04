NAME = bitfinex-api-go

#############################
### Swagger

gen-docs:
	@echo "Generating documentation"
	godocdown ./v2/websocket > ./docs/ws_v2.md
	godocdown ./v2/rest > ./docs/rest_v2.md
	godocdown ./v1/ > ./docs/v1.md
