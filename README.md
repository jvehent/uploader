# A dummy upload server

Run with:
```bash
$   DESTDIR=/var/uploads/ \
    BASEURL=https://example.net/files/ \
    UPLOADURL=https://example.net/u/uploads \
    ./upload_server
```

Additionally, use an nginx config like:
```
    location ~/u(?<uri>.+) {
        proxy_pass http://localhost:5050$uri;
    }
```
which will expose the html form at https://example.net/u/
