# scantopl

Automatically send [scanservjs](https://github.com/sbs20/scanservjs) scanned document to [paperless-ng](https://github.com/jonaswinkler/paperless-ng)

## How to configure

```
Usage of /usr/bin/scantopl:
  -config string
        path to config file
  -pltoken string
        Paperless auth token , generated through admin (default "xxxxxxxxxxxxxxxxxx")
  -plurl string
        The paperless instance URL without trailing / (default "http://localhost:8080")
  -scandir string
        Scanserjs ouput directory (default "/home/scanservjs/output")
```

or you can use envvar : SCANDIR, PLTOKEN, PLURL

provide the paperless-ng url , the paperless-ng token and the scanservjs output dir ( or bind to /output in docker) 

## How to use it

* Scan something
* if you want to send it to paperless-ng , go in the scanservjs file section and rename file to add prefix **pl_** ( test_scan.pdf -> pl_test_scan.pdf)
* the file is submitted with name "test_scan" ( remove prefix and extension automatically) then remove source file is deleted 

## How it work

* listen for file creation in the scanservjs output dir
* if a newly created file start with **pl_** , upload it to paperless 
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
