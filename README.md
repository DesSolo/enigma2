### ENIGMA
Simple online solution to transfer disposable links, like [onetimesecret.com](onetimesecret.com)

### Run example
```shell
docker run --rm -p 9000:9000 dessolo/enigma2:latest
```

#### Environments variables
|Env variable|Default value|Description|
|---|---|---|
|LISTEN_PORT|9000|_server listen port_|
|TOKEN_BYTES|20|_lenght of url secret hash_|
|UNIQ_KEY_RETRIES|3|_how many times to try to keep a secret in storage_|
|RESPONSE_ADDRESS|http://127.0.0.1:9000|_the server returns the full path to seret, specify the address_|
|SECRET_STORAGE|Memory|_where to keep secrets_|
|Redis storage|
|REDIS_ADDRESS|localhost:6379|_redis server address_|
|REDIS_PASSWORD|_empty_|_redis server password_|
|REDIS_DATABASE|0|_redis server database_|

#### Avalible secret storages
##### Memory
_val_ *Memory*  
Keeping all secrets in memory
> :warning: **Attention! not for productions use!!!**: Use another storages, like redis
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