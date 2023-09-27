package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Note struct {
	Id      int
	Name    string
	Status  string
	Content string
}

type User struct {
	Id       int
	Name     string
	Surname  string
	Password string
	Status   string
}

type Tasks struct {
	Id      int
	Name    string
	Status  string
	Content []string
}

var notes = []Note{
	{
		Id:      1,
		Name:    "Welcome",
		Status:  "important",
		Content: "# h1 \n ## h2 \n ### h3 \n #### h4 \n ##### h5 \n ###### h6 \n **Bold** \n *Italic* \n > quoted content \n\n [link](https://joelbonetr.com/) \n\n You can use inline `code` as well as code blocks: \n ```js \n  const arr = new Array(); \n  ``` \n Lists: \n - Orange \n - Peach \n - Banana \n\n Adding images: \n\n ![JavaScript](https://www.iconninja.com/files/541/586/346/command-language-software-develop-code-programming-javascript-icon.png)",
	},
	{
		Id:      2,
		Name:    "Readme",
		Status:  "study",
		Content: ``,
	},
}

var users = []User{
	{
		Id:       1,
		Name:     "Timur",
		Surname:  "Ayushiev",
		Password: "",
		Status:   "moderator",
	},
	{
		Id:       2,
		Name:     "Varvara",
		Surname:  "Talankina",
		Password: "",
		Status:   "redactor",
	},
}

var tasks = []Tasks{
	{
		Id:      1,
		Name:    "Task #1",
		Status:  "Important",
		Content: []string{"Learn Golang", "Work in MOEX"},
	},
	{
		Id:     2,
		Name:   "Task #2",
		Status: "Unimportant",
		Content: []string{
			"Buy vegetables",
			"Buy groceries",
		},
	},
	{
		Id:     3,
		Name:   "Welcome",
		Status: "Importnant",
		Content: []string{
			"MOEX работа в четверг, 11:00",
			"MOEX подписание документов на стажировку",
		},
	},
}

type SearchResult struct {
	Title string
	Link  string
	Type  string
}

func performSearch(query string) []SearchResult {
	var results []SearchResult

	// Search through notes
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(note.Content), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Title: note.Name,
				Link:  fmt.Sprintf("/notes/md/%d", note.Id),
				Type:  "Note",
			})
		}
	}

	// Search through tasks
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Title: task.Name,
				Link:  fmt.Sprintf("/notes/todo/%d", task.Id),
				Type:  "Task",
			})
		}
	}

	return results
}

func (a *Application) StartServer() {
	log.Println("Server starting")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"css": "/styles/style.css",
		})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"css": "/styles/style.css",
		})
	})

	router.GET("/notes", func(c *gin.Context) {
		markdowns, err := a.repository.GetAllNotes(1)
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "notes.tmpl", gin.H{
			"css":   "/styles/style.css",
			"Notes": markdowns,
			"Tasks": tasks,
		})
	})

	router.GET("/notes/md/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatal(err)
		}

		markdown, err := a.repository.GetMarkdownById(uint(id))
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "notesById.tmpl", gin.H{
			"css":   "/styles/style.css",
			"Notes": notes,
			"Tasks": tasks,
			"Note":  markdown,
			"Users": users,
		})
	})

	router.POST("/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("q", ""))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = a.repository.DeleteMarkdownById(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		markdowns, err := a.repository.GetAllNotes(1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "notes.tmpl", gin.H{
			"css":   "/styles/style.css",
			"Notes": markdowns,
			"Tasks": tasks,
		})
	})

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
