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

# Platforms and Archs

| Platform/Arch | Supported | Tested |  OS Tested   |
|:--------------|:---------:|:------:|:------------:|
| linux/amd64   |     ✅     |   👍   | Ubuntu 22.04 |
| linux/arm64   |     ✅     |   👎   | Ubuntu 22.04 |
| linux/armv5   |     ⛔     |  N/A   |     N/A      |
| linux/armv6   |     ⛔     |  N/A   |     N/A      |
| linux/armv7   |     ⛔     |  N/A   |     N/A      |
| windows/amd64 |     ✅     |   👍   |  Windows 11  |
| windows/arm64 |     ✅     |   👎   |  Windows 11  |
| windows/386   |     ✅     |   🫣   |     N/A      |
| darwin/arm64  |     ✅     |   👍   |  MacOS 12.6  |
| darwin/amd64  |     ✅     |   👍   |  MacOS 13.0  |

Meaning:
* ✅: yes, it's supported, yay!
* ⛔: no release targeting the specific platform yet
* 👍: yes, manually tested, meaning that the `.so` file could successfully be loaded in the platform/arch
* 🫣: not tested yet...
* 👎: manual tests failed, but more tests needs to be conducted.

# References

* https://teemukanstren.com/2018/04/16/trying-to-learn-ecdsa-and-golang/
* https://github.com/elastic/golang-crossbuild
* https://github.com/freewind-demos/call-go-function-from-java-demo