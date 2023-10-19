postgres_user := "user"

# migrate up the migrations against the cluster db
migrate_up_s1 POSTGRES_PWD:
	migrate -path backend/migrations/service1 -database postgres://{{postgres_user}}:{{POSTGRES_PWD}}@localhost:30001/{{postgres_user}}?sslmode=disable up

migrate_up_ghapp POSTGRES_PWD:
	migrate -path backend/migrations/github-app -database postgres://{{postgres_user}}:{{POSTGRES_PWD}}@localhost:30006/{{postgres_user}}?sslmode=disable up

migrate_down_ghapp POSTGRES_PWD:
	migrate -path backend/migrations/github-app -database postgres://{{postgres_user}}:{{POSTGRES_PWD}}@localhost:30006/{{postgres_user}}?sslmode=disable down


migrate_force_ghapp POSTGRES_PWD VER:
	migrate -path backend/migrations/github-app -database postgres://{{postgres_user}}:{{POSTGRES_PWD}}@localhost:30006/{{postgres_user}}?sslmode=disable force {{VER}}

## prune the entire cluster
k8s_prune PROJECT: 
	kubectl delete deploy,services,statefulset,pods -l project={{PROJECT}}

## create kafka client 
kafka_client_create:
	kubectl run kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.6.0-debian-11-r0 --namespace default --command -- sleep infinity 
	kubectl cp --namespace default config/client.properties kafka-client:/tmp/client.properties 
	kubectl exec --tty -i kafka-client --namespace default -- bash

kafka_client_setup:
	kubectl cp --namespace default config/client.properties kafka-client:/tmp/client.properties 

kafka_client_consume TOPIC:
	kubectl exec --tty -i kafka-client --namespace default -- kafka-console-consumer.sh --consumer.config /tmp/client.properties --bootstrap-server kafka.default.svc.cluster.local:9092 --topic {{TOPIC}} --from-beginning

protoc_gen_ts:
	rm -rf frontend/proto
	cp -r proto frontend/proto
	cd frontend && npx buf generate proto
	rm -rf frontend/proto

protoc_gen_go:
	protoc \
	--go_out backend/service/ \
	./proto/index.proto    

# Generate protoc for ts and go
protoc_gen: protoc_gen_ts protoc_gen_go