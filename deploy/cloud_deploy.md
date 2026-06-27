# Field Day App Cloud Deployment Guide

## 1. Prepare Your Cloud Instance

- **Create a cloud instance** (Hetzner, DigitalOcean, etc.) with Fedora 44+.
- **Install required packages:**

  ```sh
  sudo dnf install -y nginx firewalld
  ```

## 2. Build and Upload the Application

- **Build the binary** (from your local machine):

  ```sh
  GOOS=linux GOARCH=amd64 go build -o fieldday main.go
  ```

- **Create service user and directories** on the server:

  ```sh
  sudo useradd -r -s /sbin/nologin nfarl
  sudo mkdir -p /opt/fieldday /var/lib/fieldday
  sudo chown -R nfarl:nfarl /opt/fieldday /var/lib/fieldday
  ```

- **Upload files** (from your local machine):

  ```sh
  scp fieldday root@<SERVER_IP>:/opt/fieldday/fieldday
  scp members-2026.csv root@<SERVER_IP>:/var/lib/fieldday/
  scp deploy/pavelanni.dev.conf root@<SERVER_IP>:/etc/nginx/conf.d/
  scp deploy/fieldday.service root@<SERVER_IP>:/etc/systemd/system/
  ```

- **Set permissions:**

  ```sh
  sudo chmod +x /opt/fieldday/fieldday
  sudo chown nfarl:nfarl /opt/fieldday/fieldday /var/lib/fieldday/members-2026.csv
  ```

## 3. Set Up SSL with Cloudflare Origin Certificate

The `pavelanni.dev` zone uses **Full (Strict)** SSL mode, so the
origin server needs a valid certificate. Cloudflare Origin
Certificates work for this — they're trusted by Cloudflare's edge
and last up to 15 years.

- **Create an Origin Certificate** (if you don't have one):
  In Cloudflare → SSL/TLS → Origin Server → Create Certificate.
  Use the wildcard `*.pavelanni.dev` so it covers all subdomains.

- **Install on the server:**

  ```sh
  sudo mkdir -p /etc/nginx/ssl
  sudo vim /etc/nginx/ssl/fullchain.pem   # paste the certificate
  sudo vim /etc/nginx/ssl/privkey.pem     # paste the private key
  sudo chmod 600 /etc/nginx/ssl/privkey.pem
  sudo chown root:root /etc/nginx/ssl/privkey.pem
  ```

  **Save the private key** — Cloudflare only shows it once.

## 4. Configure NGINX

The included `pavelanni.dev.conf` handles both HTTP (port 80) and
HTTPS (port 443) with the Origin Certificate, proxying to the app
on localhost:3000.

- **Remove default config** (optional):

  ```sh
  sudo rm -f /etc/nginx/conf.d/default.conf
  ```

- **Test and enable:**

  ```sh
  sudo nginx -t
  sudo systemctl enable --now nginx
  ```

## 5. Configure Firewall and SELinux

- **Firewall** (allow only HTTP/HTTPS):

  ```sh
  sudo systemctl enable --now firewalld
  sudo firewall-cmd --permanent --add-service=http
  sudo firewall-cmd --permanent --add-service=https
  sudo firewall-cmd --reload
  ```

- **SELinux** (allow NGINX to proxy):

  ```sh
  sudo setsebool -P httpd_can_network_connect 1
  ```

- **Verify port 3000 is closed externally:**

  ```sh
  curl --connect-timeout 5 http://<SERVER_IP>:3000
  ```

## 6. Start the Application

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now fieldday
```

Check logs: `sudo journalctl -u fieldday -f`

## 7. Cloudflare DNS

- **Set an A record** for `fieldday.pavelanni.dev` pointing to your
  server's IP.
- **Enable the orange cloud** (proxy) in Cloudflare.
- SSL/TLS mode should remain **Full (Strict)** (zone-wide setting).

## 8. (Optional) Snapshot to Save Costs

- **Shut down:** `sudo shutdown now`
- Take a snapshot in your cloud provider's dashboard.
- Delete the instance to stop billing (snapshot storage is cheap).
- Restore from snapshot before Field Day.
- **After restoring:** update the DNS A record if the IP changes.

---

## Updating for a New Year

1. Update `deploy/fieldday.service`: change the database filename
   (e.g., `fd2026.db` → `fd2027.db`) and members CSV path.
2. Rebuild and re-upload the binary.
3. Upload the new members CSV.
4. Restart: `sudo systemctl restart fieldday`
5. The Origin Certificate (`*.pavelanni.dev`) is valid until 2041 —
   no need to reissue.

## Useful Commands

| Task | Command |
|---|---|
| App logs | `journalctl -u fieldday -f` |
| NGINX status | `systemctl status nginx` |
| Firewall rules | `firewall-cmd --list-all` |
| SELinux status | `sestatus` |
| Restart app | `systemctl restart fieldday` |
