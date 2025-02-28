#!/bin/sh
go install github.com/air-verse/air@latest
go install github.com/a-h/templ/cmd/templ@latest
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs

