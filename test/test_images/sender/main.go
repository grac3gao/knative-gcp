/*
Copyright 2019 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	//"time"
	//
	//"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	DBName string `envconfig:"DB_NAME" required:"true"`
	DBUser string `envconfig:"DB_USER" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		panic(fmt.Sprintf("Failed to process env var: %s", err))
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:@tcp(127.0.0.1:3306)/%s", env.DBUser, env.DBName))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("connected")
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, age INT, first_name TEXT, last_name TEXT, email TEXT NOT NULL);")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

	db.Exec("INSERT INTO users (age, email, first_name, last_name) VALUES (30, 'jon@calhoun.io', 'Jonathan', 'Calhoun');")
	db.Exec("INSERT INTO users (age, email, first_name, last_name) VALUES (52, 'bob@smith.io', 'Bob', 'Smith');")
	db.Exec("INSERT INTO users (age, email, first_name, last_name) VALUES (15, 'jerryjr123@gmail.com', 'Jerry', 'Seinfeld');")

	sqlStatement := `SELECT id, email FROM users WHERE first_name="Jonathan";`
	var email string
	var id int
	//// Replace 3 with an ID from your database or another random
	//// value to test the no rows use case.
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&id, &email); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(id, email)
	default:
		panic(err)
	}
}
