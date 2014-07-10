
BINS=web query socket_client query_xml hello

GO=go build

all:	$(BINS)

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

clean:
	rm -f $(BINS)
