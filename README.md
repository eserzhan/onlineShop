# Shop
### To launch the application:


1. Clone the repository
 


2. Create .env file in root directory and add following values:
```dotenv
PASSWORD=<Password>
JWT.SigningKey=<RandomSymbols>
Password.salt=<RandomSymbols>
```

3. Up the Postgres db, you can do that by Docker:
   
```dotenv
docker run --name postgres-container -e POSTGRES_PASSWORD=<Password> -p 5432:5432 -d postgres
```
4. Install migrate CLI and apply migrations:

```dotenv
migrate -database "postgres://postgres:<Password>@localhost:5432/postgres?sslmode=disable" -path "db/migrations" down  
```
