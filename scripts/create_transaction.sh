#!/bin/sh
curl -X POST http://localhost:8080 \
   -H 'Content-Type: application/json' \
   -d '{"BPM":60}'
