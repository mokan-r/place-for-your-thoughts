package controller

import (
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/jmoiron/sqlx"
	"github.com/mokan-r/place-for-your-thoughts/internal/adapter/db"
	"github.com/mokan-r/place-for-your-thoughts/internal/adapter/db/postgresdb"
	"github.com/mokan-r/place-for-your-thoughts/internal/model"
	"github.com/sirupsen/logrus"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	itemsOnPage  = 3
	roundedLeft  = `rounded-l`
	roundedRight = `rounded-r`
)

type Handler struct {
	db db.Storager
}

func New(client *sqlx.DB) *Handler {
	return &Handler{db: postgresdb.New(client)}
}

func (h *Handler) StartRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	limiter := tollbooth.NewLimiter(100, nil)

	handler.GET(`/health`, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	auth := handler.Group(`/admin`, gin.BasicAuth(gin.Accounts{
		os.Getenv(`LOGIN`): os.Getenv(`PASSWORD`),
	}))
	{
		auth.GET(`/`, h.AdminPage)
		auth.POST(`/`, h.AddPost)
	}
	handler.GET(`/`, tollbooth_gin.LimitHandler(limiter), h.MainPage)
	handler.GET(`/superheroes`, tollbooth_gin.LimitHandler(limiter), h.Entry)
}

func (h *Handler) Run() {
	r := gin.Default()
	h.StartRouter(r)
	r.LoadHTMLGlob(`./resources/html/*`)
	r.Static(`/assets`, `./resources/assets`)
	err := r.Run(`:8888`)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}

func (h *Handler) AddPost(c *gin.Context) {
	var m model.Post
	c.Request.ParseForm()
	m.Name = c.Request.PostFormValue(`name`)
	m.Text = c.Request.PostFormValue(`text`)

	if err := h.db.AddPost(m); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.AdminPage(c)
}

func (h *Handler) ValidateNumber(num string) (res int, err error) {
	res, err = strconv.Atoi(num)
	return
}

func (h *Handler) MainPage(c *gin.Context) {
	if len(c.Errors) > 0 {
		ErrorResponse(c, http.StatusTooManyRequests, errors.New(`429 Too many Requests`).Error())
		return
	}
	page := c.Query(`page`)
	pageNumber, err := h.ValidateNumber(page)
	if err != nil {
		pageNumber = 1
	}
	entriesCount, err := h.db.GetEntriesCount()
	if err != nil {
		c.HTML(http.StatusNotFound, `404.html`, gin.H{})
		return
	}
	if entriesCount == 0 {
		c.HTML(http.StatusNotFound, `emptyPage.html`, gin.H{})
	}

	maxPage := int(math.Ceil(float64(entriesCount) / float64(itemsOnPage)))
	if pageNumber > maxPage || pageNumber <= 0 {
		c.HTML(http.StatusNotFound, `404.html`, gin.H{})
		return
	}

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	m, _ := h.db.GetEntriesWithOffset(itemsOnPage, (pageNumber-1)*3)
	c.HTML(http.StatusOK, `index.html`, gin.H{
		`Data`:           m,
		`EntryStart`:     fmt.Sprintf(`%d`, (pageNumber-1)*3+1),
		`EntryEnd`:       fmt.Sprintf(`%d`, (pageNumber-1)*3+len(m)),
		`EntryCount`:     fmt.Sprintf(`%d`, entriesCount),
		`NextPageButton`: GetPageButton(maxPage, pageNumber, 1),
		`PrevPageButton`: GetPageButton(maxPage, pageNumber, -1),
	})
}

func GetPageButton(max int, current int, offset int) (button template.HTML) {
	disabled := false
	if offset > 0 {
		if max == current {
			disabled = true
		}
		button = GetPageButtonHTML(current+offset, `Next`, roundedRight, disabled)
	}

	if offset < 0 {
		if current+offset <= 0 {
			disabled = true
		}
		button = GetPageButtonHTML(current+offset, `Prev`, roundedLeft, disabled)
	}

	return
}

func GetPageButtonHTML(linkPage int, text string, rounded string, disabled bool) template.HTML {
	disabledAttr := ``
	if disabled {
		disabledAttr = `disabled`
	}

	return template.HTML(fmt.Sprintf(`<button onclick="window.location.href='/?page=%d'" %v
			class="px-4
			py-2
			text-sm
			font-medium
			text-white
			bg-primary-800
			border-0
			border-primary-500
			%v
			hover:bg-primary-900
			dark:bg-primary-800
			dark:border-primary-500
			dark:text-primary-400
			dark:hover:bg-primary-700
			dark:hover:text-white
			disabled:bg-gray-400
			disabled:border-gray-400
			dark:disabled:hover:bg-gray-400
			dark:disabled:text-white">
            %s
                    </button>`, linkPage, disabledAttr, rounded, text))
}

func (h *Handler) Entry(c *gin.Context) {
	idQuery := c.Query(`id`)
	entry, err := h.db.GetEntry(idQuery)
	if err != nil {
		c.HTML(http.StatusNotFound, `404.html`, gin.H{})
		return
	}

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	md := []byte(entry.Text)
	html := markdown.ToHTML(md, parser, nil)

	c.HTML(http.StatusOK, `entry.html`, gin.H{
		`Name`: entry.Name,
		`Text`: template.HTML(html),
	})
}

func (h *Handler) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, `admin.html`, gin.H{})
}
