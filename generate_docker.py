#!/usr/bin/env python3
"""Generate docker-compose.yaml files for open source databases with Docker Hub images.

Reads docker_matches.txt (slug|image|port) and databases.json,
creates docker/{slug}/docker-compose.yaml for each match.
"""

import json
import os
import subprocess
import sys

# Database-specific configurations: (extra_ports, env_vars, volumes, command, notes)
# env_vars is a dict, volumes is a list of "name:/path" pairs
DB_CONFIG = {
    "postgresql": {
        "env": {"POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "postgres"},
        "volumes": ["pgdata:/var/lib/postgresql/data"],
    },
    "mysql": {
        "env": {"MYSQL_ROOT_PASSWORD": "mysql", "MYSQL_DATABASE": "mydb"},
        "volumes": ["mysqldata:/var/lib/mysql"],
    },
    "mariadb": {
        "env": {"MARIADB_ROOT_PASSWORD": "mariadb", "MARIADB_DATABASE": "mydb"},
        "volumes": ["mariadbdata:/var/lib/mysql"],
    },
    "mongodb": {
        "env": {"MONGO_INITDB_ROOT_USERNAME": "root", "MONGO_INITDB_ROOT_PASSWORD": "mongo"},
        "volumes": ["mongodata:/data/db"],
    },
    "redis": {
        "volumes": ["redisdata:/data"],
    },
    "valkey": {
        "volumes": ["valkeydata:/data"],
    },
    "keydb": {
        "volumes": ["keydbdata:/data"],
    },
    "memcached": {},
    "cassandra": {
        "env": {"CASSANDRA_CLUSTER_NAME": "TestCluster"},
        "volumes": ["cassandradata:/var/lib/cassandra"],
    },
    "neo4j": {
        "ports": ["7687:7687"],
        "env": {"NEO4J_AUTH": "neo4j/password"},
        "volumes": ["neo4jdata:/data"],
    },
    "influxdb": {
        "env": {"DOCKER_INFLUXDB_INIT_MODE": "setup", "DOCKER_INFLUXDB_INIT_USERNAME": "admin",
                "DOCKER_INFLUXDB_INIT_PASSWORD": "password", "DOCKER_INFLUXDB_INIT_ORG": "myorg",
                "DOCKER_INFLUXDB_INIT_BUCKET": "mybucket"},
        "volumes": ["influxdata:/var/lib/influxdb2"],
    },
    "elasticsearch": {
        "ports": ["9300:9300"],
        "env": {"discovery.type": "single-node", "xpack.security.enabled": "false"},
        "volumes": ["esdata:/usr/share/elasticsearch/data"],
    },
    "clickhouse": {
        "ports": ["9000:9000"],
        "volumes": ["clickhousedata:/var/lib/clickhouse"],
    },
    "cockroachdb": {
        "ports": ["8080:8080"],
        "command": "start-single-node --insecure",
        "volumes": ["crdbdata:/cockroach/cockroach-data"],
    },
    "tidb": {},
    "tikv": {},
    "vitess": {},
    "rethinkdb": {
        "ports": ["8080:8080"],
        "volumes": ["rethinkdata:/data"],
    },
    "arangodb": {
        "env": {"ARANGO_ROOT_PASSWORD": "arangodb"},
        "volumes": ["arangodata:/var/lib/arangodb3"],
    },
    "orientdb": {
        "ports": ["2424:2424"],
        "env": {"ORIENTDB_ROOT_PASSWORD": "orientdb"},
        "volumes": ["orientdata:/orientdb/databases"],
    },
    "dgraph": {
        "command": "dgraph zero",
    },
    "scylladb": {
        "volumes": ["scylladata:/var/lib/scylla"],
    },
    "couchdb": {
        "env": {"COUCHDB_USER": "admin", "COUCHDB_PASSWORD": "couchdb"},
        "volumes": ["couchdata:/opt/couchdb/data"],
    },
    "couchbase": {
        "ports": ["8091:8091", "8092:8092", "8093:8093", "11210:11210"],
    },
    "questdb": {
        "ports": ["9009:9009", "8812:8812"],
        "volumes": ["questdata:/var/lib/questdb"],
    },
    "surrealdb": {
        "command": "start --user root --pass surrealdb file:/data",
        "volumes": ["surrealdata:/data"],
    },
    "etcd": {
        "env": {"ALLOW_NONE_AUTHENTICATION": "yes"},
    },
    "foundationdb": {},
    "hazelcast": {},
    "ignite": {},
    "druid": {},
    "solr": {
        "volumes": ["solrdata:/var/solr"],
    },
    "meilisearch": {
        "env": {"MEILI_MASTER_KEY": "masterKey"},
        "volumes": ["meilidata:/meili_data"],
    },
    "typesense": {
        "command": "--data-dir /data --api-key=typesense",
        "volumes": ["typesensedata:/data"],
    },
    "nats": {},
    "eventstore": {
        "env": {"EVENTSTORE_INSECURE": "true", "EVENTSTORE_RUN_PROJECTIONS": "All",
                "EVENTSTORE_CLUSTER_SIZE": "1"},
    },
    "timescaledb": {
        "env": {"POSTGRES_PASSWORD": "postgres"},
        "volumes": ["tsdata:/var/lib/postgresql/data"],
    },
    "yugabytedb": {
        "ports": ["7000:7000", "9000:9000", "15433:15433"],
        "command": "bin/yugabyted start --daemon=false",
    },
    "percona-server-for-mysql": {
        "env": {"MYSQL_ROOT_PASSWORD": "percona"},
        "volumes": ["perconadata:/var/lib/mysql"],
    },
    "firebird": {
        "env": {"FIREBIRD_DATABASE": "mydb.fdb", "FIREBIRD_USER": "firebird",
                "FIREBIRD_PASSWORD": "firebird", "ISC_PASSWORD": "masterkey"},
        "volumes": ["firebirddata:/firebird/data"],
    },
    "tarantool": {
        "volumes": ["tarantooldata:/var/lib/tarantool"],
    },
    "geode": {},
    "cratedb": {
        "env": {"CRATE_HEAP_SIZE": "512m"},
        "command": "-Cdiscovery.type=single-node",
        "volumes": ["cratedata:/data"],
    },
    "ravendb": {
        "env": {"RAVEN_Security_UnsecuredAccessAllowed": "PublicNetwork",
                "RAVEN_Setup_Mode": "None"},
        "volumes": ["ravendata:/opt/RavenDB/Server/RavenData"],
    },
    "janusgraph": {},
    "presto": {},
    "trino": {},
    "stardog": {},
    "weaviate": {
        "env": {"AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED": "true",
                "PERSISTENCE_DATA_PATH": "/var/lib/weaviate"},
        "volumes": ["weaviatedata:/var/lib/weaviate"],
    },
    "milvus": {
        "env": {"ETCD_USE_EMBED": "true", "COMMON_STORAGETYPE": "local"},
        "volumes": ["milvusdata:/var/lib/milvus"],
    },
    "qdrant": {
        "ports": ["6334:6334"],
        "volumes": ["qdrantdata:/qdrant/storage"],
    },
    "chroma": {
        "volumes": ["chromadata:/chroma/chroma"],
    },
    "prometheus": {
        "volumes": ["promdata:/prometheus"],
    },
    "victoriametrics": {
        "volumes": ["vmdata:/victoria-metrics-data"],
    },
    "immudb": {
        "volumes": ["immudata:/var/lib/immudb"],
    },
    "manticore-search": {
        "volumes": ["manticoredata:/var/lib/manticore"],
    },
    "tdengine": {
        "ports": ["6041:6041"],
        "volumes": ["tdenginedata:/var/lib/taos"],
    },
    "ferretdb": {
        "env": {"FERRETDB_POSTGRESQL_URL": "postgres://postgres:postgres@postgres:5432/ferretdb"},
    },
    "dolt": {
        "env": {"DOLT_ROOT_PASSWORD": "dolt"},
        "volumes": ["doltdata:/var/lib/dolt"],
    },
    "arcadedb": {
        "env": {"JAVA_OPTS": "-Darcadedb.server.rootPassword=arcadedb"},
        "volumes": ["arcadedata:/home/arcadedb/databases"],
    },
    "edgedb": {
        "env": {"EDGEDB_SERVER_SECURITY": "insecure_dev_mode"},
        "volumes": ["edgedata:/var/lib/edgedb/data"],
        "ports": ["5656:5656"],
    },
    "oceanbase": {
        "env": {"MODE": "mini"},
    },
    "risingwave": {
        "command": "playground",
    },
    "starrocks": {},
    "doris": {},
    "databend": {},
    "greptimedb": {},
    "pocketbase": {
        "volumes": ["pbdata:/pb_data"],
    },
    "skytable": {},
    "aerospike": {
        "volumes": ["aerospikedata:/opt/aerospike/data"],
    },
    "quickwit": {},
    "tigerbeetle": {},
    "zookeeper": {
        "ports": ["2181:2181"],
    },
    "garnet": {},
    "dragonflydb": {
        "volumes": ["dragonflydata:/data"],
    },
    "kvrocks": {},
    "pika": {},
    "griddb": {},
    "gun": {},
    "objectbox": {},
    "zombodb": {},
    "hydra": {
        "env": {"POSTGRES_PASSWORD": "postgres"},
        "volumes": ["hydradata:/var/lib/postgresql/data"],
    },
    "neon": {},
    "opengauss": {
        "env": {"GS_PASSWORD": "OpenGauss@123"},
        "volumes": ["gaussdata:/var/lib/opengauss/data"],
    },
    "polardb": {},
    "kuzu": {},
    "materialize": {},
    "redpanda": {
        "command": "redpanda start --smp 1 --memory 512M --overprovisioned",
    },
    "singlestore": {
        "env": {"ROOT_PASSWORD": "singlestore", "SINGLESTORE_LICENSE": ""},
    },
    "loki": {},
    "grafana-tempo": {},
    "openldap": {
        "env": {"LDAP_ORGANISATION": "Example", "LDAP_DOMAIN": "example.org",
                "LDAP_ADMIN_PASSWORD": "admin"},
    },
    "fauna": {},
    "apache-pinot": {
        "command": "QuickStart -type batch",
    },
    "apache-iotdb": {},
}

# Load matches
matches = []
with open("docker_matches.txt") as f:
    for line in f:
        line = line.strip()
        if not line:
            continue
        slug, image, port = line.split("|")
        matches.append((slug, image, int(port)))

# Load database names from cache
db_names = {}
with open("databases.json") as f:
    for d in json.load(f):
        db_names[d["slug"]] = d["name"]

print(f"Generating docker-compose.yaml for {len(matches)} databases...")

results = []

for slug, image, port in matches:
    name = db_names.get(slug, slug)
    config = DB_CONFIG.get(slug, {})
    env = config.get("env", {})
    volumes = config.get("volumes", [])
    extra_ports = config.get("ports", [])
    command = config.get("command", "")
    service_name = slug.replace("-", "_")

    # Build docker-compose.yaml
    lines = ["services:", f"  {service_name}:", f"    image: {image}:latest"]
    if command:
        lines.append(f"    command: {command}")

    # Ports
    all_ports = []
    if port > 0:
        all_ports.append(f"{port}:{port}")
    all_ports.extend(extra_ports)
    if all_ports:
        lines.append("    ports:")
        for p in all_ports:
            lines.append(f'      - "{p}"')

    # Environment
    if env:
        lines.append("    environment:")
        for k, v in env.items():
            lines.append(f"      {k}: \"{v}\"")

    # Volumes (service level)
    if volumes:
        lines.append("    volumes:")
        for v in volumes:
            lines.append(f"      - {v}")

    # Restart policy
    lines.append("    restart: unless-stopped")

    # Top-level volumes
    if volumes:
        lines.append("")
        lines.append("volumes:")
        for v in volumes:
            vol_name = v.split(":")[0]
            lines.append(f"  {vol_name}:")

    content = "\n".join(lines) + "\n"

    # Write file
    dirpath = f"docker/{slug}"
    os.makedirs(dirpath, exist_ok=True)
    filepath = f"{dirpath}/docker-compose.yaml"
    with open(filepath, "w") as f:
        f.write(content)

    results.append((slug, name, image, filepath))
    print(f"  created {filepath}")

# Write results for commit script
with open("docker_results.txt", "w") as f:
    for slug, name, image, filepath in results:
        f.write(f"{slug}|{name}|{image}|{filepath}\n")

print(f"\nDone: {len(results)} docker-compose files created")
