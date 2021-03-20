Prometheus IPMI Exporter
===

Requirements
---
* Go compiler
* `ipmitool` package

Building
---

```
$ make
mkdir -p build/{linux,solaris}
GOOS=linux go build -o build/linux/ipmi_exporter main.go
GOOS=solaris go build -o build/solaris/ipmi_exporter main.go
```

Usage
---

```
Usage of ipmi_exporter:
  -ipmi-passwd string
        [IPMI_PASSWD] IPMI password
  -ipmi-target string
        [IPMI_TARGET] IPMI target address
  -ipmi-user string
        [IPMI_USER] IPMI username (default "admin")
  -listen-address string
        [LISTEN_ADDRESS] Address on which to expose metrics and web interface. (default "0.0.0.0:9100")
```

Run directly:
```
export IPMI_TARGET=10.0.16.112
export IPMI_USER=admin
export IPMI_PAASWD=****
./ipmi_exporter
```

As a SMF service:
```
# cp build/solaris/ipmi_exporter /opt/custom/smf/bin/ipmi_exporter
# svccfg import ipmi_exporter.xml
# svccfg -s ipmi-exporter
svc:/ipmi-exporter> setprop ipmi/target="10.0.16.112"
svc:/ipmi-exporter> setprop ipmi/user="admin"
svc:/ipmi-exporter> setprop ipmi/passwd="****"
svc:/ipmi-exporter> end
# svccfg -s ipmi-exporter:default refresh
# svcadm restart ipmi-exporter
```
