
BINS=simpleWeb web query socket_client query_xml hello goRoutines libtest.so

GO=go build
CC=gcc -g

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

goRoutines:	goRoutines.go libtest.so
	$(GO) goRoutines.go

test.o:	test.c test.h Makefile
	$(CC) -fPIC -c test.c -o test.o

libtest.so:	test.o
	$(CC) -shared -Wl,-soname,libtest.so.1 -o libtest.so.1.0 test.o
	ln -sf libtest.so.1.0 libtest.so.1
	ln -sf libtest.so.1 libtest.so

tester.o:	tester.c test.h
	$(CC) -c tester.c -o tester.o

tester:	tester.o test.h libtest.so
	$(CC) tester.o -L. -ltest -o tester
clean:
	rm -f $(BINS) *~
