package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name string 
	Email string 
}

func newContact(name, email string) Contact {
	return Contact{
		Name: name,
		Email: email,
	}
}

type Contacts []Contact // type alias for contact slice

type Data struct {
	Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("John Doe", "johndoe@email.com"),
			newContact("Claire Doe", "clairedoe@email.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	x := gin.Default()

	page := newPage()

	renderer := newTemplate()
	x.SetHTMLTemplate(renderer.tmpl)

	x.LoadHTMLGlob("views/*.html")

	// Define routes
	x.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", page)
	})

	x.POST("/contacts", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email

			formData.Errors["email"] = "Email already exists"
			c.HTML(422, "form", formData)
			return
		}

		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		// TODO: Save the contact to the database

		c.HTML(http.StatusOK, "form", newFormData())
		c.HTML(http.StatusOK, "oob-contact", contact)
	})

	// Run the server
	fmt.Println("Server is running at port 6969")
	x.Run(":6969")
}