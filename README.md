### ENIGMA
![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/dessolo/enigma2)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/dessolo/enigma2)
![Docker Pulls](https://img.shields.io/docker/pulls/dessolo/enigma2)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/dessolo/enigma2)
![GitHub](https://img.shields.io/github/license/dessolo/enigma2)

Simple online solution to transfer disposable links, like [onetimesecret.com](onetimesecret.com)

### Run example
```shell
docker run --rm -p 9000:9000 dessolo/enigma2:latest
```

#### Environments variables
##### Global
|Env variable|Default value|Description|
|---|---|---|
|LISTEN_PORT|9000|_server listen port_|
|TOKEN_BYTES|20|_lenght of url secret hash_|
|UNIQ_KEY_RETRIES|3|_how many times to try to keep a secret in storage_|
|RESPONSE_ADDRESS|http://127.0.0.1:9000|_the server returns the full path to seret, specify the address_|
|SECRET_STORAGE|Memory|_where to keep secrets_ (see [avalible secret storages](https://github.com/DesSolo/enigma2#avalible-secret-storages))|
##### Redis storage
|Env variable|Default value|Description|
|---|---|---|
|REDIS_ADDRESS|localhost:6379|_redis server address_|
|REDIS_PASSWORD|_empty_|_redis server password_|
|REDIS_DATABASE|0|_redis server database_|

#### Avalible secret storages
##### Memory
_val_ *Memory*  
Keeping all secrets in memory
> :warning: **Attention! not for productions use!!!**: old secrets are removed only on request, please choose another storages, like redis
##### Redis
_val_ *Redis*  
> Env: REDIS_ADDRESS, REDIS_PASSWORD, REDIS_DATABASE,
#### Build project
##### docker
```shell
docker build -t enigma .
```
##### podman
```shell
podman build -t enigma .
```
#### binary
```shell
make build
```

#### Run compiled project
##### docker
```shell
docker run --rm --name="enigma" -p 9000:9000 localhost/enigma
```
##### podman
```shell
podman run --rm --name="enigma" -p 9000:9000 localhost/enigma
```
#### binary
```shell
./enigma
```

#### Customisation
You can change the [html templates](https://github.com/DesSolo/enigma2/tree/master/templates) to suit your needs.  
[index.html](https://github.com/DesSolo/enigma2/blob/master/templates/index.html) - _main page_  
[get.html](https://github.com/DesSolo/enigma2/blob/master/templates/get.html) - _rendered when a secret is requested_  
> :warning: **get.html must contain "%s" to insert the secret into the template**