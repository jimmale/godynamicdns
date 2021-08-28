#!/bin/bash

## Get it to notice the systemd unit file
systemctl daemon-reload

echo "Edit the /etc/godynamicdns/config.toml file and then run"
echo "  sudo systemctl enable godynamicdns.service"
echo "  sudo systemctl start godynamicdns.service"