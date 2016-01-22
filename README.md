# DNS config for NSQD

Dynamically configure nsqd with IP address from DNS record.

```
docker run --rm -it harlow/nsqd-discovery \
  -lookupd-dns-address $LOOKUPD_DNS_ADDRESS \
  -nsqd-http-address $NSQD_HTTP_ADDRESS:4151
```
