#!/usr/bin/bash

# Use este script para executar testes locais

# Need to have the following structure:
# tools
# ├── gatling       --> https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.10.3/gatling-charts-highcharts-bundle-3.10.3-bundle.zip
# └── load-test     --> https://github.com/zanfranceschi/rinha-de-backend-2024-q1/tree/main/load-test

TOOLS_DIR="$(pwd)/bin/tools"

export GATLING_HOME=$TOOLS_DIR/gatling
RESULTS_WORKSPACE="$TOOLS_DIR/load-test/user-files/results"
GATLING_BIN_DIR=$TOOLS_DIR/gatling/bin
GATLING_WORKSPACE="$TOOLS_DIR/load-test/user-files"

runGatling() {
    sh $GATLING_BIN_DIR/gatling.sh -rm local -s RinhaBackendCrebitosSimulation \
        -rd "Rinha de Backend - 2024/Q1: Crébito" \
        -rf $RESULTS_WORKSPACE \
        -sf "$GATLING_WORKSPACE/simulations"
}

startTest() {
    for i in {1..20}; do
        # 2 requests to wake the 2 api instances up :)
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        runGatling && \
        break || sleep 2;
    done
}

startTest
