# Docker Compose Collection

One-command setup for **72 open source databases**. Every compose file in this directory is ready to run, no configuration needed.

```sh
cd docker/postgresql && docker compose up -d
```

That's it. Pick a database, `cd` into its folder, and spin it up.

## What's here

Each subdirectory contains a single `docker-compose.yaml` with sensible defaults: correct ports, initial credentials, persistent volumes, and automatic restart. Everything uses the `latest` tag so you always get the current stable release.

These aren't production configs. They're designed to get you a running instance in seconds for development, learning, or evaluation. Credentials are intentionally simple (usually the database name as the password).

## Quick reference

| Database | Image | Port | Run |
|----------|-------|------|-----|
| [PostgreSQL](postgresql/) | `postgres` | 5432 | `cd postgresql && docker compose up -d` |
| [MySQL](mysql/) | `mysql` | 3306 | `cd mysql && docker compose up -d` |
| [MariaDB](mariadb/) | `mariadb` | 3306 | `cd mariadb && docker compose up -d` |
| [MongoDB](mongodb/) | `mongo` | 27017 | `cd mongodb && docker compose up -d` |
| [Redis](redis/) | `redis` | 6379 | `cd redis && docker compose up -d` |
| [Valkey](valkey/) | `valkey/valkey` | 6379 | `cd valkey && docker compose up -d` |
| [KeyDB](keydb/) | `eqalpha/keydb` | 6379 | `cd keydb && docker compose up -d` |
| [Cassandra](cassandra/) | `cassandra` | 9042 | `cd cassandra && docker compose up -d` |
| [Neo4j](neo4j/) | `neo4j` | 7474 | `cd neo4j && docker compose up -d` |
| [ClickHouse](clickhouse/) | `clickhouse/clickhouse-server` | 8123 | `cd clickhouse && docker compose up -d` |
| [Elasticsearch](elasticsearch/) | `elasticsearch` | 9200 | `cd elasticsearch && docker compose up -d` |
| [InfluxDB](influxdb/) | `influxdb` | 8086 | `cd influxdb && docker compose up -d` |
| [CockroachDB](cockroachdb/) | `cockroachdb/cockroach` | 26257 | `cd cockroachdb && docker compose up -d` |
| [CouchDB](couchdb/) | `couchdb` | 5984 | `cd couchdb && docker compose up -d` |
| [ArangoDB](arangodb/) | `arangodb` | 8529 | `cd arangodb && docker compose up -d` |
| [RethinkDB](rethinkdb/) | `rethinkdb` | 28015 | `cd rethinkdb && docker compose up -d` |
| [SurrealDB](surrealdb/) | `surrealdb/surrealdb` | 8000 | `cd surrealdb && docker compose up -d` |
| [TiDB](tidb/) | `pingcap/tidb` | 4000 | `cd tidb && docker compose up -d` |
| [Milvus](milvus/) | `milvusdb/milvus` | 19530 | `cd milvus && docker compose up -d` |
| [Qdrant](qdrant/) | `qdrant/qdrant` | 6333 | `cd qdrant && docker compose up -d` |
| [Chroma](chroma/) | `chromadb/chroma` | 8000 | `cd chroma && docker compose up -d` |
| [Weaviate](weaviate/) | `semitechnologies/weaviate` | 8080 | `cd weaviate && docker compose up -d` |
| [MeiliSearch](meilisearch/) | `getmeili/meilisearch` | 7700 | `cd meilisearch && docker compose up -d` |
| [QuestDB](questdb/) | `questdb/questdb` | 9000 | `cd questdb && docker compose up -d` |
| [TDengine](tdengine/) | `tdengine/tdengine` | 6030 | `cd tdengine && docker compose up -d` |
| [Prometheus](prometheus/) | `prom/prometheus` | 9090 | `cd prometheus && docker compose up -d` |
| [VictoriaMetrics](victoriametrics/) | `victoriametrics/victoria-metrics` | 8428 | `cd victoriametrics && docker compose up -d` |

See below for the full list of all 72 databases.

## By category

### Relational

The classics and the newcomers. SQL is alive and well.

PostgreSQL, MySQL, MariaDB, ClickHouse, CockroachDB, TiDB, YugabyteDB, Vitess, Dolt, Firebird, OceanBase, StarRocks, Doris, Databend, RisingWave, GreptimeDB, QuestDB, TDengine, Manticore Search, PocketBase, TigerBeetle, Hydra, Neon, CrateDB, Trino

### Document stores

Flexible schemas for when you don't want to think about migrations.

MongoDB, CouchDB, Couchbase, FerretDB, ArangoDB, RethinkDB, Elasticsearch, Solr, MeiliSearch, Quickwit, RavenDB, SurrealDB, Chroma, ArcadeDB

### Key-value

The fastest way to store and retrieve data. Great for caching, sessions, and configuration.

Redis, Valkey, KeyDB, Memcached, etcd, FoundationDB, Aerospike, Hazelcast, Ignite, Geode, Tarantool, Immudb, Kvrocks, Pika, Garnet, Skytable, GridDB

### Graph

For data where relationships are the point.

Neo4j, DGraph, JanusGraph, ArangoDB, OrientDB, Kuzu, GUN, SurrealDB, ArcadeDB

### Time-series

Purpose-built for metrics, IoT, and monitoring data.

InfluxDB, QuestDB, TDengine, Prometheus, VictoriaMetrics, GreptimeDB, Druid

### Vector

The new kids. Built for embeddings, similarity search, and AI applications.

Milvus, Qdrant, Chroma, Weaviate

### Streaming and analytics

Real-time data processing and analytical queries.

Druid, Trino, StarRocks, Doris, Databend, RisingWave, Quickwit, Redpanda (via Kafka protocol)

## The numbers

```
Databases by data model
  Key/Value              ██████████████████████████████████████ 26
  Relational             █████████████████████████████ 21
  Document / XML         ██████████████████████ 16
  Graph                  ████████████ 9
  Column Family          ██████ 5
  Vector                 ████ 3
  Object-Relational      ████ 3
  Object-Oriented        ██ 2

Implementation language
  C++                    ██████████████████████████████████████ 22
  Java                   ██████████████████████████████ 18
  Go                     ██████████████████████████ 16
  C                      ███████████████████████ 14
  Rust                   ██████████████ 9
  Python                 ██████ 4
  JavaScript             ████ 3

Country of origin
  United States          ██████████████████████████████████████ 37
  China                  ██████████ 10
  United Kingdom         █████ 5
  Canada                 ███ 3
  Germany                ███ 3
  Sweden                 ███ 3

When they started
  2010s                  ██████████████████████████████████████ 37
  2000s                  █████████████████ 17
  2020s                  ████████████████ 16
  1990s                  █ 1
  1980s                  █ 1
```

Most of these databases were born in the 2010s, during the explosion of distributed systems and cloud computing. The 2020s cohort is heavily weighted toward vector databases and Rust-based systems, reflecting the current trends in AI and systems programming.

C++ and Java still dominate, but Go and Rust are gaining ground fast. Almost every database started after 2018 is written in one of those two.

## All 72 databases

| # | Database | Category | Image | Port |
|---|----------|----------|-------|------|
| 1 | [Aerospike](aerospike/) | Key/Value | `aerospike` | - |
| 2 | [ArangoDB](arangodb/) | Multi-model | `arangodb` | 8529 |
| 3 | [ArcadeDB](arcadedb/) | Multi-model | `arcadedata/arcadedb` | 2480 |
| 4 | [Cassandra](cassandra/) | Wide-column | `cassandra` | 9042 |
| 5 | [Chroma](chroma/) | Vector | `chromadb/chroma` | 8000 |
| 6 | [ClickHouse](clickhouse/) | Columnar | `clickhouse/clickhouse-server` | 8123 |
| 7 | [CockroachDB](cockroachdb/) | Distributed SQL | `cockroachdb/cockroach` | 26257 |
| 8 | [Couchbase](couchbase/) | Document | `couchbase` | 8091 |
| 9 | [CouchDB](couchdb/) | Document | `couchdb` | 5984 |
| 10 | [CrateDB](cratedb/) | Distributed SQL | `crate` | 4200 |
| 11 | [Databend](databend/) | Cloud warehouse | `datafuselabs/databend` | 8000 |
| 12 | [DGraph](dgraph/) | Graph | `dgraph/dgraph` | 8080 |
| 13 | [Dolt](dolt/) | Version-controlled SQL | `dolthub/dolt-sql-server` | 3306 |
| 14 | [Doris](doris/) | Analytics | `apache/doris` | 8030 |
| 15 | [Druid](druid/) | Analytics | `apache/druid` | 8888 |
| 16 | [Elasticsearch](elasticsearch/) | Search | `elasticsearch` | 9200 |
| 17 | [etcd](etcd/) | Key/Value | `quay.io/coreos/etcd` | 2379 |
| 18 | [FerretDB](ferretdb/) | Document | `ferretdb/ferretdb` | 27017 |
| 19 | [Firebird](firebird/) | Relational | `jacobalberty/firebird` | 3050 |
| 20 | [FoundationDB](foundationdb/) | Key/Value | `foundationdb/foundationdb` | 4500 |
| 21 | [Garnet](garnet/) | Cache | `ghcr.io/microsoft/garnet` | 6379 |
| 22 | [Geode](geode/) | In-memory | `apachegeode/geode` | 10334 |
| 23 | [GreptimeDB](greptimedb/) | Time-series | `greptime/greptimedb` | 4000 |
| 24 | [GridDB](griddb/) | Time-series | `griddb/griddb` | 10001 |
| 25 | [GUN](gun/) | Decentralized | `gundb/gun` | 8765 |
| 26 | [Hazelcast](hazelcast/) | In-memory | `hazelcast/hazelcast` | 5701 |
| 27 | [Hydra](hydra/) | Columnar Postgres | `ghcr.io/hydradatabase/hydra` | 5432 |
| 28 | [Ignite](ignite/) | In-memory | `apacheignite/ignite` | 10800 |
| 29 | [Immudb](immudb/) | Immutable | `codenotary/immudb` | 3322 |
| 30 | [InfluxDB](influxdb/) | Time-series | `influxdb` | 8086 |
| 31 | [JanusGraph](janusgraph/) | Graph | `janusgraph/janusgraph` | 8182 |
| 32 | [KeyDB](keydb/) | Key/Value | `eqalpha/keydb` | 6379 |
| 33 | [Kuzu](kuzu/) | Graph | `kuzudb/kuzu` | 8000 |
| 34 | [Kvrocks](kvrocks/) | Key/Value | `apache/kvrocks` | 6666 |
| 35 | [Manticore Search](manticore-search/) | Search | `manticoresearch/manticore` | 9306 |
| 36 | [MariaDB](mariadb/) | Relational | `mariadb` | 3306 |
| 37 | [MeiliSearch](meilisearch/) | Search | `getmeili/meilisearch` | 7700 |
| 38 | [Memcached](memcached/) | Cache | `memcached` | 11211 |
| 39 | [Milvus](milvus/) | Vector | `milvusdb/milvus` | 19530 |
| 40 | [MongoDB](mongodb/) | Document | `mongo` | 27017 |
| 41 | [MySQL](mysql/) | Relational | `mysql` | 3306 |
| 42 | [Neo4j](neo4j/) | Graph | `neo4j` | 7474 |
| 43 | [Neon](neon/) | Serverless Postgres | `ghcr.io/neondatabase/neon` | 5432 |
| 44 | [ObjectBox](objectbox/) | Embedded | `objectboxio/admin` | 8081 |
| 45 | [OceanBase](oceanbase/) | Distributed SQL | `oceanbase/oceanbase-ce` | 2881 |
| 46 | [OrientDB](orientdb/) | Multi-model | `orientdb` | 2480 |
| 47 | [Pika](pika/) | Key/Value | `pikadb/pika` | 9221 |
| 48 | [PocketBase](pocketbase/) | Backend-as-DB | `spectado/pocketbase` | 8090 |
| 49 | [PostgreSQL](postgresql/) | Relational | `postgres` | 5432 |
| 50 | [Prometheus](prometheus/) | Monitoring | `prom/prometheus` | 9090 |
| 51 | [QuestDB](questdb/) | Time-series | `questdb/questdb` | 9000 |
| 52 | [Quickwit](quickwit/) | Search | `quickwit/quickwit` | 7280 |
| 53 | [RavenDB](ravendb/) | Document | `ravendb/ravendb` | 8080 |
| 54 | [Redis](redis/) | Key/Value | `redis` | 6379 |
| 55 | [RethinkDB](rethinkdb/) | Document | `rethinkdb` | 28015 |
| 56 | [RisingWave](risingwave/) | Streaming SQL | `risingwavelabs/risingwave` | 4566 |
| 57 | [Skytable](skytable/) | Key/Value | `skytable/skytable` | 2003 |
| 58 | [Solr](solr/) | Search | `solr` | 8983 |
| 59 | [StarRocks](starrocks/) | Analytics | `starrocks/allin1-ubuntu` | 9030 |
| 60 | [SurrealDB](surrealdb/) | Multi-model | `surrealdb/surrealdb` | 8000 |
| 61 | [Tarantool](tarantool/) | In-memory | `tarantool/tarantool` | 3301 |
| 62 | [TDengine](tdengine/) | Time-series | `tdengine/tdengine` | 6030 |
| 63 | [TiDB](tidb/) | Distributed SQL | `pingcap/tidb` | 4000 |
| 64 | [TigerBeetle](tigerbeetle/) | Financial | `tigerbeetle/tigerbeetle` | 3001 |
| 65 | [TiKV](tikv/) | Key/Value | `pingcap/tikv` | 20160 |
| 66 | [Trino](trino/) | Query engine | `trinodb/trino` | 8080 |
| 67 | [Valkey](valkey/) | Key/Value | `valkey/valkey` | 6379 |
| 68 | [VictoriaMetrics](victoriametrics/) | Time-series | `victoriametrics/victoria-metrics` | 8428 |
| 69 | [Vitess](vitess/) | MySQL scaling | `vitess/lite` | 15306 |
| 70 | [Weaviate](weaviate/) | Vector | `semitechnologies/weaviate` | 8080 |
| 71 | [YugabyteDB](yugabytedb/) | Distributed SQL | `yugabytedb/yugabyte` | 5433 |
| 72 | [Zookeeper](zookeeper/) | Coordination | `zookeeper` | 2181 |

## Tips

**Start something up fast:**
```sh
cd docker/redis && docker compose up -d
```

**Stop and remove everything:**
```sh
cd docker/redis && docker compose down -v
```

**Check logs:**
```sh
cd docker/redis && docker compose logs -f
```

**Try a few databases side by side** (they all use different ports, so no conflicts):
```sh
cd docker/postgresql && docker compose up -d
cd ../mongodb && docker compose up -d
cd ../redis && docker compose up -d
```
