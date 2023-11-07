# Details

Date : 2023-11-06 12:19:30

Directory /Users/sebastianflajszer/projects/go/ailingo-backend

Total : 55 files,  3537 codes, 146 comments, 682 blanks, all 4365 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [.idea/ailingo.iml](/.idea/ailingo.iml) | XML | 9 | 0 | 0 | 9 |
| [.idea/codeStyles/codeStyleConfig.xml](/.idea/codeStyles/codeStyleConfig.xml) | XML | 5 | 0 | 0 | 5 |
| [.idea/dataSources.xml](/.idea/dataSources.xml) | XML | 19 | 0 | 0 | 19 |
| [.idea/modules.xml](/.idea/modules.xml) | XML | 8 | 0 | 0 | 8 |
| [.idea/sqldialects.xml](/.idea/sqldialects.xml) | XML | 8 | 0 | 0 | 8 |
| [.idea/vcs.xml](/.idea/vcs.xml) | XML | 6 | 0 | 0 | 6 |
| [Dockerfile](/Dockerfile) | Docker | 11 | 0 | 2 | 13 |
| [README.md](/README.md) | Markdown | 14 | 0 | 3 | 17 |
| [cmd/api/main.go](/cmd/api/main.go) | Go | 17 | 0 | 7 | 24 |
| [config/config.go](/config/config.go) | Go | 70 | 2 | 13 | 85 |
| [docker-compose.dev.yaml](/docker-compose.dev.yaml) | YAML | 23 | 3 | 0 | 26 |
| [go.mod](/go.mod) | Go Module File | 30 | 0 | 4 | 34 |
| [go.sum](/go.sum) | Go Checksum File | 438 | 0 | 1 | 439 |
| [internal/app/app.go](/internal/app/app.go) | Go | 144 | 16 | 33 | 193 |
| [internal/controller/ai.go](/internal/controller/ai.go) | Go | 110 | 3 | 19 | 132 |
| [internal/controller/me.go](/internal/controller/me.go) | Go | 244 | 8 | 41 | 293 |
| [internal/controller/studyset.go](/internal/controller/studyset.go) | Go | 372 | 11 | 56 | 439 |
| [internal/domain/ai.go](/internal/domain/ai.go) | Go | 17 | 3 | 6 | 26 |
| [internal/domain/datastore.go](/internal/domain/datastore.go) | Go | 12 | 0 | 3 | 15 |
| [internal/domain/definition.go](/internal/domain/definition.go) | Go | 35 | 4 | 9 | 48 |
| [internal/domain/profile.go](/internal/domain/profile.go) | Go | 12 | 0 | 4 | 16 |
| [internal/domain/studysession.go](/internal/domain/studysession.go) | Go | 24 | 2 | 6 | 32 |
| [internal/domain/studyset.go](/internal/domain/studyset.go) | Go | 62 | 5 | 9 | 76 |
| [internal/domain/task.go](/internal/domain/task.go) | Go | 15 | 3 | 5 | 23 |
| [internal/domain/translation.go](/internal/domain/translation.go) | Go | 11 | 4 | 5 | 20 |
| [internal/domain/user.go](/internal/domain/user.go) | Go | 28 | 0 | 8 | 36 |
| [internal/gpt/repo.go](/internal/gpt/repo.go) | Go | 95 | 12 | 20 | 127 |
| [internal/gpt/worker.go](/internal/gpt/worker.go) | Go | 32 | 3 | 9 | 44 |
| [internal/mysql/datastore.go](/internal/mysql/datastore.go) | Go | 48 | 4 | 16 | 68 |
| [internal/mysql/db.go](/internal/mysql/db.go) | Go | 11 | 2 | 3 | 16 |
| [internal/mysql/definition.go](/internal/mysql/definition.go) | Go | 86 | 5 | 23 | 114 |
| [internal/mysql/errors.go](/internal/mysql/errors.go) | Go | 5 | 0 | 3 | 8 |
| [internal/mysql/profile.go](/internal/mysql/profile.go) | Go | 43 | 1 | 10 | 54 |
| [internal/mysql/studysession.go](/internal/mysql/studysession.go) | Go | 102 | 0 | 24 | 126 |
| [internal/mysql/studyset.go](/internal/mysql/studyset.go) | Go | 228 | 15 | 45 | 288 |
| [internal/mysql/task.go](/internal/mysql/task.go) | Go | 15 | 0 | 7 | 22 |
| [internal/mysql/user.go](/internal/mysql/user.go) | Go | 80 | 4 | 20 | 104 |
| [internal/usecase/definition.go](/internal/usecase/definition.go) | Go | 118 | 3 | 35 | 156 |
| [internal/usecase/errors.go](/internal/usecase/errors.go) | Go | 18 | 5 | 8 | 31 |
| [internal/usecase/gpt.go](/internal/usecase/gpt.go) | Go | 29 | 3 | 8 | 40 |
| [internal/usecase/profile.go](/internal/usecase/profile.go) | Go | 71 | 0 | 19 | 90 |
| [internal/usecase/studysession.go](/internal/usecase/studysession.go) | Go | 61 | 2 | 16 | 79 |
| [internal/usecase/studyset.go](/internal/usecase/studyset.go) | Go | 101 | 1 | 29 | 131 |
| [internal/usecase/translate.go](/internal/usecase/translate.go) | Go | 42 | 0 | 11 | 53 |
| [internal/usecase/user.go](/internal/usecase/user.go) | Go | 35 | 0 | 14 | 49 |
| [internal/webhook/clerk.go](/internal/webhook/clerk.go) | Go | 131 | 0 | 26 | 157 |
| [pkg/apiutil/apiutil.go](/pkg/apiutil/apiutil.go) | Go | 49 | 3 | 9 | 61 |
| [pkg/apiutil/error.go](/pkg/apiutil/error.go) | Go | 20 | 1 | 5 | 26 |
| [pkg/auth/auth.go](/pkg/auth/auth.go) | Go | 63 | 3 | 15 | 81 |
| [pkg/deepl/deepl.go](/pkg/deepl/deepl.go) | Go | 54 | 2 | 13 | 69 |
| [pkg/deepl/models.go](/pkg/deepl/models.go) | Go | 12 | 3 | 4 | 19 |
| [pkg/httpserver/server.go](/pkg/httpserver/server.go) | Go | 72 | 0 | 17 | 89 |
| [pkg/openai/models.go](/pkg/openai/models.go) | Go | 29 | 6 | 7 | 42 |
| [pkg/openai/openai.go](/pkg/openai/openai.go) | Go | 93 | 4 | 22 | 119 |
| [sql/init.sql](/sql/init.sql) | SQL | 50 | 0 | 10 | 60 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)