# NeoFS API Shared Lib

This repo provides a shared library with key functionalities from 
[neofs-api-go](https://github.com/nspcc-dev/neofs-api-go) for multiple platforms and architectures.

The purpose is to avoid re-implementation of complex and risky functions of NeoFS in other languages, 
such as Java, Kotlin, Python, Typescript, etc.

Below you can find key functions exported by this shared library:

- [`SignServiceMessage(key *ecdsa.PrivateKey, msg interface{})`](https://github.com/nspcc-dev/neofs-api-go/blob/master/signature/sign.go#L147)
- [`VerifyServiceMessage(msg interface{})`](https://github.com/nspcc-dev/neofs-api-go/blob/master/signature/sign.go#L227)

# Build

Make sure you have [docker](https://docker.com) installed.

Then, execute:

```shell
bash cross.sh
```

The output, i.e., `.h` and `.so` files, will be placed in the `./libs` folder.

# Supported Platforms and Archs

| Platform/Arch     | Supported |
| :-------------    | :---------: |
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