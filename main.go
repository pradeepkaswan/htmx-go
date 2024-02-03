package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

var id = 0
type Contact struct {
	Id int
	Name string 
	Email string 
}

func newContact(name, email string) Contact {
	id++ 
	return Contact{
		Id : id,
		Name: name,
		Email: email,
	}
}

type Contacts []Contact // type alias for contact slice

type Data struct {
	Contacts Contacts
}

func (d *Data) indexOf(id int) int {
	for i, contact := range d.Contacts {
		if contact.Id == id {
			return i
		}
	}
	return -1
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

	x.Static("/images", "images")
	x.Static("/css", "css")

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

		c.HTML(http.StatusOK, "form", newFormData())
		c.HTML(http.StatusOK, "oob-contact", contact)
	})

	x.DELETE("/contacts/:id", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		index := page.Data.indexOf(id)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}

		page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)

		c.Status(http.StatusOK)
	})

	// Run the server
	fmt.Println("Server is running at port 6969")
	x.Run(":6969")
}