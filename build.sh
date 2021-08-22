#!/user/bin/env zsh
set -xe

# install packages and dependencies
go get github.com/gin-gonic/gin

go get github.com/go-playground/validator/v10

#build command
go build -o bin/application server.go