build:
	go build -o fieldday main.go

build-raspi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o fieldday main.go

stop:
        if systemctl list-units | grep fieldday; then \
            systemctl stop fieldday; \
        else \
            echo "Service fieldday is not installed"; \
        fi
start:
	systemctl daemon-reload
	systemctl start fieldday

install: stop
	mkdir -p /opt/fieldday
	mkdir -p /var/lib/fieldday
	chown -R nfarl:nfarl /var/lib/fieldday
	chown -R nfarl:nfarl /opt/fieldday
	cp fieldday /opt/fieldday
	cp deploy/fieldday.service /etc/systemd/system
	systemctl daemon-reload
	systemctl start fieldday

user:
	deploy/mkuser.sh
	mkdir -p ~nfarl/.config/lxsession/LXDE-pi/
	cp deploy/autostart  ~nfarl/.config/lxsession/LXDE-pi/
	chown -R nfarl ~nfarl/.config

all: build install user start
