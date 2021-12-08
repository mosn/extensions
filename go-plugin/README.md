## build plugin usage:

make `module` plugin=`target`

`module`:

- codec build protocol module
- trans build transcoder module
- sf build stream filter module

`target`:

- xr build xr protocol

example:

```
make codec plugin=xr
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