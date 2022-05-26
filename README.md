# todo-app
Todo-app written mainly in Golang. PostgreSQL is used for all storage except sessions, which we use Redis for. 

1. Make sure you have `git` and `docker` installed.

2. Run:

```
git clone https://github.com/prabhdatnoor/todo-app/
docker-compose up
```

3. Debug:
```
docker-compose --file docker-compose-debug.yml up
```
DB runs at port 5438 external, backend is on 5001 external. Feel free to change it to whatever you like in the `docker-compose.yml` and `backend/Dockerfile `

Currently `db/password.txt` is not needed and the default "postgres" password is used for the DB
