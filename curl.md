```
curl --unix-socket /var/run/docker.sock http:/info
```

```
curl --unix-socket /var/run/docker.sock http:/images/json
```

```
curl --unix-socket /var/run/docker.sock http:/containers/json
```

```
curl --unix-socket /var/run/docker.sock \
     -X POST http:/containers/create \
     -H "Content-Type: application/json" \
     -d '{"Image":"alpine", "Cmd":["/bash/bin"]}' \
     -v
```




```
curl http://localhost:9090/v1/docker/images/json
```

```
curl -X POST http://localhost:9090/v1/docker/containers/create \
     -H "Content-Type: application/json" \
     -d '{"Image":"alpine", "Cmd":["/bash/bin"]}' \
     -v
```
