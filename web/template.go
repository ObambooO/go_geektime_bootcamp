package web

import "context"

type TemplateEngine interface {
	// Render 渲染模板
	// tplName 模板的名字，按名索引
	// data 渲染页面的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}
