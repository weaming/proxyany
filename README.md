![](logo.png)

## Feature

- Reverse proxy one site to any host, any port
- Built-in HTTPS certification from let's encrypt (force 443 port)
- Rewrite response headers and text body, replace it from the target domain to your domain

## Usage

```
Usage of proxyany:
  -from string
    	your reverse proxy target url, including path is allowed, then your visit path will be appended to it (default "https://www.google.com")
  -to string
    	local bind [<host>]:<port> (default ":20443")
  -https
    	HTTPS mode, auto certification from let's encrypt
  -domain string
    	domain allowed to access, all sub domains will be allowed too (default "bitsflow.org")
```

```sh
# proxy github
$ proxyany -from https://github.com/weaming/ -to :20443

# proxy google
$ proxyany -from https://www.google.com -https

# proxy https://static.rust-lang.org to speed up downloading rust installer
$ proxyany -from https://static.rust-lang.org -https -domain rustup.bitsflow.org
```

You could run `proxyany -from https://<domain> -https -domain <domain>` first,
then modify your hosts records to point the domain to your server,
so that you could visit the site bypassing the firewall.
Your server only acts as an transparent proxy.
