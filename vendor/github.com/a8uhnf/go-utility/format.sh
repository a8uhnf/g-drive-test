#!/usr/bin/env bash

goimports -w queue stack
gofmt -s -w queue stack
