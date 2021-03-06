package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/kgosse/shop-back/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Product)
// DB Table: Plural (products)
// Resource: Plural (Products)
// Path: Plural (/products)
// View Template Folder: Plural (/templates/products/)

// ProductsResource is the resource for the Product model
type ProductsResource struct {
	buffalo.Resource
}

// List gets all Products. This function is mapped to the path
// GET /products
func (v ProductsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	products := &models.Products{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Products from the DB
	if err := q.All(products); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(products))
}

// Show gets the data for one Product. This function is mapped to
// the path GET /products/{product_id}
func (v ProductsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Product
	product := &models.Product{}

	// To find the Product the parameter product_id is used.
	if err := tx.Find(product, c.Param("product_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, product))
}

// New renders the form for creating a new Product.
// This function is mapped to the path GET /products/new
func (v ProductsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Product{}))
}

// Create adds a Product to the DB. This function is mapped to the
// path POST /products
func (v ProductsResource) Create(c buffalo.Context) error {
	// Allocate an empty Product
	product := &models.Product{}

	// Bind product to the html form elements
	if err := c.Bind(product); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(product)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, product))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Product was created successfully")

	// and redirect to the products index page
	return c.Render(201, r.Auto(c, product))
}

// Edit renders a edit form for a Product. This function is
// mapped to the path GET /products/{product_id}/edit
func (v ProductsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Product
	product := &models.Product{}

	if err := tx.Find(product, c.Param("product_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, product))
}

// Update changes a Product in the DB. This function is mapped to
// the path PUT /products/{product_id}
func (v ProductsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Product
	product := &models.Product{}

	if err := tx.Find(product, c.Param("product_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Product to the html form elements
	if err := c.Bind(product); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(product)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, product))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Product was updated successfully")

	// and redirect to the products index page
	return c.Render(200, r.Auto(c, product))
}

// Destroy deletes a Product from the DB. This function is mapped
// to the path DELETE /products/{product_id}
func (v ProductsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Product
	product := &models.Product{}

	// To find the Product the parameter product_id is used.
	if err := tx.Find(product, c.Param("product_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(product); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Product was destroyed successfully")

	// Redirect to the products index page
	return c.Render(200, r.Auto(c, product))
}
