#!/bin/bash

gofmt -w .
goimports -local github.com/bohenriksen2020/ms-orders-api -w .
