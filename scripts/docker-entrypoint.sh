#! /bin/sh

make gitHooks
make build
swag init
make watch

make run