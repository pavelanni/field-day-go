# WiFi Access Point Setup (OrangePi 3)

The OrangePi 3 runs a WiFi access point so that devices (phones, ESP32
displays) can connect to the Field Day registration app without
any external network infrastructure.

## Quick Reference: Changing SSID and Password

Edit `/etc/hostapd/hostapd.conf` on the OrangePi:

```
ssid=FieldDay
wpa_passphrase=Nfarl-2026
```

Change these two values, then restart:

```sh
sudo systemctl restart hostapd
```

## How It Works

Four services cooperate to provide the access point:

| Service | Config File | Role |
|---------|-------------|------|
| hostapd | `/etc/hostapd/hostapd.conf` | Broadcasts the WiFi network (SSID, WPA2) |
| dnsmasq | `/etc/dnsmasq.conf` | Assigns IP addresses (DHCP) and forwards DNS |
| interfaces | `/etc/network/interfaces.d/orangepi.ap.nat` | Gives wlan0 its static IP (172.24.1.1) |
| nftables | managed by iptables-nft | NAT masquerade from wlan0 to ethernet (end0) |

NetworkManager is told to ignore wlan0 in
`/etc/NetworkManager/conf.d/10-ignore-interfaces.conf`
so it doesn't interfere with hostapd.

## Network Layout

```
WiFi clients ─── wlan0 (172.24.1.1) ─── OrangePi ─── end0 (DHCP) ─── upstream network
                  AP: FieldDay                         ethernet
                  DHCP: 172.24.1.50-150
```

- WiFi clients get addresses in 172.24.1.50 -- 172.24.1.150 (12-hour lease)
- The Field Day app is reachable at `http://fieldday.lan:3000/new`
  (or `http://172.24.1.1:3000/new` as fallback)
- If ethernet is connected, clients also get internet access via NAT
- IP forwarding is enabled in `/etc/sysctl.d/90-ip-forward.conf`

## Local DNS

dnsmasq serves `fieldday.lan` pointing to 172.24.1.1 so that WiFi
clients can use `http://fieldday.lan:3000` instead of the bare IP.
This is configured in `/etc/dnsmasq.conf`:

```
address=/fieldday.lan/172.24.1.1
```

To change the hostname, edit that line and restart dnsmasq:

```sh
sudo systemctl restart dnsmasq
```

**Why `.lan` and not `.local` or `.internal`?** The `.local` TLD is
reserved for mDNS (Avahi/Bonjour) and causes lookup delays on some
devices. The `.internal` TLD is an IANA reserved name that some OSes
handle specially. The `.lan` suffix is a common convention for
private networks served by dnsmasq and has no protocol conflicts.

## Port 80 Redirect

An nftables rule redirects port 80 to the app's port 3000 for WiFi
clients, so mobile users can visit `http://fieldday.lan` without
specifying a port. The kiosk browser is hardcoded to port 3000 and
is unaffected (the redirect only applies to traffic arriving on
wlan0).

This rule is persisted in `/etc/iptables.ipv4.nat` and restored at
boot by the `orangepi-restore-iptables` systemd service.

## Unisoc WiFi Driver Notes

The OrangePi 3 uses a Unisoc (Spreadtrum) WiFi chip with the
`sprdwl_ng` kernel module. This driver does **not** support several
standard hostapd settings. The following must stay commented out or
removed from `hostapd.conf`:

- `rts_threshold` -- causes "Could not set RTS threshold" error
- `fragm_threshold` -- same class of unsupported ioctl
- `beacon_int`, `dtim_period` -- cause "Failed to set beacon parameters"
- WMM parameters (`wmm_ac_*`) -- not needed and may cause issues

If hostapd fails to start repeatedly, the Unisoc firmware can enter
a bad state (`sprdwl_uninit_fw failed!` in dmesg). The only recovery
is a full reboot -- the SDIO bus cannot be reset by module
unload/reload alone.

## Full hostapd.conf

For reference, the working minimal config:

```ini
interface=wlan0
driver=nl80211
ssid=FieldDay
hw_mode=g
channel=6
wmm_enabled=0
auth_algs=1
macaddr_acl=0
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=Nfarl-2026
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
```

## Troubleshooting

**hostapd won't start after config change**: Check logs with
`sudo journalctl -u hostapd -n 20`. If you see firmware errors in
`sudo dmesg | grep unisoc`, reboot the OrangePi.

**Clients connect but get no IP**: Check that dnsmasq is running
(`systemctl status dnsmasq`) and that wlan0 has 172.24.1.1
(`ip addr show wlan0`).

**No internet through the AP**: Verify IP forwarding is on
(`cat /proc/sys/net/ipv4/ip_forward` should be `1`) and that ethernet
has a default route (`ip route show default`).
