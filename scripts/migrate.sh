#!/bin/bash

if [ -f .env ]; then
    source .env
fi

DIRECTION="up"

if [ $# -eq 1 ]; then
	DIRECTION=$1
fi

cd sql/schema
goose turso $DATABASE_URL $DIRECTION
