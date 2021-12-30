# kfchc / Kafka Connect HealthCheck

Kafka Connect (connectors / tasks) HealthCheck For AWS ALB and more

## commands

###  connectors:gen_server_config

generating Kafka Connect connectors config.

```bash
$ kfchc gsc --path ./
```

```json
{
  "servers": [
    {
      "connectServer": "http://127.0.0.1:8083",
      "connectors": [
        "replace_me",
        "replace_me"
      ]
    }
  ]
}
```

### connectors:health_check

Kafka Connect connectors and clients using the REST interface.

```bash
$ kfchc c:conns --config_file ./path/to/servers.json
```

## example

```bash
$ docker compose exec kafka /opt/bitnami/kafka/bin/kafka-topics.sh --create --if-not-exists --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1 --topic test_topic_json
$ docker compose exec kafka /opt/bitnami/kafka/bin/kafka-console-producer.sh --topic test_topic_json --bootstrap-server kafka:9092 
```

```bash
$ curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d '{
"name": "file_sink",
"config": {
"connector.class": "org.apache.kafka.connect.file.FileStreamSinkConnector",
"tasks.max": "1",
"file": "/tmp/output.txt",
"topics": "test_topic_json",
"key.converter": "org.apache.kafka.connect.storage.StringConverter",
"value.converter": "org.apache.kafka.connect.storage.StringConverter"
}
}'

$ curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d '{
"name": "file_json_sink",
"config": {
"connector.class": "org.apache.kafka.connect.file.FileStreamSinkConnector",
"tasks.max": "1",
"file": "/tmp/json_output.txt",
"topics": "test_topic_json",
"value.converter": "org.apache.kafka.connect.json.JsonConverter",
"value.converter.schemas.enable": "false",
"key.converter": "org.apache.kafka.connect.json.JsonConverter",
"key.converter.schemas.enable": "false"
}
}'
```

## protobuf

```bash
$ protoc --proto_path=protobuf --go_out=pbdef --go_opt=paths=source_relative protobuf/config.proto
```
