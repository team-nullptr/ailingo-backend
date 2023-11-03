Repo provides you with `docker-compose.dev.yml` so that you can run the app easily in a containerized environment.

To get started create a `.env` file according to `.env.template` file. You can just copy and rename it. Then run
```shell
docker compose -f docker-compose.dev.yaml up db
```
To start the database. When started run
```shell
docker compose -f docker-compose.dev.yaml up backend
```

Then remove all clerk users, expose backend via ngrok and update webhook address.

If you want to start just add `backend` or `db` to the end of the previous command.
```shell
docker compose -f docker-compose.dev.yaml up db
```