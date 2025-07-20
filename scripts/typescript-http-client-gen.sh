#!/bin/sh

docker run -v ../web/src/api:/app/output \
-v ../backend/cmd/api/docs/swagger.yaml:/app/openapi.yaml \
--user $(id -u):$(id -g) \
mcr.microsoft.com/openapi/kiota generate --language typescript