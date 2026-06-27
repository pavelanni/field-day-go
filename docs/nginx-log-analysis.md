# NGINX Log Analysis: Bot and Scanner Activity

Analysis of bot traffic hitting the Field Day server within hours of
deployment. The server runs a Go binary behind nginx on Fedora 44,
proxied through Cloudflare. None of these attacks can affect the
server — it has no PHP, Apache, Docker, or IoT firmware.

## Scan 1: 2026-06-27 (first 4 hours)

- **Server IP:** 5.161.60.14 (Hetzner)
- **Period:** 15:13 – 19:07 UTC (3h 54m)
- **Total requests:** 139 (approximately 30 legitimate)
- **Attack IPs:** 6 unique sources

### Reconnaissance scanners

| Scanner | IP | Description |
|---|---|---|
| CensysInspect | 66.132.172.189 | Internet-wide service indexing |
| Shodan | 176.65.139.66 | Internet-wide service indexing |
| Palo Alto Cortex Xpanse | 147.185.132.252 | Enterprise attack surface management |
| zgrab | 103.203.57.3, 40.124.173.6 | Open-source application scanner |
| Tencent Cloud bots | 43.153.74.75, 43.130.174.37, 43.130.57.46 | Spoofed iPhone user-agents from Chinese cloud IPs |

These are well-known internet scanners that continuously index every
public IP. They catalog services but don't attempt exploitation.

### Attack attempts

Two IPs (`167.86.80.37` and `212.127.90.201`) ran identical 30+
request attack scripts — likely the same botnet. User-agent:
`libredtail-http`. A third IP (`138.197.87.21`) ran a separate IoT
router attack.

#### 1. Apache path traversal (CVE-2021-41773 / CVE-2021-42013)

```
POST /cgi-bin/.%2e/.%2e/.%2e/.%2e/.%2e/bin/sh
POST /cgi-bin/%%32%65%%32%65/%%32%65%%32%65/.../bin/sh
```

Path traversal using URL-encoded `..` sequences to escape the web
root and reach `/bin/sh`. The double-encoded variant (`%%32%65`)
bypasses Apache's incomplete initial fix. Only affects Apache httpd
2.4.49–2.4.50.

#### 2. PHPUnit eval-stdin.php (CVE-2017-9841)

```
GET /vendor/phpunit/phpunit/src/Util/PHP/eval-stdin.php
GET /{laravel,cms,admin,blog,...}/vendor/phpunit/.../eval-stdin.php
```

Sprays ~30 common web app directory prefixes looking for PHPUnit's
`eval-stdin.php`, which executes arbitrary PHP from the POST body.
Nine years old but still the most common PHP vulnerability scan
because developers frequently expose `vendor/` directories.

#### 3. PHP-CGI argument injection (CVE-2024-4577)

```
POST /?%ADd+allow_url_include%3d1+%ADd+auto_prepend_file%3dphp://input
```

The `%AD` soft-hyphen character is interpreted as a dash by PHP-CGI
on Windows, turning query strings into CLI arguments. This enables
`allow_url_include=1` and `auto_prepend_file=php://input`, allowing
arbitrary PHP execution via the POST body. Disclosed June 2024.

#### 4. ThinkPHP RCE (CVE-2018-20062)

```
GET /index.php?s=/index/\think\app/invokefunction&function=call_user_func_array&vars[0]=md5&vars[1][]=Hello
```

Fingerprinting probe for ThinkPHP (Chinese PHP framework). The
`md5("Hello")` call is a test — if the response contains the known
hash, the bot sends the real payload. Common in Asian botnet
scanning.

#### 5. PHP pearcmd Local File Inclusion

```
GET /index.php?lang=../../../../../../../../usr/local/lib/php/pearcmd&+config-create+/&/<?echo(md5("hi"));?>+/tmp/index1.php
GET /index.php?lang=../../../../../../../../tmp/index1
```

Two-step attack: first request abuses PHP's built-in `pearcmd.php`
via Local File Inclusion to write a webshell to `/tmp/index1.php`.
Second request includes that file to execute it. This technique
gained popularity in 2022–2023 because it weaponizes a tool already
on the system without needing file upload.

#### 6. Docker API exposure

```
GET /containers/json
```

Probing for an exposed Docker daemon REST API. If found without
authentication, attackers can create privileged containers, mount
the host filesystem, and gain full host access.

#### 7. BOA web server / router login

```
POST /boaform/admin/formLogin
```

From `138.197.87.21`. Targeting BOA, a lightweight HTTP server used
in IoT devices and home routers. Credential stuffing spray hoping
to find unpatched routers or cameras.

#### 8. SAP NetWeaver (CVE-2025-31324)

```
GET /developmentserver/metadatauploader
```

Probing for a critical (CVSS 10.0) unauthenticated file upload in
SAP NetWeaver Visual Composer. This CVE is from April 2025 — only
two months old at the time of this scan, already being mass-scanned.

### Observations

- All attack traffic arrived via the **bare IP** (bypassing
  Cloudflare), not via `fieldday.pavelanni.dev`.
- The `libredtail-http` botnet runs a comprehensive multi-exploit
  script covering PHP, Apache, Docker, and Chinese frameworks in
  a single pass.
- The CVE age range spans 2017–2025, showing that botnets maintain
  a long tail of exploits — old vulnerabilities remain profitable
  because unpatched servers persist.
- The newest exploit (SAP CVE-2025-31324) was being mass-scanned
  within two months of disclosure.

---

## Scan 2: TBD (after ~24 hours)

*To be added before server shutdown.*
