#!/bin/bash

echo "Updating code from Git..."
sudo git fetch origin && sudo git reset --hard origin/main

echo "Building the application ... "
go build -buildvcs=false

echo "Reloading systemd..."
sudo systemctl daemon-reload

echo "Restarting the peec service..."
sudo systemctl restart api

echo "Deployment completed successfully."
