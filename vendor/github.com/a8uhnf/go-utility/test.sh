#!/bin/bash

pushd queue
go test
popd

pushd stack
go test
popd

pushd filepath
go test
popd