
BINS=simpleWeb web query socket_client query_xml hello goRoutines

GO=go build

all:	$(BINS)

simpleWeb:	simpleWeb.go ./src/apc_ups/apc_ups.go
	$(GO) simpleWeb.go

web:	web.go
	$(GO) web.go

query:	query.go
	$(GO) query.go

socket_client:	socket_client.go
	$(GO) socket_client.go

query_xml:	query_xml.go
	$(GO) query_xml.go

hello:	hello.go
	$(GO) hello.go

goRoutines:	goRoutines.go
	$(GO) goRoutines.go

clean:
	rm -f $(BINS) *~
