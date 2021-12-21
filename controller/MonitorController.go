package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/service"
	"net/http"
)

type Monitor Controller

func (m Monitor) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/monitor/getList", http.MethodGet, m.GetList),
		core.NewRoute("/monitor/getTotal", http.MethodGet, m.GetTotal),
		core.NewRoute("/monitor/check", http.MethodPost, m.Check).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/monitor/add", http.MethodPost, m.Add).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/monitor/edit", http.MethodPut, m.Edit).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/monitor/toggle", http.MethodPut, m.Toggle).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
		core.NewRoute("/monitor/remove", http.MethodDelete, m.Remove).Roles(core.RoleAdmin, core.RoleManager, core.RoleGroupManager),
	}
}

func (Monitor) GetList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	monitorList, err := model.Monitor{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Monitors model.Monitors `json:"list"`
		}{Monitors: monitorList},
	}
}

func (Monitor) GetTotal(gp *core.Goploy) core.Response {
	total, err := model.Monitor{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Monitor) Check(gp *core.Goploy) core.Response {
	type ReqData struct {
		URL string `json:"url" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if err := (service.Gnet{URL: reqData.URL}.Ping()); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Message: "Connected"}
}

func (Monitor) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		URL          string `json:"url" validate:"required"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		NotifyTimes  uint16 `json:"notifyTimes" validate:"gt=0"`
		Description  string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Monitor{
		NamespaceID:  gp.Namespace.ID,
		Name:         reqData.Name,
		URL:          reqData.URL,
		Second:       reqData.Second,
		Times:        reqData.Times,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
		NotifyTimes:  reqData.NotifyTimes,
		Description:  reqData.Description,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Monitor) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		Name         string `json:"name" validate:"required"`
		URL          string `json:"url" validate:"required"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		NotifyTimes  uint16 `json:"notifyTimes" validate:"gt=0"`
		Description  string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Monitor{
		ID:           reqData.ID,
		Name:         reqData.Name,
		URL:          reqData.URL,
		Second:       reqData.Second,
		Times:        reqData.Times,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
		NotifyTimes:  reqData.NotifyTimes,
		Description:  reqData.Description,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Monitor) Toggle(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID}).ToggleState(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Monitor) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
