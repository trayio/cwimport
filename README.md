# cwimport

Send data from Prometheus to CloudWatch.

## Warning

Data received from Prometheus seems to be model.Vector no matter what the query is.
However this is probably wrong assumption that's why we log if data is of different type.

## Configuration

Configuration file is in [HCL](https://github.com/hashicorp/hcl), which also support JSON, but this hasn't been tested (yet).

For a sample configuration see `sample-config.hcl`. For now all fields are required (aka there are no defaults).

## Running

```
$ ./cwimport -h
Usage of ./cwimport:
  -config string
        Configuration file (default "config.hcl")
  -t    Test configuration and exit
```

Or using docker:

```
$ docker run --rm tray/cwimport:latest -h
Usage of /bin/cwimport:
  -config string
        Configuration file (default "config.hcl")
  -t    Test configuration and exit
```
