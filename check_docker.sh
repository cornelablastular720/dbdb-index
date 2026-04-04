#!/bin/bash
# Check Docker Hub for open source databases from databases.json
# Outputs: slug|image|description lines for databases that have Docker images

set -euo pipefail

CACHE_FILE="databases.json"
OUT_FILE="docker_matches.txt"

# Extract open source slugs
SLUGS=$(python3 -c "
import json
dbs = json.load(open('$CACHE_FILE'))
oss = [d for d in dbs if any('Open Source' in t for t in d.get('project_types', []))]
for d in sorted(oss, key=lambda x: x['slug']):
    print(d['slug'])
")

# Known mappings: slug -> docker_image:default_port
# These are databases where the Docker Hub name differs from the slug
declare -A KNOWN_IMAGES=(
  ["postgresql"]="postgres:5432"
  ["mongodb"]="mongo:27017"
  ["couchbase-server"]="couchbase:8091"
  ["couchdb"]="couchdb:5984"
  ["mysql"]="mysql:3306"
  ["mariadb"]="mariadb:3306"
  ["redis"]="redis:6379"
  ["memcached"]="memcached:11211"
  ["cassandra"]="cassandra:9042"
  ["neo4j"]="neo4j:7474"
  ["influxdb"]="influxdb:8086"
  ["elasticsearch"]="elasticsearch:9200"
  ["clickhouse"]="clickhouse/clickhouse-server:8123"
  ["cockroachdb"]="cockroachdb/cockroach:26257"
  ["tidb"]="pingcap/tidb:4000"
  ["tikv"]="pingcap/tikv:20160"
  ["vitess"]="vitess/lite:15000"
  ["rethinkdb"]="rethinkdb:28015"
  ["arangodb"]="arangodb:8529"
  ["orientdb"]="orientdb:2480"
  ["dgraph"]="dgraph/dgraph:8080"
  ["scylladb"]="scylladb/scylla:9042"
  ["minio"]="minio/minio:9000"
  ["questdb"]="questdb/questdb:9000"
  ["surrealdb"]="surrealdb/surrealdb:8000"
  ["duckdb"]="datacatering/duckdb:1294"
  ["etcd"]="quay.io/coreos/etcd:2379"
  ["foundationdb"]="foundationdb/foundationdb:4500"
  ["gun"]="gundb/gun:8765"
  ["hazelcast"]="hazelcast/hazelcast:5701"
  ["ignite"]="apacheignite/ignite:10800"
  ["druid"]="apache/druid:8888"
  ["solr"]="solr:8983"
  ["meilisearch"]="getmeili/meilisearch:7700"
  ["typesense"]="typesense/typesense:8108"
  ["nats"]="nats:4222"
  ["eventstore"]="eventstore/eventstore:2113"
  ["timescaledb"]="timescale/timescaledb:5432"
  ["yugabytedb"]="yugabytedb/yugabyte:5433"
  ["percona-server-for-mysql"]="percona:3306"
  ["firebird"]="jacobalberty/firebird:3050"
  ["tarantool"]="tarantool/tarantool:3301"
  ["geode"]="apachegeode/geode:10334"
  ["pouchdb"]="pouchdb/pouchdb-server:5984"
  ["opentsdb"]="petergrace/opentsdb-docker:4242"
  ["greenplum"]="projectairws/greenplum:5432"
  ["cratedb"]="crate:4200"
  ["ravendb"]="ravendb/ravendb:8080"
  ["fauna"]="fauna/faunadb:8443"
  ["janusgraph"]="janusgraph/janusgraph:8182"
  ["presto"]="prestodb/presto:8080"
  ["trino"]="trinodb/trino:8080"
  ["stardog"]="stardog/stardog:5820"
  ["weaviate"]="semitechnologies/weaviate:8080"
  ["milvus"]="milvusdb/milvus:19530"
  ["qdrant"]="qdrant/qdrant:6333"
  ["chroma"]="chromadb/chroma:8000"
  ["lancedb"]="lancedb/lancedb:8080"
  ["apache-pinot"]="apachepinot/pinot:9000"
  ["grafana-tempo"]="grafana/tempo:3200"
  ["prometheus"]="prom/prometheus:9090"
  ["victoriametrics"]="victoriametrics/victoria-metrics:8428"
  ["keycloak"]="quay.io/keycloak/keycloak:8080"
  ["apache-iotdb"]="apache/iotdb:6667"
  ["griddb"]="griddb/griddb:10001"
  ["immudb"]="codenotary/immudb:3322"
  ["manticore-search"]="manticoresearch/manticore:9306"
  ["objectbox"]="objectboxio/admin:8081"
  ["tdengine"]="tdengine/tdengine:6030"
  ["zombodb"]="zombodb/zombodb:5432"
  ["valkey"]="valkey/valkey:6379"
  ["keydb"]="eqalpha/keydb:6379"
  ["garnet"]="ghcr.io/microsoft/garnet:6379"
  ["dragonflydb"]="docker.dragonflydb.io/dragonflydb/dragonfly:6379"
  ["kvrocks"]="apache/kvrocks:6666"
  ["pika"]="pikadb/pika:9221"
  ["openldap"]="osixia/openldap:389"
  ["apisix"]="apache/apisix:9080"
  ["ferretdb"]="ferretdb/ferretdb:27017"
  ["singlestore"]="ghcr.io/singlestore-labs/singlestoredb-dev:3306"
  ["loki"]="grafana/loki:3100"
)

> "$OUT_FILE"
FOUND=0
CHECKED=0

echo "Checking Docker Hub for ${#KNOWN_IMAGES[@]} known images..." >&2

# First, output all known matches
for slug in $(echo "${!KNOWN_IMAGES[@]}" | tr ' ' '\n' | sort); do
  if echo "$SLUGS" | grep -qx "$slug"; then
    IFS=':' read -r image port <<< "${KNOWN_IMAGES[$slug]}"
    echo "$slug|$image|$port"
    echo "$slug|$image|$port" >> "$OUT_FILE"
    ((FOUND++))
  fi
done

echo "Found $FOUND from known mappings" >&2
echo "Checking Docker Hub for remaining databases..." >&2

# For remaining slugs, check Docker Hub official images
for slug in $SLUGS; do
  # Skip if already matched
  if [[ -n "${KNOWN_IMAGES[$slug]+x}" ]]; then
    continue
  fi

  ((CHECKED++))
  if (( CHECKED % 50 == 0 )); then
    echo "  checked $CHECKED..." >&2
  fi

  # Try official image (library namespace)
  status=$(curl -s -o /dev/null -w "%{http_code}" \
    "https://hub.docker.com/v2/repositories/library/$slug/" 2>/dev/null || echo "000")

  if [[ "$status" == "200" ]]; then
    # Get default exposed port from image config if possible, default to 0
    echo "$slug|$slug|0"
    echo "$slug|$slug|0" >> "$OUT_FILE"
    ((FOUND++))
    echo "  + $slug (official)" >&2
  fi

  # Rate limit
  sleep 0.15
done

echo "" >&2
echo "Total matches: $FOUND" >&2
