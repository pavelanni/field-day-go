# Field Day App Cloud Deployment Guide

## 1. Prepare Your Cloud Instance

- **Create a DigitalOcean droplet** (Fedora 41 recommended).
- **Update system:**

  ```sh
  sudo dnf update
  ```

## 2. Install Required Software

- **NGINX:**

  ```sh
  sudo dnf install nginx
  sudo systemctl enable --now nginx
  ```

- **Your Application:**
  - Deploy your app to run on `localhost:3000` (use `systemd`, `pm2`, or Docker as needed).

## 3. Configure NGINX

- **Copy your NGINX config** (e.g., `pavelanni.dev.conf`) to `/etc/nginx/conf.d/`.
- **Remove default config** (optional):

  ```sh
  sudo rm /etc/nginx/conf.d/default.conf
  ```

- **Test and reload:**

  ```sh
  sudo nginx -t
  sudo systemctl reload nginx
  ```

## 4. Set Up SSL with Cloudflare Origin Certificate

- **In Cloudflare dashboard:**
  Go to SSL/TLS → Origin Server → Create Certificate.
- **On your server:**

  ```sh
  sudo mkdir -p /etc/nginx/ssl
  sudo vim /etc/nginx/ssl/fullchain.pem
  sudo vim /etc/nginx/ssl/privkey.pem
  ```

  Paste the certificate and key.

- **Set permissions:**

  ```sh
  sudo chmod 600 /etc/nginx/ssl/privkey.pem
  sudo chown root:root /etc/nginx/ssl/privkey.pem
  ```

## 5. SELinux Configuration (Fedora)

- **Allow NGINX to connect to your app:**

  ```sh
  sudo setsebool -P httpd_can_network_connect 1
  ```

## 6. Enable and Configure firewalld

- **Install and start:**

  ```sh
  sudo dnf install firewalld
  sudo systemctl enable --now firewalld
  ```

- **Allow only HTTP/HTTPS:**
  ```sh
  sudo firewall-cmd --permanent --add-service=http
  sudo firewall-cmd --permanent --add-service=https
  sudo firewall-cmd --reload
  ```

- **Verify port 3000 is closed externally:**

  ```sh
  curl http://your-server-ip:3000
  nmap -p 3000 your-server-ip
  ```

## 7. Cloudflare DNS and Proxy

- **Set an A record** for your domain (e.g., `fieldday.pavelanni.dev`) to your droplet’s IP.
- **Enable the orange cloud** (proxy) in Cloudflare.
- **Set SSL/TLS mode** to **Full (strict)**.

## 8. (Optional) Take a Snapshot to Save Costs

- **Shut down your droplet:**

  ```sh
  sudo shutdown now
  ```

- **In DigitalOcean dashboard:**
  Go to Droplet → Snapshots → Take Snapshot.
- **Delete the droplet** to stop billing for compute (you’ll only pay for snapshot storage).
- **Restore from snapshot** before Field Day.

## 9. After Restoring from Snapshot

- **Update DNS** if your droplet’s IP changes.
- **Test your app and NGINX.**

---

## Troubleshooting

- **Permission denied connecting to upstream:**
  Run: `sudo setsebool -P httpd_can_network_connect 1`
- **Port 3000 accessible from outside:**
  Check `firewalld` rules and ensure only HTTP/HTTPS are open.

---

## Useful Commands

- **Check NGINX status:**
  `sudo systemctl status nginx`
- **Check firewall rules:**
  `sudo firewall-cmd --list-all`
- **Check SELinux status:**
  `sestatus`

---
