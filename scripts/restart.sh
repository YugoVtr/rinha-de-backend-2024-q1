#!/bin/bash

go test -run ^TestResetDB$ -count=1 -tags=integration github.com/yugovtr/rinha-de-backend-2024-q1/test
docker compose restart api01 api02
