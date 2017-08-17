# service usage amount api.

## build and run

```sh
make

export DATAFOUNDRY_API_SERVER=$DF_API_SERVER
export DATAFOUNDRY_API_TOKEN=$DF_TOKEN
export HADOOP_AMOUNT_BASEURL=http://$REST_API_SERVER:PORT/ocmanager/v1/api/quota
export DATAFOUNDRY_ADMIN_USER=xxxx
export DATAFOUNDRY_ADMIN_PASS=xxxx


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
-e DATAFOUNDRY_ADMIN_PASS=$DF_ADMIN_USER \
-e DATAFOUNDRY_ADMIN_USER=$DF_ADMIN_PASS \
--name svc-amount \
-d svc-amount-agent:latest 
```
