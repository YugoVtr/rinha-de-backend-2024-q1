#!/bin/bash

go test -v -count=1 --tags=integration -run ^TestIntegration$ github.com/yugovtr/rinha-de-backend-2024-q1/server
