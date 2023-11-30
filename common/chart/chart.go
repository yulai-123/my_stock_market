package chart

import (
	"context"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/model/do"
	"net/http"
	"sort"
	"strings"
)

type Chart struct {
	line      *charts.Line
	DataList  [][]*do.ResultData
	WithMoney bool
}

func NewChart(ctx context.Context, title string) *Chart {
	chart := &Chart{}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove",
		}),
	)

	chart.line = line
	return chart
}

func (c *Chart) AddResultData(resultDataList []*do.ResultData) {
	c.DataList = append(c.DataList, resultDataList)
}

func (c *Chart) Render(ctx context.Context) {
	xMap := make(map[string]bool)
	for _, resultDataList := range c.DataList {
		for _, resultData := range resultDataList {
			xMap[resultData.TradeDate] = true
		}
	}
	xList := util.MapKeysString(xMap)
	sort.Slice(xList, func(i, j int) bool {
		return strings.Compare(xList[i], xList[j]) <= 0
	})
	c.line.SetXAxis(xList)

	for _, resultDataList := range c.DataList {
		tsCode := resultDataList[0].TSCode
		for _, x := range xList {
			flag := false
			for _, resultData := range resultDataList {
				if x == resultData.TradeDate {
					flag = true
					break
				}
			}
			if !flag {
				resultDataList = append(resultDataList, &do.ResultData{
					Result:    0,
					Money:     0,
					TradeDate: x,
					TSCode:    tsCode,
				})
			}
		}
		sort.Slice(resultDataList, func(i, j int) bool {
			return strings.Compare(resultDataList[i].TradeDate, resultDataList[j].TradeDate) <= 0
		})

		items := make([]opts.LineData, 0)
		for _, resultData := range resultDataList {
			items = append(items, opts.LineData{
				Value: resultData.Result,
			})
		}
		c.line.AddSeries(tsCode+"result", items)

		if c.WithMoney {
			items := make([]opts.LineData, 0)
			for _, resultData := range resultDataList {
				items = append(items, opts.LineData{
					Value: resultData.Money,
				})
			}
			c.line.AddSeries(tsCode+"money", items)
		}
	}

	c.line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c.line.Render(writer)
	})

	logrus.Infof("生成完成")
	http.ListenAndServe(":8081", nil)
}
