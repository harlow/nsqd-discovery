# DNS config for NSQD

Dynamically configure nsqd with IP addresses from DNS record.

### Command Line Options

```
-lookupd-tcp-port=4160: <port> the TCP port to connect to nsqlookupd
-lookupd-dns-address=: <addr> to provide A record with nsqlookupd IP Addresses
-nsqd-http-address="0.0.0.0:4151": <addr>:<port> of the nsqd to configure
```

### Run

```
docker run --rm -it harlow/nsqd-discovery \
  --lookupd-dns-address $LOOKUPD_DNS_ADDRESS \
  --nsqd-http-address $NSQD_HTTP_ADDRESS:4151
```

### Tiny Docker Image

The image is around ~5MB. Thanks to this [post from Travis Reeder](
http://www.iron.io/blog/2015/07/an-easier-way-to-create-tiny-golang-docker-images.html).

A build script has been included for convince.
