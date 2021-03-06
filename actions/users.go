package actions

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/kgosse/shop-back/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// LoginRequest represents a login form.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MyCustomClaims associates Role and jwt.StandardClaims
type MyCustomClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// Login logs in a User. This function is
// mapped to the path GET /auth/login
func (v UsersResource) Login(c buffalo.Context) error {
	var req LoginRequest
	err := c.Bind(&req)

	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	pwd := req.Password
	if len(pwd) == 0 {
		return c.Error(http.StatusBadRequest, errors.New("Invalid password"))
	}

	email := req.Email
	if checkmail.ValidateFormat(email) != nil {
		return c.Error(http.StatusBadRequest, errors.New("Invalid email"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	log.Printf("email = %v and password = %v\n", req.Email, req.Password)
	if err := tx.Where("email = ? and password = ?", req.Email, req.Password).First(user); err != nil {
		log.Println(err)
		return c.Error(http.StatusBadRequest, errors.New("Invalid credentials"))
	}

	claims := MyCustomClaims{
		"member",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    fmt.Sprintf("%s.api.shop", envy.Get("GO_ENV", "development")),
			Id:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(envy.Get("JWT_SECRET", "JesusIsGod"))
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return fmt.Errorf("could not sign token, %v", err)
	}

	return c.Render(200, r.JSON(map[string]interface{}{"token": tokenString, "user": user}))
}

// LoginAdmin logs in an admin. This function is
// mapped to the path GET admin/auth/login
func (v UsersResource) LoginAdmin(c buffalo.Context) error {
	var req LoginRequest
	err := c.Bind(&req)

	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	pwd := req.Password
	if len(pwd) == 0 {
		return c.Error(http.StatusBadRequest, errors.New("Invalid password"))
	}

	email := req.Email
	if checkmail.ValidateFormat(email) != nil {
		return c.Error(http.StatusBadRequest, errors.New("Invalid email"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	log.Printf("email = %v and password = %v\n", req.Email, req.Password)
	if err := tx.Eager().Where("email = ? and password = ?", req.Email, req.Password).First(user); err != nil {
		log.Println(err)
		return c.Error(http.StatusBadRequest, errors.New("Invalid credentials"))
	}

	// The user should be an admin
	isAdmin := false
	for _, r := range user.Roles {
		if r.Role == "admin" {
			isAdmin = true
		}
	}

	if isAdmin == false {
		return errors.WithStack(errors.New("you are not admin"))
	}

	claims := MyCustomClaims{
		"admin",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    fmt.Sprintf("%s.api.shop", envy.Get("GO_ENV", "development")),
			Id:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(envy.Get("JWT_SECRET", "JesusIsGod"))
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return fmt.Errorf("could not sign token, %v", err)
	}

	return c.Render(200, r.JSON(map[string]interface{}{"token": tokenString, "user": user}))
}

// List gets all Users. This function is mapped to the path
// GET /users
func (v UsersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	users := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Users from the DB
	if err := q.All(users); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(users))
}

// Show gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
func (v UsersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, user))
}

// New renders the form for creating a new User.
// This function is mapped to the path GET /users/new
func (v UsersResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.User{}))
}

// Create adds a User to the DB. This function is mapped to the
// path POST /users
func (v UsersResource) Create(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Bind user to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	// user.Role = nulls.NewString("member")

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	if user.Password != user.ConfirmPassword {
		return c.Render(422, r.JSON(map[string]interface{}{"error": "passwords do not match!"}))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.JSON(map[string]interface{}{"error": verrs}))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "User was created successfully")

	user.Password = ""
	user.ConfirmPassword = ""
	// and redirect to the users index page
	return c.Render(201, r.JSON(user))
}

// Edit renders a edit form for a User. This function is
// mapped to the path GET /users/{user_id}/edit
func (v UsersResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, user))
}

// Update changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
func (v UsersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind User to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, user))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "User was updated successfully")

	// and redirect to the users index page
	return c.Render(200, r.Auto(c, user))
}

// Destroy deletes a User from the DB. This function is mapped
// to the path DELETE /users/{user_id}
func (v UsersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(user); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "User was destroyed successfully")

	// Redirect to the users index page
	return c.Render(200, r.Auto(c, user))
}
