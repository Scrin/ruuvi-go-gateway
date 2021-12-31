# ruuvi-go-gateway

ruuvi-go-gateway is a software that tries to replicate [Ruuvi Gateway](https://ruuvi.com/gateway/) MQTT and HTTP POST (for custom servers) features so that users without a physical Ruuvi Gateway can still use software and tools created for the gateway, such as [RuuviBridge](https://github.com/Scrin/RuuviBridge/).

### Features

- Supports publishing BLE data to MQTT in real time in same format as Ruuvi Gateway
- Supports sending latest BLE data via HTTP POST in same format as Ruuvi Gateway
- Can send either just Ruuvi data or all scanned BLE data (configurable, like with the Gateway)

### Requirements

- Linux-based OS (the bluetooth stack varies too greatly between operating systems and it would be simply too much work to support all of them separately)
- Bluetooth adapter supporting Bluetooth Low Energy
- 10MB of disk space
- 20MB of RAM (typical usage around 10MB)

### Configuration

Check [config.sample.yml](./config.sample.yml) for a sample config. By default the gateway assumes to find a file called `config.yml` in the current working directory, but that can be overridden with `-config /path/to/config.yml` command line flag.

### Installation

Recommended method is using Docker with the prebuilt dockerimage: [ghcr.io/scrin/ruuvi-go-gateway](https://ghcr.io/scrin/ruuvi-go-gateway) for which you can use the provided [composefile](./docker-compose.yml)

Without docker you can download prebuilt binaries from the [releases](https://github.com/Scrin/ruuvi-go-gateway/releases) page. For production use it's recommended to set up as a service.

Note that running the standalone binaries without root requires some extra capabilities be set to the binary to grant it permissions to scan for ble, this can be done with:

```sh
sudo setcap 'cap_net_raw,cap_net_admin+eip' ./ruuvi-go-gateway
```
