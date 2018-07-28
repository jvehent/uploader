# A dummy upload server

Run with:
```bash
$   DESTDIR=/var/uploads/ \
    BASEURL=https://example.net/files/ \
    UPLOADURL=https://example.net/u/uploads \
    ./uploader
```

Additionally, use an nginx config like:
```
    location ~ /u(?<qs>.+) {
        auth_basic "Uploader credentials please";
        auth_basic_user_file /etc/nginx/htpasswd;
        proxy_pass http://127.0.0.1:5050$qs;
        proxy_read_timeout 1200;
    }
```
which will expose the html form at https://example.net/u/
