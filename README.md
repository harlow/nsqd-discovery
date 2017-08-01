# DNS config for NSQD

Dynamically configure `nsqd` with IP addresses from DNS A-record. Read more about [NSQ dynamic service discovery with DNS records](https://medium.com/@harlow/nsq-service-discovery-with-dns-records-de9d759db150).

### Set up

1. Create a private zone in AWS Route 53
2. Create a Record Set with an address
3. Add an A record w/ a list of IP's for `nsqlookupd`

Note: Step 3 should be automated. When `nsqlookupd` boots the IP should be automatcially added to the above A record.

### Command Line Options

```
-config-http-address=: <addr>:<port> of the HTTP config endpoint
-lookupd-dns-address=: <addr> of DNS A record that provides nsqlookupd IP Addresses
-lookupd-tcp-port=4160: <port> of nsqlookupd TCP address
```

### Run

```
docker run --rm -it harlow/nsqd-discovery \
  --lookupd-dns-address $LOOKUPD_DNS_ADDRESS \
  --config-http-address $CONFIG_HTTP_ADDRESS
```

### Development

After verifying the desired functionality build the docker image in preparation for deployment.

The image is around ~5MB. Thanks to this [post from Travis Reeder](
http://www.iron.io/blog/2015/07/an-easier-way-to-create-tiny-golang-docker-images.html).

    make build

### Deploy

    make deploy


