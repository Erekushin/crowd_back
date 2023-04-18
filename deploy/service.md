[Unit]
Description=Template Backend Go
After=syslog.target
StartLimitIntervalSec=0

[Service]
Type=simple
User=service
WorkingDirectory=/home/service/template-backend/
ExecStart=/home/service/template-backend/go-backend.gerege.mn

StandardOutput=syslog
SyslogIdentifier=template-backend

SuccessExitStatus=143
TimeoutStopSec=10
Restart=on-failure
RestartSec=10

StandardError=syslog

[Install]
WantedBy=multi-user.target