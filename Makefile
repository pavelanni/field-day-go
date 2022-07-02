build:
	go build -o fieldday main.go

install:
	sudo mkdir -p /var/local/fieldday
	sudo cp fieldday deploy/start.sh /var/local/fieldday
	sudo cp -a static/ templates/ /var/local/fieldday
	sudo cp deploy/fieldday.service /lib/systemd/system

user:
	sudo deploy/mkuser.sh
	sudo mkdir -p ~nfarl/.config/lxsession/LXDE-pi/
	sudo cp deploy/autostart  ~nfarl/.config/lxsession/LXDE-pi/
	sudo chown -R nfarl ~nfarl/.config

all: build install user
