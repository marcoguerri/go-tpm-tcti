### go-tpm-tcti
Library which implements basic TCTI-like layer for [go-tpm](https://github.com/google/go-tpm). 
Supports abrmd only for the moment, by implementing Reader/Writer interface and proxying all 
commands to abrmd via dbus. Support is only minimal to what go-tpm requires: a full TCTI interface
it not relevant (e.g. cancel, polling handles support are not relevant for go-tpm, yet).

### Usage
The library implements a `Broker` object that can be passed directly to `go-tpm/tpm2` commands.
See for example [nvread.go](https://github.com/marcoguerri/go-tpm-tcti/blob/main/examples/nvread.go),
which reads EK certificate from NV index `0x1c0000a`:

```
$ go run examples/nvread.go -path /tmp/cert.der
2021/04/25 11:31:41 creating new broker connection...
2021/04/25 11:31:41 certificate written in /tmp/cert.der
$ openssl x509 -inform der -in /tmp/cert.pem -noout -text                                                                                        
Certificate:
    Data:
        Version: 3 (0x2)
[...]
```

### License
The library is distributed under BSD 2-Clause license. See LICENSE file for full text.
