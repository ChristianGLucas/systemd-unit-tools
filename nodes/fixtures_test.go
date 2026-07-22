package nodes_test

// Shared fixtures used across this package's tests. serviceFixture,
// timerFixture and socketFixture are modeled closely on real, widely
// deployed unit files (Debian/Ubuntu's nginx.service, a dnf/yum-style
// makecache.timer, and a Docker-style API socket) so that expected values
// in the tests are hand-derived from realistic unit-file text, not from
// running the code under test and capturing whatever it happens to produce.

const serviceFixture = `[Unit]
Description=A high performance web server and reverse proxy server
Documentation=man:nginx(8) https://nginx.org/en/docs/
After=network.target remote-fs.target nss-lookup.target
Wants=network-online.target
After=network-online.target

[Service]
Type=forking
PIDFile=/run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t -q -g "daemon on; master_process on;"
ExecStart=/usr/sbin/nginx -g "daemon on; master_process on;"
ExecReload=/usr/sbin/nginx -g "daemon on; master_process on;" -s reload
ExecStop=-/sbin/start-stop-daemon --quiet --stop --retry QUIT/5 --pidfile /run/nginx.pid
TimeoutStopSec=5
KillMode=mixed
Restart=on-failure
RestartSec=5s
User=www-data
Group=www-data
WorkingDirectory=/var/www
Environment=NGINX_ENV=production DEBUG=0
Environment="EXTRA=hello world"
EnvironmentFile=-/etc/default/nginx
NoNewPrivileges=true
ProtectSystem=full
PrivateTmp=true
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`

const timerFixture = `[Unit]
Description=dnf makecache timer

[Timer]
OnBootSec=10min
OnUnitActiveSec=1h
OnCalendar=*-*-* 6,18:00
OnCalendar=*-*-* 12:00
Persistent=true

[Install]
WantedBy=timers.target
`

const socketFixture = `[Unit]
Description=Docker Socket for the API

[Socket]
ListenStream=/run/docker.sock
ListenStream=127.0.0.1:2375
SocketMode=0660
SocketUser=root
SocketGroup=docker

[Install]
WantedBy=sockets.target
`

// malformedNoClosingBracket: an unterminated section header (no "]" anywhere
// in the remaining text) — the go-systemd lexer hits EOF looking for it.
const malformedNoClosingBracket = "[Unit\nDescription=broken\n"

// malformedGarbageAfterHeader: trailing non-comment garbage on a
// section-header line, which go-systemd explicitly rejects (a systemd
// section header must have nothing else on its line).
const malformedGarbageAfterHeader = "[Unit] extra garbage\nDescription=x\n"

// malformedNoEqualsSign: a non-comment, non-blank line inside a section with
// no "=" before end-of-line — go-systemd requires every such line to be a
// key=value assignment.
const malformedNoEqualsSign = "[Unit]\nNoEqualsSignHere\n"
