#!/bin/bash

# Need to have the following structure:
# tools
# ├── gatling       --> https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.10.3/gatling-charts-highcharts-bundle-3.10.3-bundle.zip
# └── load-test     --> https://github.com/zanfranceschi/rinha-de-backend-2024-q1/tree/main/load-test

mkdir -p bin/tools
cd bin/tools

# load-test
curl -LO https://codeload.github.com/zanfranceschi/rinha-de-backend-2024-q1/tar.gz/main
tar -xz --strip=1 rinha-de-backend-2024-q1-main/load-test < main
rm main

# gatling
curl -LO https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.10.3/gatling-charts-highcharts-bundle-3.10.3-bundle.zip
unzip gatling-charts-highcharts-bundle-3.10.3-bundle.zip
mv gatling-charts-highcharts-bundle-3.10.3 gatling
rm gatling-charts-highcharts-bundle-3.10.3-bundle.zip
