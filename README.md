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

Or you could run it in a Docker container, as follows:
```bash
ip addr add 192.168.1.2/24 dev br0
docker run -d -it \
        --mount type=bind,source=/home/user/go/bin/uploader,target=/opt/uploader \
        --mount type=bind,source=/srv/public-share/,target=/tmp/public-share/ \
        -p 192.168.1.2:5050:5050 \
        -e "DESTDIR=/var/uploads/" \
        -e "BASEURL=https://example.net/files/" \
        -e "UPLOADURL=https://example.net/u/uploads" \
        ubuntu:latest \
        /opt/uploader
```

Note that Golang will parse uploaded files larger than 100MB into files stored
in /tmp by default, so the size of `/tmp` limits that max size of files that can
be uploaded. Use an alternate `/tmp` with more disk space if needed, by setting
the `TMPDIR` env var.
