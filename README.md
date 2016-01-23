# DNS config for NSQD

Dynamically configure `nsqd` or `nsqadmin` with IP addresses from DNS A-record.

### Command Line Options

```
-lookupd-dns-address=: <addr> of DNS A record that provides nsqlookupd IP Addresses
-lookupd-tcp-port=4160: <port> of nsqlookupd TCP address
-config-http-address=: <addr>:<port> of the HTTP config endpoint
```

### Run

```
docker run --rm -it harlow/nsqd-discovery \
  --lookupd-dns-address $LOOKUPD_DNS_ADDRESS \
  --config-http-address $CONFIG_HTTP_ADDRESS
```

### Tiny Docker Image

The image is around ~5MB. Thanks to this [post from Travis Reeder](
http://www.iron.io/blog/2015/07/an-easier-way-to-create-tiny-golang-docker-images.html).

A build script has been included for convince.
