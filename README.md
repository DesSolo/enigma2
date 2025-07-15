### ENIGMA
![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/dessolo/enigma2)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/dessolo/enigma2)
![Docker Pulls](https://img.shields.io/docker/pulls/dessolo/enigma2)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/dessolo/enigma2)
![GitHub](https://img.shields.io/github/license/dessolo/enigma2)

OpenSource self-hosted solution to transfer disposable links, like [onetimesecret.com](http://onetimesecret.com)

### Examples
see [deployments](/deployments/) examples

### Run container
```shell
docker run --rm -p 8080:8080 dessolo/enigma2:latest
```

### Features
- Encrypt secrets on AES256
- Scalability
- Easy to rebrand (see [Customisation](#Customisation))

#### Configuration
For more details see [example](/examples/config.yml) config.  
Default config file path is `/etc/enigma2/config.yml`. You can override it by specifying the path in the `CONFIG_FILE_PATH` environment variable.

#### Available secret storages
##### Memory
_val:_ *memory*  
Keeping all secrets in memory
> :warning: **Attention! not for productions use!!!**: old secrets are removed only on request, please choose another storages, like redis
##### Redis
_val:_ *redis*  
Keeping all secrets in [redis](https://redis.io/)
#### Build project
##### docker
```shell
docker build -t enigma .
```
#### binary
```shell
make build-server
```

#### Run compiled project
##### docker
```shell
docker run --rm --name="enigma" -p 8080:8080 localhost/enigma
```
#### binary
```shell
./enigma2_server_{VERSION}_{OS}_{ARCH}
```

#### Customisation
You can change the [html templates](/templates) to suit your needs.  
[index.html](/templates/index.html) - _main page_  
[view_secret.html](/templates/view_secret.html) - _rendered when a secret is requested_  
> :warning: **get.html must contain "%s" to insert the secret into the template**