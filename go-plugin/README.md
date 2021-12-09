## build plugin usage:

make `module` plugin=`target`

`module`:

- codec build protocol module
- trans build transcoder module
- filter build stream filter module

`target`:

- xr build xr protocol

example:

``` shell
# build codec only
make codec plugin=xr

# build code and copy filter
# make filter plugin=auth (build stream filter auth plugin)
make codec plugin=xr filter=auth

# build code and copy transcoder
# make trans plugin=xr2sp (build trancoder plugin)
make codec plugin=xr transcoder=xr2sp
``` 

## build sidecar:

```shell
mkdir -p ~/go/src/mosn.io
git clone https://github.com/YIDWang/mosn.git
cd mosn
git checkout test-test

cd ~/go/src/github.com/mosn/extensions/go-plugin
make sidecar
```

## start and stop sidecar

```shell
cd ./etc/script
sh start.sh module target

# eg:
sh start.sh codec xr

# eg:
sh stop.sh 
```