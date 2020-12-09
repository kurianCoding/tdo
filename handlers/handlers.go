package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"todo/models"

	logr "todo/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	logrus "github.com/sirupsen/logrus"
)

const (
	SERVER     = "internal server error"
	DONE       = "success"
	NOEXIST    = "task does not exist"
	JSON       = "bad json"
	ID_INVALID = "id invalid"
)

var (
	log *logr.Logger
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
	logr.InitLogger(Log)
	log = logr.New()
	fmt.Println("init")
}

type CreateTask struct {
	ID          string
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	SubTask     []models.Task
	Parent      string
	Type        string
}

// CreateTaskHandler create new task
// TODO: add validation for time

// CreateTaskHandler creates Task
// CreateTaskHandler godoc
// @Summary Create a Task
// @Description add by json Task
// @Tags CreateTask
// @Param task body CreateTask true "Create Task"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Error 400 string
// @Error 500 error
// @Router /task [post]
func CreateTaskHandler(c *gin.Context) {
	// create new task
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	type CreateTaskPl struct {
		CreateTask
	}
	var t1 CreateTask
	err = json.Unmarshal(bytes, &t1)
	t := models.Task(t1) // converting user struct to database ready struct

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, JSON)
		return
	}
	u := uuid.New()
	t.ID = u.String()
	res, err := models.CreateTask(t)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.Debug(res.GetUids())
	c.JSON(http.StatusOK, map[string]string{"id": u.String()})
	return
}

// UpdateTaskHandler updates task
// UpdateTaskHandler create new task
// TODO: add validation for time

// UpdateTaskHandler creates Task
// UpdateTaskHandler godoc
// @Summary Update a Task
// @Description add by json Task
// @Tags CreateTask
// @Param task body CreateTask true "Create Task"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Error 400 string
// @Error 500 error
// @Router /task/{taskId} [put]
func UpdateTaskHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ID_INVALID)
		return
	}
	// TODO: check if id is alpha numeric
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var t models.Task
	t.ID = id // settting to param in case not specified
	err = json.Unmarshal(payload, &t)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSON)
		return
	}
	res, err := models.UpdateTask(t)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, SERVER)
	}
	log.Debug(res)
	c.JSON(http.StatusOK, DONE)

	return
}

// GetAllHandler gets all tasks currently in the system
// TODO: add filters
func GetAllHandler(c *gin.Context) {
	res, err := models.GetAll()
	// c.JSON(http.StatusOK, fmt.Sprintf("%s", res.Json)
	// return
	log.Info(res)
	type payload struct {
		Data []models.Task `json:"getTasks"`
	}

	var pl payload
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(res.GetJson(), &pl)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Debug((pl.Data))
	c.JSON(http.StatusOK, pl.Data)
	return
}

// GetTaskHandler get task by id gets task by parameter id
func GetTaskHandler(c *gin.Context) {
	id := string(c.Param("id"))
	if id == "" {
		log.Error(fmt.Errorf("blank id"))
		c.JSON(http.StatusBadRequest, ID_INVALID)
		return
	}

	log.Debug(id)
	res, err := models.GetTask(id)
	log.Debug(res)
	var parse map[string]interface{}
	err = json.Unmarshal(res.Json, &parse)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(parse["getTasks"].([]interface{})) == 0 {
		c.JSON(http.StatusNotFound, NOEXIST)
		return
	}
	c.JSON(http.StatusOK, parse["getTasks"])
	return
}

// DeleteTaskHandler delete task
func DeleteTaskHandler(c *gin.Context) {
	id := c.Param("id")

	res, err := models.GetUID(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, ID_INVALID)
		return
	}
	log.Debug(res)

	var pars map[string]interface{}

	err = json.Unmarshal(res.Json, &pars)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if len(pars["getTasks"].([]interface{})) == 0 {
		c.JSON(http.StatusNotFound, NOEXIST)
		return
	}
	log.Info(pars)
	var ids []string
	for _, val := range pars["getTasks"].([]interface{}) {
		for _, val1 := range val.(map[string]interface{}) {
			ids = append(ids, string(val1.(string)))
		}

	}
	log.Debug(ids)
	for _, id := range ids {
		res, err = models.DeleteEdges(id, "task.task", "task.description", "task.title", "task.time")
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	res, err = models.DeleteTask(ids)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Debug(res)

	log.Debug(res)
	c.JSON(http.StatusOK, DONE)
	return

}

// SubTaskHandler put given task as subtask
// TODO: convert this into upsert mutation
func SubTaskHandler(c *gin.Context) {
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var t models.Task
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad json")
		return
	}
	u := uuid.New()
	t.ID = u.String()

	if t.Parent == "" {
		c.JSON(http.StatusBadRequest, "parent id empty")
		return
	}

	uid, err := models.GetUID(t.Parent)

	var pars map[string]interface{}

	err = json.Unmarshal(uid.Json, &pars)
	if len(pars["getTasks"].([]interface{})) == 0 {
		c.JSON(http.StatusNotFound, NOEXIST)
		return
	}

	parentUID := string(pars["getTasks"].([]interface{})[0].(map[string]interface{})["uid"].(string))
	log.Debug(parentUID)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.Debug(t)
	res, err := models.CreateTask(t)
	ts := res.GetUids()
	var tp string
	for _, val := range ts {
		tp = val
	}
	childUID := tp
	log.Debug(childUID)

	res, err = models.Link(parentUID, childUID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.Debug(res)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, DONE)
	return
}
