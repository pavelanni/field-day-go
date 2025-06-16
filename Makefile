build:
	go build -o fieldday main.go

build-raspi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o fieldday main.go

stop:
	sudo systemctl stop fieldday || echo "Service fieldday is not installed"
start:
	sudo systemctl daemon-reload
	sudo systemctl start fieldday

install: stop
	sudo mkdir -p /opt/fieldday
	sudo mkdir -p /var/lib/fieldday
	sudo chown -R nfarl:nfarl /var/lib/fieldday
	sudo chown -R nfarl:nfarl /opt/fieldday
	sudo cp fieldday /opt/fieldday
	sudo cp deploy/fieldday.service /etc/systemd/system
	sudo systemctl daemon-reload
	sudo systemctl enable fieldday
	sudo systemctl start fieldday
	@if command -v getenforce >/dev/null 2>&1 && [ "$$(getenforce)" = "Enforcing" ]; then \
		echo "SELinux is enforcing, setting httpd_can_network_connect"; \
		sudo setsebool -P httpd_can_network_connect 1; \
	else \
		echo "SELinux is not enforcing or not installed, skipping setsebool"; \
	fi

user:
	sudo deploy/mkuser.sh
	sudo mkdir -p ~nfarl/.config/lxsession/LXDE-pi/
	sudo cp deploy/autostart  ~nfarl/.config/lxsession/LXDE-pi/
	sudo chown -R nfarl:nfarl ~nfarl/.config

all: build install user start
