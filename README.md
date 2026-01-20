# scantopl

Automatically send [scanservjs](https://github.com/sbs20/scanservjs) scanned document to [paperless-ngx](https://github.com/paperless-ngx/paperless-ngx)

## How to configure

```
Usage of /usr/bin/scantopl:
  -config string
        path to config file
  -pltoken string
        Paperless auth token , generated through admin (default "xxxxxxxxxxxxxxxxxx")
  -plurl string
        The paperless instance URL without trailing / (default "http://localhost:8080")
  -prefix string
        The file prefix to look for. Use an empty string for all files. (default "pl_")
  -scandir string
        Scanserjs ouput directory (default "/home/scanservjs/output")
```

or you can use envvar : SCANDIR, PLTOKEN, PLURL, PREFIX

provide the paperless-ngx url , the paperless-ngx token and the scanservjs output dir ( or bind to /output in docker) 

## How to use it

* Scan something
* To send it to paperless-ngx, either rename the file to include the configured prefix (default is `pl_`, e.g., `test_scan.pdf` -> `pl_test_scan.pdf`), or run `scantopl` with `-prefix ""` to upload all new files regardless of prefix.
* The file is submitted with its name (prefix and extension are automatically removed if a prefix is used). The source file is then deleted.

## How it work

* listen for file creation in the scanservjs output dir
* if a newly created file starts with the configured prefix (or any file if prefix is empty), upload it to paperless 
* If uploaded succefully, remove file from scanservjs output

## Install

### go binary

Have a working go env then

```
go install github.com/Celedhrim/scantopl@master
``` 

### Docker

```
$ docker run --rm \
  -v /your/host/scanservjs/output:/output \
  -e PLURL=https://paperless.yourdomain.instance \
  -e PLTOKEN=XXXXXXXXXXXX \
  ghcr.io/celedhrim/scantopl:master
```
