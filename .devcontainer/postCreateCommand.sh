#!/bin/sh
go install github.com/air-verse/air@latest
go install github.com/a-h/templ/cmd/templ@latest
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo bash -
sudo apt-get install -y nodejs

