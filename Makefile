build:
	/usr/local/go/bin/go build -o fieldday main.go

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
	mkdir -p /var/local/fieldday
	cp fieldday deploy/start.sh /var/local/fieldday
	cp -a static/ templates/ /var/local/fieldday
	cp deploy/fieldday.service /lib/systemd/system
	systemctl daemon-reload
	systemctl start fieldday

user:
	deploy/mkuser.sh
	mkdir -p ~nfarl/.config/lxsession/LXDE-pi/
	cp deploy/autostart  ~nfarl/.config/lxsession/LXDE-pi/
	chown -R nfarl ~nfarl/.config

all: build install user start
