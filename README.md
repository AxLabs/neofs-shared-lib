# NeoFS Shared Library

This repo provides a shared library with key functionalities for interacting with NeoFS. It is built
with [cgo](https://go.dev/blog/cgo) and can be used on multiple platforms and architectures.

The purpose is to avoid re-implementation of complex and risky functions of NeoFS in other
languages, such as Java, Kotlin, Python, Typescript, etc.

In order to use this shared library, have a look at the specification of all
provided [exported methods](./EXPORTED_METHODS.md).

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
| linux/arm64   |     ✅     |   👍   | Ubuntu 22.04 |
| linux/armv5   |     ⛔     |  N/A   |     N/A      |
| linux/armv6   |     ⛔     |  N/A   |     N/A      |
| linux/armv7   |     ✅     |   🫣   |     N/A      |
| windows/amd64 |     ✅     |   👍   |  Windows 11  |
| windows/arm64 |     ⛔     |  N/A   |     N/A      |
| windows/386   |     ✅     |   🫣   |     N/A      |
| darwin/arm64  |     ✅     |   👍   |  MacOS 12.6  |
| darwin/amd64  |     ✅     |   👍   |  MacOS 13.0  |

Meaning:

* ✅: yes, it's supported, yay!
* ⛔: no release targeting the specific platform yet
* 👍: yes, manually tested, meaning that the `.so` file could successfully be loaded in the
  platform/arch
* 🫣: not tested yet...
* 👎: manual tests failed, but more tests needs to be conducted.

# References

* https://teemukanstren.com/2018/04/16/trying-to-learn-ecdsa-and-golang/
* https://github.com/elastic/golang-crossbuild
* https://github.com/freewind-demos/call-go-function-from-java-demo