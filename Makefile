build:
	go build -o ruuvi-go-gateway cmd/ruuvi-go-gateway/main.go

install: build
	install -m 755 ruuvi-go-gateway /usr/bin/
	install -d /etc/ruuvi-go-gateway

systemd-install:
	install -m 644 service/ruuvi-go-gateway.service /lib/systemd/system/ruuvi-go-gateway.service
	systemctl enable ruuvi-go-gateway.service
	@echo "Please install a configuration file to /etc/ruuvi-go-gateway/config.yml, because"
	@echo "the systemd service uses it from there. You can start the service manually with:"
	@echo "    systemctl start ruuvi-go-gateway.service"

clean:
	rm -f ruuvi-go-gateway

all: build install systemd-install
