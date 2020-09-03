### ENIGMA
Simple online solution to transfer disposable links, like [onetimesecret.com](onetimesecret.com)

#### Environments
|Env variable|Default value|Description|
|---|---|---|
|LISTEN_PORT|9000|_server listen port_|
|TOKEN_BYTES|20|_lenght of url secret hash_|
|UNIQ_KEY_RETRIES|3|_how many times to try to keep a secret in storage_|
|RESPONSE_ADDRESS|http://127.0.0.1:9000|_the server returns the full path to seret, specify the address_|
|SECRET_STORAGE|Memory|_where to keep secrets_|

#### Avalible secret storages
Memory

#### Build image
##### docker
```shell
docker build -t enigma2 .
```
##### podman
```shell
podman build -t enigma2 .
```

### Build binary
```shell
make build
ls bin
```