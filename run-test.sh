#!/bin/sh

docker-compose -f docker-compose.test.yml build
docker-compose -f docker-compose.test.yml run --rm test
rc=$?
docker-compose -f docker-compose.test.yml logs
docker-compose -f docker-compose.test.yml down

echo "Leaving test with exit code $rc"
exit $rc
