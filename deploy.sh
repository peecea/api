#!/bin/bash

echo "Updating code from Git..."
cd /api && sudo git fetch origin && sudo git reset --hard origin/master

echo "Building the application ... "
cd /api && go build -buildvcs=false

echo "Reloading systemd..."
sudo systemctl daemon-reload

echo "Restarting the peec service..."
sudo systemctl restart api

echo "Deployment completed successfully."
