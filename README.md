
# Banking System in Golang 

Banking System, a classic web server designed to provide essential banking functionalities. This system allows users to sign up and log in, providing a seamless experience for managing their accounts.

PostgreSQL Database is used to store the data in relational manner.

**Features**

User Authentication: Secure user signup and login mechanisms ensure the confidentiality of user accounts.

**Account Management:**

○ Create new accounts with ease.
○ View account summary to keep track of transactions and balances.

○ Conveniently add funds to your account.

○ Effortlessly transfer funds to another account.



## Environment Variables

To run this project, you will need to add the following environment variables 

`JWT_SECRET` using command for windows
```
set JWT_SECRET="...."
```


And update your Postgres connection string in storage.go file: 

```
"postgres://postgres:<password>@localhost/<db_name>?sslmode=disable"
```



## Installation

Install go packages with following commands

```
    github.com/golang-jwt/jwt/v5
```
```
    golang.org/x/crypto
```

```
    github.com/gorilla/mux
```
```
    github.com/lib/pq
```
    
## Running the Server

Use the below command to run the server to play around!
```
go run .
```

