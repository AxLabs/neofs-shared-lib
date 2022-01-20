# neofs-api-shared-lib

This repo provides a shared library with key functionalities from 
[neofs-api-go](https://github.com/nspcc-dev/neofs-api-go) for multiple platforms and architectures.

Below you can find key functions exported by this shared library:

- [`SignServiceMessage(key *ecdsa.PrivateKey, msg interface{})`](https://github.com/nspcc-dev/neofs-api-go/blob/master/signature/sign.go#L147)
- [`VerifyServiceMessage(msg interface{})`](https://github.com/nspcc-dev/neofs-api-go/blob/master/signature/sign.go#L227)

# Build

Make sure you have [docker](https://docker.com) installed and execute:

```shell
sh cross.sh
```

# Supported Platforms and Archs

| Platform/Arch     | Supported |
| -------------     | --------- |
| linux/amd64       | ✅ |
| linux/i386        | ⛔ |
| linux/arm64       | ✅ |
| linux/armv5       | ⛔ |
| linux/armv6       | ⛔ |
| linux/armv7       | ⛔ |
| windows/amd64     | ✅ |
| windows/386       | ✅ |
| darwin/arm64      | ✅ |
| darwin/amd64      | ✅ |

# References

* https://teemukanstren.com/2018/04/16/trying-to-learn-ecdsa-and-golang/
* https://github.com/elastic/golang-crossbuild
* https://github.com/freewind-demos/call-go-function-from-java-demo