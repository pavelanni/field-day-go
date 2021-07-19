build:
	go build -o fieldday cmd/main.go

install:
	mkdir -p /var/local/fieldday
	cp fieldday deploy/start.sh /var/local/fieldday
	cp -a static/ templates/ /var/local/fieldday
	cp deploy/fieldday.service /lib/systemd/system

all: build install
