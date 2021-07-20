build:
	go build -o fieldday cmd/main.go

install:
	mkdir -p /var/local/fieldday
	cp fieldday deploy/start.sh /var/local/fieldday
	cp -a static/ templates/ /var/local/fieldday
	cp deploy/fieldday.service /lib/systemd/system

user:
	deploy/mkuser.sh
	mkdir -p ~registration/.config/lxsession/LXDE-pi/
	cp deploy/autostart  ~registration/.config/lxsession/LXDE-pi/
	chown -R registration ~registration/.config/lxsession/LXDE-pi/

all: build install
