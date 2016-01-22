# DNS config for NSQD

Dynamically configure nsqd with IP addresses from DNS record.

### Command Line Options

```
-lookupd-dns-address=: <addr> to provide A record with nsqlookupd IP Addresses
-nsqd-http-address="0.0.0.0:4151": <addr>:<port> of the nsqd to configure
```

### Run

```
docker run --rm -it harlow/nsqd-discovery \
  --lookupd-dns-address $LOOKUPD_DNS_ADDRESS \
  --nsqd-http-address $NSQD_HTTP_ADDRESS:4151
```
