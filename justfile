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