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
	SERVER         = "internal server error"
	DONE           = "success"
	NOEXIST        = "task does not exist"
	NOEXIST_PARENT = "parent task does not exist"
	JSON           = "bad json"
	ID_INVALID     = "id invalid"
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

/* CreateTaskHandler create new task
TODO: add validation for time*/

// CreateTask creates Task
// CreateTask godoc
// @Summary Create a Task
// @Description add by json Task
// @Tags CreateTask
// @Param task body models.CreateTaskPl true "Create Task"
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
	var t1 models.CreateTaskPl
	err = json.Unmarshal(bytes, &t1)
	t3 := models.CreateTaskDB{} // converting user struct to database ready struct
	t4 := models.CreateTaskPlInt(t1)
	t := models.Task{t4, t3} // converting user struct to database ready struct
	//t := models.Task(t2)        // converting user struct to database ready struct

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

// UpdateTaskHandler Updates Task
// UpdateTaskHandler godoc
// @Summary Update a Task
// @Description add by json Task
// @Tags UpdateTask
// @Param taskId path string true "Update Task"
// @Param task body models.CreateTaskPl true "Create Task"
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
	var t1 models.CreateTaskPl
	err = json.Unmarshal(payload, &t1)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSON)
		return
	}
	t3 := models.CreateTaskPlInt(t1)
	t2 := models.CreateTaskDB{}
	t := models.Task{t3, t2}

	t.ID = id // settting to param in case not specified
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

// GetAllHandler Updates Task
// GetAllHandler godoc
// @Summary Get All Tasks
// @Description Get by json Task
// @Tags GetAll
// @Accept  json
// @Produce json
// @Success 200 {object} []models.CreateTaskPl
// @Error 400 string
// @Error 500 error
// @Router /task [get]
func GetAllHandler(c *gin.Context) {
	res, err := models.GetAll()
	// c.JSON(http.StatusOK, fmt.Sprintf("%s", res.Json)
	// return
	log.Info(res)
	type GetAllPl struct {
		models.CreateTaskPl
		ID string `json:"id"`
	}
	type payload struct {
		Data []models.TaskPl `json:"getTasks"`
	}

	var pl payload
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

// GetTaskHandler get task by parameter id
// GetTask Gets Task
// GetTask godoc
// @Summary Get a Task
// @Description Get by id taskId
// @Tags GetTask
// @Param taskId path string true "Get Task"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Error 400 string
// @Error 500 error
// @Router /task/{taskId} [get]
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
// DeleteTaskHandler Updates Task
// DeleteTaskHandler godoc
// @Summary Delete a Task
// @Description delete by id Task
// @Tags DeleteTask
// @Param taskId path string true "Delete Task"
// @Produce  json
// @Success 200 {object} string
// @Error 400 string
// @Error 500 error
// @Router /task/{taskId} [delete]
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

// SubTaskHandler Updates Task
// SubTaskHandler godoc
// @Summary Update a Task
// @Description add by json subtask
// @Tags SubTask
// @Param taskId path string true "Subtask Task"
// @Param task body models.CreateTaskPl true "Create Task"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Error 400 string
// @Error 500 error
// @Router /subtask/{taskId} [post]
func SubTaskHandler(c *gin.Context) {
	id := c.Param("id")
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var t1 models.CreateTaskPl
	err = json.Unmarshal(bytes, &t1)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSON)
		return
	}

	t3 := models.CreateTaskPlInt(t1)
	t2 := models.CreateTaskDB{} // Converting User side struct to db struct
	t := models.Task{t3, t2}

	t.Parent = id //parentID is a parameter
	u := uuid.New()
	t.ID = u.String()

	if t.Parent == "" {
		// subtask cannot be added if parent does not exist
		c.JSON(http.StatusBadRequest, NOEXIST_PARENT)
		return
	}

	uid, err := models.GetUID(t.Parent) // get parents uid

	var pars map[string]interface{}

	err = json.Unmarshal(uid.Json, &pars)
	if len(pars["getTasks"].([]interface{})) == 0 {
		c.JSON(http.StatusNotFound, NOEXIST_PARENT)
		return
	}

	parentUID := string(pars["getTasks"].([]interface{})[0].(map[string]interface{})["uid"].(string))
	log.Info(parentUID)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, SERVER)
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
	log.Info(childUID)

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
