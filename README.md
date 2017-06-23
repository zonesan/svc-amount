# service usage amount api.

## build and run

```sh
make

export DATAFOUNDRY_API_SERVER=<dfapiserver>
export DATAFOUNDRY_API_TOKEN=<dftoken>
export HADOOP_AMOUNT_BASEURL=<http://192.168.1.113:8080/ocmanager/v1/api/quota/>

./bin/linux/svc-amount
```

## build docker image and run docker images

```sh
make images

docker run \
-p 8080:8080 \
-e DATAFOUNDRY_API_SERVER=$DATAFOUNDRY_API_SERVER \
-e DATAFOUNDRY_API_TOKEN=$DATAFOUNDRY_API_TOKEN \
-e HADOOP_AMOUNT_BASEURL=$HADOOP_AMOUNT_BASEURL \
--name svc-amount \
-d svc-amount-agent:latest 
```