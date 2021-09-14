#! /bin/bash

PROJECT="boilerplate"

# run docker compose
docker-compose up -d consul db 

# wait for consul container be ready
while ! curl --request GET -sL --url 'http://localhost:8590/' > /dev/null 2>&1; do printf .; sleep 1; done
# setting KV, dependency of app
curl --request PUT --data-binary @config.local.json http://localhost:8590/v1/kv/${PROJECT}

docker-compose up -d --build ${PROJECT}
