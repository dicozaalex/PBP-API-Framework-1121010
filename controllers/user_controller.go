package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo"
)

func GetAllUsers(c echo.Context) error {
	db := connect()
	defer db.Close()
	query := "SELECT * FROM users"

	name := c.QueryParam("name")
	age := c.QueryParam("age")
	if len(name) != 0 {
		query += " WHERE name='" + name + "'"
	}
	if len(age) != 0 {
		if len(name) != 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age='" + age + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		SendErrorResponse(c.Echo(), c, 400, "Failed query")
		return err
	}
	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			c.Echo().Logger.Fatal(err)
		} else {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		SendErrorResponse(c.Echo(), c, 400, "No data selected")
	} else {
		var response UsersResponse
		response.Status = 200
		response.Message = "success"
		response.Data = users
		err := json.NewEncoder(c.Response()).Encode(response)
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
	return nil
}

func InsertUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	name := c.FormValue("name")
	age, _ := strconv.Atoi(c.FormValue("age"))
	address := c.FormValue("address")
	email := c.FormValue("email")
	password := c.FormValue("password")
	usertype, _ := strconv.Atoi(c.FormValue("usertype"))

	_, errQuery := db.Exec("INSERT INTO users(name, age, address, email, password, usertype) VALUES(?,?,?,?,?,?)",
		name,
		age,
		address,
		email,
		password,
		usertype)

	if errQuery == nil {
		var response UserResponse
		response.Status = 200
		response.Message = "success"
		response.Data = User{Name: name, Age: age, Address: address, Email: email, Password: password, UserType: usertype}
		err := json.NewEncoder(c.Response()).Encode(response)
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	} else {
		fmt.Println(errQuery)
		SendErrorResponse(c.Echo(), c, 400, "Insert Failed!")
	}
	return nil
}

func UpdateUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	id := c.Param("id")

	name := c.FormValue("name")
	age, _ := strconv.Atoi(c.FormValue("age"))
	address := c.FormValue("address")
	email := c.FormValue("email")
	password := c.FormValue("password")

	age_str := strconv.Itoa(age)
	query := "UPDATE users SET"
	if len(name) != 0 {
		query += " name='" + name + "'"
	}
	if age != 0 {
		if len(name) != 0 {
			query += ","
		}
		query += " age='" + age_str + "'"
	}
	if len(address) != 0 {
		if age != 0 || len(name) != 0 {
			query += ","
		}
		query += " address='" + address + "'"
	}
	if len(email) != 0 {
		if len(address) != 0 || age != 0 || len(name) != 0 {
			query += ","
		}
		query += " email='" + email + "'"
	}
	if len(password) != 0 {
		if len(email) != 0 || len(address) != 0 || age != 0 || len(name) != 0 {
			query += ","
		}
		query += " password='" + password + "'"
	}
	query += " WHERE id='" + id + "'"
	_, errQuery := db.Exec(query)
	if errQuery != nil {
		SendErrorResponse(c.Echo(), c, 400, "Error update")
	}

	// select user untuk ditampilkan
	selectQuery, err := db.Prepare("SELECT * FROM users WHERE id=?")
	if err != nil {
		SendErrorResponse(c.Echo(), c, 400, "Error select")
	}

	var user User
	err = selectQuery.QueryRow(id).Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		fmt.Println(err)
		SendErrorResponse(c.Echo(), c, 400, "Error scan")
	}

	if err == nil {
		var response UserResponse
		response.Status = 200
		response.Message = "success"
		response.Data = user
		err := json.NewEncoder(c.Response()).Encode(response)
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	} else {
		fmt.Println(errQuery)
		SendErrorResponse(c.Echo(), c, 400, "Insert Failed!")
	}
	return nil
}

func DeleteUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	id := c.Param("id")
	_, errQuery := db.Exec("DELETE FROM users WHERE id=?",
		id,
	)

	if errQuery == nil {
		var response UserResponse
		response.Status = 200
		response.Message = "delete success"
		err := json.NewEncoder(c.Response()).Encode(response)
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	} else {
		fmt.Println(errQuery)
		SendErrorResponse(c.Echo(), c, 400, "Insert Failed!")
	}
	return nil
}

func SendErrorResponse(e *echo.Echo, c echo.Context, status int, message string) {
	var response ErrorResponse
	response.Status = status
	response.Message = message
	json.NewEncoder(c.Response()).Encode(response)

	e.Logger.Fatal(response)
}
