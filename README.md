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
- **Secure**:
    - Secrets are encrypted using AES-256.
    - Secrets are deleted immediately after being viewed.
    - An optional confirmation step prevents accidental secret exposure.
- **Flexible Configuration**:
    - Configure the application with a single `config.yml` file.
    - Pluggable storage backends: use in-memory for testing or Redis for production.
    - Set a lifetime for secrets.
    - Customize secret token length.
- **Easy to Deploy and Scale**:
    - Ready-to-use Docker image and Docker Compose example.
    - Stateless server design allows for horizontal scaling.
    - Healthcheck endpoint for service monitoring.
- **Developer Friendly**:
    - Easily rebrand and customize HTML templates (see [Customisation](#Customisation)).
    - Structured request logging with request IDs for easier debugging.
    - Graceful shutdown to prevent data loss.

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
