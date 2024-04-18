#!/bin/bash

TOPIC="im-push-topic"

# 删除 Kafka topic
docker-compose exec kafka kafka-topics.sh --delete --topic $TOPIC --bootstrap-server kafka:9092

# 等待一段时间以确保 topic 已经被删除
sleep 5

# 创建 Kafka topic
docker-compose exec kafka kafka-topics.sh --create --topic $TOPIC --partitions 1 --replication-factor 1 --bootstrap-server kafka:9092

# 启动 Kafka 消费者来消费 topic
docker-compose exec kafka kafka-console-consumer.sh --topic $TOPIC --from-beginning --bootstrap-server kafka:9092 &

# 启动 Kafka 生产者来向 topic 发送消息
docker-compose exec kafka kafka-console-producer.sh --topic $TOPIC --broker-list kafka:9092


#protoc --proto_path=. --go_out=. --go-grpc_out=. --validate_out="lang=go:." ./api/logic/logic.proto

