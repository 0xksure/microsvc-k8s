postgres_user := "user"

# migrate up the migrations against the cluster db
migrate_up POSTGRES_PWD:
	migrate -path backend/migrations -database postgres://{{postgres_user}}:{{POSTGRES_PWD}}@localhost:30001/{{postgres_user}}?sslmode=disable up