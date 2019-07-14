package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Project struct
type Project struct{}

// GetList project list
func (project *Project) GetList(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	projectList, err := model.Project{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Project: projectList}}
	response.Json(w)
}

// GetDetail project detail
func (project *Project) GetDetail(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		ProjectServers model.ProjectServers `json:"projectServerMap"`
		ProjectUsers   model.ProjectUsers   `json:"projectUserMap"`
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.Json(w)
		return
	}
	projectServersModel := model.ProjectServers{}
	if err := projectServersModel.GetServerByProjectID(uint32(id)); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectUsersModel := model.ProjectUsers{}
	if err := projectUsersModel.GetUserByProjectID(uint32(id)); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{ProjectServers: projectServersModel, ProjectUsers: projectUsersModel}}
	response.Json(w)
}

// Add one project
func (project *Project) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Name      string   `json:"name"`
		URL       string   `json:"url"`
		Path      string   `json:"path"`
		ServerIDs []uint32 `json:"serverIds"`
		UserIDs   []uint32 `json:"userIds"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectID, err := model.Project{
		Name:       reqData.Name,
		URL:        reqData.URL,
		Path:       reqData.Path,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID:  projectID,
			ServerID:   serverID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}
	err = projectServersModel.AddMany()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID:  projectID,
			UserID:     userID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}
	err = projectUsersModel.AddMany()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	response := core.Response{Message: "添加成功"}
	response.Json(w)
}

// Edit one Project
func (project *Project) Edit(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID   uint32 `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Path string `json:"path"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	err = model.Project{
		ID:         reqData.ID,
		Name:       reqData.Name,
		URL:        reqData.URL,
		Path:       reqData.Path,
		UpdateTime: time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.Json(w)
}

// Create new repository
func (project *Project) Create(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectData, err := model.Project{
		ID: reqData.ID,
	}.GetData()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	err = model.Project{
		ID:     reqData.ID,
		Status: 1,
	}.ChangeStatus()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	path := "./repository/" + projectData.Name
	repo := projectData.URL

	// clone repository async
	go func(id uint32, path, repo string) {
		projectModel := model.Project{
			ID: id,
		}
		err = os.RemoveAll(path)
		if err != nil {
			projectModel.Status = 3
			_ = projectModel.ChangeStatus()
			fmt.Println(err)
			return
		}
		cmd := exec.Command("git", "clone", repo, path)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			projectModel.Status = 3
			_ = projectModel.ChangeStatus()
			fmt.Println(err)
			return
		}
		projectModel.Status = 2
		_ = projectModel.ChangeStatus()
	}(reqData.ID, path, repo)

	response := core.Response{Message: "初始化中，请稍后"}
	response.Json(w)
}
