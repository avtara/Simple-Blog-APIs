package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type (
	blogModel struct {
		gorm.Model
		Title    string `json:"title"`
		Content  string `json:"content"`
		Tags     string `json:"tags"`
		Status   int    `json:"status"`
		AuthorID int    `json:"authorid"`
	}

	transformedBlog struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Tags      string    `json:"tags"`
		Status    int       `json:"status"`
		AuthorID  string    `json:"authorid"`
		CreatedAt time.Time `json:"createdat"`
		UpdatedAt time.Time `json:"updateat"`
	}
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "xxxxxxxxx") // enter your uri database
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&blogModel{})
}

func main() {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/blogs/", createBlog)
		v1.GET("/blogs/", fetchAllBlog)
		v1.GET("/blogs/:id", fetchSingleBlog)
		v1.PUT("/blogs/:id", updateBlog)
		v1.DELETE("/blogs/:id", deleteBlog)

	}
	router.Run()
}

func createBlog(c *gin.Context) {
	authorID, _ := strconv.Atoi(c.PostForm("authorid"))
	status, _ := strconv.Atoi(c.PostForm("status"))
	blog := blogModel{Title: c.PostForm("title"), Content: c.PostForm("content"), Tags: c.PostForm("tags"), Status: status, AuthorID: authorID}

	db.Save(&blog)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Success creating blog!", "resourceId": blog.ID})
}

func fetchAllBlog(c *gin.Context) {
	var blogs []blogModel
	var _blogs []transformedBlog
	var author string

	db.Find(&blogs)

	if len(blogs) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Blogs Available !"})
		return
	}

	for _, item := range blogs {
		if item.AuthorID == 1 {
			author = "Muhammad Avtara Khrisna"
		} else {
			author = "Admin"
		}
		_blogs = append(_blogs, transformedBlog{ID: item.ID, Title: item.Title, Content: item.Content, Tags: item.Tags, Status: item.Status, AuthorID: author, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _blogs})
}

func fetchSingleBlog(c *gin.Context) {
	var blog blogModel
	var author string
	blogID := c.Param("id")

	db.First(&blog, blogID)

	if blog.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Blogs Available !"})
		return
	}

	if blog.AuthorID == 1 {
		author = "Muhammad Avtara Khrisna"
	} else {
		author = "Admin"
	}

	_blog := transformedBlog{ID: blog.ID, Title: blog.Title, Content: blog.Content, Tags: blog.Tags, Status: blog.Status, AuthorID: author, CreatedAt: blog.CreatedAt, UpdatedAt: blog.UpdatedAt}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _blog})
}

func updateBlog(c *gin.Context) {
	var blog blogModel
	blogID := c.Param("id")

	db.First(&blog, blogID)

	if blog.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Blogs Available !"})
		return
	}

	db.Model(&blog).Update("title", c.PostForm("title"))
	db.Model(&blog).Update("content", c.PostForm("content"))
	db.Model(&blog).Update("tags", c.PostForm("tags"))
	status, _ := strconv.Atoi(c.PostForm("status"))
	db.Model(&blog).Update("status", status)
	authorID, _ := strconv.Atoi(c.PostForm("authorid"))
	db.Model(&blog).Update("author_id", authorID)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success updating blog !"})
}

func deleteBlog(c *gin.Context) {
	var blog blogModel
	blogID := c.Param("id")

	db.First(&blog, blogID)

	if blog.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Blogs Available !"})
		return
	}

	db.Delete(&blog)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success deleting blog !"})
}
