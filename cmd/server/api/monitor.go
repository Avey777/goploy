// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/monitor"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/model"
	"net/http"
	"time"
)

type Monitor API

func (m Monitor) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/monitor/getList", http.MethodGet, m.GetList).Permissions(config.ShowMonitorPage),
		server.NewRoute("/monitor/check", http.MethodPost, m.Check).LogFunc(middleware.AddOPLog),
		server.NewRoute("/monitor/add", http.MethodPost, m.Add).Permissions(config.AddMonitor).LogFunc(middleware.AddOPLog),
		server.NewRoute("/monitor/edit", http.MethodPut, m.Edit).Permissions(config.EditMonitor).LogFunc(middleware.AddOPLog),
		server.NewRoute("/monitor/toggle", http.MethodPut, m.Toggle).Permissions(config.EditMonitor).LogFunc(middleware.AddOPLog),
		server.NewRoute("/monitor/remove", http.MethodDelete, m.Remove).Permissions(config.DeleteMonitor).LogFunc(middleware.AddOPLog),
	}
}

func (Monitor) GetList(gp *server.Goploy) server.Response {
	monitorList, err := model.Monitor{NamespaceID: gp.Namespace.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Monitors model.Monitors `json:"list"`
		}{Monitors: monitorList},
	}
}

func (Monitor) Check(gp *server.Goploy) server.Response {
	type ReqData struct {
		Type    int           `json:"type" validate:"oneof=1 2 3 4 5"`
		Items   []string      `json:"items" validate:"required"`
		Timeout time.Duration `json:"timeout"`
		Process string        `json:"process"`
		Script  string        `json:"script"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if err := (monitor.Monitor{
		Type:    reqData.Type,
		Items:   reqData.Items,
		Timeout: reqData.Timeout,
		Process: reqData.Process,
		Script:  reqData.Script,
	}.Check()); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Message: "Connected"}
}

func (Monitor) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		Type         int    `json:"type" validate:"oneof=1 2 3 4 5"`
		Target       string `json:"target" validate:"required"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		Description  string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Monitor{
		NamespaceID:  gp.Namespace.ID,
		Name:         reqData.Name,
		Type:         reqData.Type,
		Target:       reqData.Target,
		Second:       reqData.Second,
		Times:        reqData.Times,
		SilentCycle:  reqData.SilentCycle,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
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

func (Monitor) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		Name         string `json:"name" validate:"required"`
		Type         int    `json:"type" validate:"oneof=1 2 3 4 5"`
		Target       string `json:"target" validate:"required"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		SilentCycle  int    `json:"silentCycle" validate:"required"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		Description  string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Monitor{
		ID:           reqData.ID,
		Name:         reqData.Name,
		Type:         reqData.Type,
		Target:       reqData.Target,
		Second:       reqData.Second,
		Times:        reqData.Times,
		SilentCycle:  reqData.SilentCycle,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
		Description:  reqData.Description,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Monitor) Toggle(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID    int64 `json:"id" validate:"gt=0"`
		State uint8 `json:"state" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID, State: reqData.State}).ToggleState(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Monitor) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}