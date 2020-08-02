package xman

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type RouteType uint8

const (
	RouteTypeMenu RouteType = iota + 1
	RouteTypeApi
)

// 注册路由
type Router struct {
	RouterType RouteType // 路由类型
	Name       string    // 名称
	Path       string    // 路径
	Action     string    // 行为，get、post...
	h          []gin.HandlerFunc
	Sub        []R // 子路由
	Hide       bool
}

func (r Router) getRouteTypeName() string {
	switch r.RouterType {
	case RouteTypeMenu:
		return "菜单"
	case RouteTypeApi:
		return "API"
	}
	return "Unknown type"
}

func (r Router) ToString() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s：%s，path:%s\n", r.getRouteTypeName(), r.Name, r.Path))
	for _, item := range r.Sub {
		b.WriteByte('\t')
		b.WriteString(item.ToString())
	}
	return b.String()
}

type R interface {
	Register(gin.IRouter)
	ToString() string
}

func NewHideMenu(name string, path string, sub []R, handler ...gin.HandlerFunc) *Menu {
	m := NewMenu(name, path, sub, handler...)
	m.Hide = true
	return m
}

func NewMenu(name string, path string, sub []R, handler ...gin.HandlerFunc) *Menu {
	m := &Menu{}
	m.RouterType = RouteTypeMenu
	m.Name = name
	m.Path = path
	m.Sub = sub
	m.h = handler
	return m
}

type Menu struct {
	Router
}

func (m Menu) Register(r gin.IRouter) {
	r1 := r.Group(m.Path, m.h...)
	for _, item := range m.Sub {
		item.Register(r1)
	}
}

func NewHideApi(action string, name string, path string, handler ...gin.HandlerFunc) *Api {
	m := NewApi(action, name, path, handler...)
	m.Hide = true
	return m
}

func NewApi(action string, name string, path string, handler ...gin.HandlerFunc) *Api {
	m := &Api{}
	m.RouterType = RouteTypeApi
	m.Name = name
	m.Path = path
	m.h = handler
	m.Action = action
	return m
}

type Api struct {
	Router
}

func (a Api) Register(r gin.IRouter) {
	switch strings.ToUpper(a.Action) {
	case "GET":
		r.GET(a.Path, a.h...)
	case "POST":
		r.POST(a.Path, a.h...)
	case "DELETE":
		r.DELETE(a.Path, a.h...)
	case "PATCH":
		r.PATCH(a.Path, a.h...)
	case "PUT":
		r.PUT(a.Path, a.h...)
	case "OPTIONS":
		r.OPTIONS(a.Path, a.h...)
	case "HEAD":
		r.HEAD(a.Path, a.h...)
	}
}
