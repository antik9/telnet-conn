## Telnet Analogue

### Get a Repo
```bash
>>> go get -u github.com/antik9/telnet-conn
```

### Run client
```bash
>>> telnet-conn --host localhost --port 80 --timeout 15 --request_timeout 2
GET / HTTP/1.0
...
```

### Parameters
 - host(as ip-address or url)
 - port(as integer)
 - timeout(as integer, limit the time of work with client)
 - request_timeout(as integer, limit the time of each request)
