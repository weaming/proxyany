![](logo.png)

## Feature

- Reverse proxy one site to any host, any port
- Built-in HTTPS certification from let's encrypt (force 443 port)
- Rewrite request headers
- Rewrite response headers and text body

## Usage

```
Usage of proxyany:
  -config string
    	file path domain mapping config in json format (default "config.json")
  -bind string
    	local bind [<host>]:<port> (default ":20443")
  -https
    	HTTPS mode, auto certification from let's encrypt
```

## Example config

```sh
$ proxyany -https
```

### proxy google

```json
[
    {"from": "google.bitsflow.org", "to": "https://google.com"}
]
```

### proxy twitter

```json
[
    {"from": "t.byteio.cn", "to": "https://twitter.com"},
    {"from": "img.byteio.cn", "to": "https://twimg.com"}
]
```
