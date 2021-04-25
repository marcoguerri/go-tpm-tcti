### go-tpm-tcti
Library which implements TCTI-like layer for go-tpm library. Support abrmd only for the moment,
by implementing Reader/Writer interface and proxying all commands to abrmd via dbus.

### Usage
The library implements a `Broker` object that can be passed directly to `go-tpm/tpm2` commands.
See examples directory.

### License
The library is distributed under BSD 2-Clause license. See LICENSE file for full text.
