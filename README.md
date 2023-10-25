Repo provides you with `docker-compose.dev.yml` so that you can run the app easily in a contenerized environment.

To get started create a `.env` file according to `.env.template` file. You can just copy and rename it. Then run
```shell
docker compose -f docker-compose.dev.yaml up
```

If you want to start just add `backend` or `db` to the end of the previous command.
```shell
docker compose -f docker-compose.dev.yaml up db
```