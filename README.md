# gcs-uploader

Tiny uploader for Google Cloud Storage

## install

Install from https://github.com/myoan/gcs-uploader/releases

## usage

```
❯ ./gcs-client -help
Usage of ./gcs-client:
  -b string
        backetname
  -conc int
        upload cuncurrency (default 4)
  -cred string
        credential path
  -in string
        input dir path
  -out string
        output dir path
```

example
```
❯ gcs-uploader -b my-bucket-name -cred credential.json -in path/to/input -out path/to/output
---------------------------------------------
bucket:     my-bucket-name
credential: credential.json
input:      path/to/input
output:     path/to/output
---------------------------------------------

upload: gs://my-bucket-name/path/to/output/foo.txt
upload: gs://my-bucket-name/path/to/output/foo/bar.txt
```