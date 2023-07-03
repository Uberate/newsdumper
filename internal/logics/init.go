package logics

import "news/pkg/staged"

const (
	FlowStartKind     = "flow-start"
	FlowGetNewsKind   = "flow-get-news"
	FlowSplitNewsKind = "flow-split-news"
	FlowSendNewsKind  = "flow-send-news"

	V1Version = "v1"

	NewsFlow  = "news-flow"
	Instances = "instance"
)

var (
	CallTime        = staged.Scopes(NewsFlow, "call-time")
	LoggerInstance  = staged.Scopes(NewsFlow, Instances, "logger")
	GetterInstances = staged.Scopes(NewsFlow, Instances, "getters")
	HookerInstances = staged.Scopes(NewsFlow, Instances, "hookers")
	NewsKey         = staged.Scopes(NewsFlow, "news")
	GroupKey        = staged.Scopes(NewsFlow, "groups")
	GroupedNewsKey  = staged.Scopes(NewsFlow, "grouped_news")
)

func init() {
	staged.StageFactory.Registry(FlowStartKind, V1Version, staged.AbsStaged(FlowStartKind, V1Version,
		StartStaged, true, true))
	staged.StageFactory.Registry(FlowGetNewsKind, V1Version, staged.AbsStaged(FlowGetNewsKind, V1Version,
		GetterStaged, true, true))
	staged.StageFactory.Registry(FlowSplitNewsKind, V1Version, staged.AbsStaged(FlowSplitNewsKind, V1Version,
		SplitStage, true, true))
	staged.StageFactory.Registry(FlowSendNewsKind, V1Version, staged.AbsStaged(FlowSendNewsKind, V1Version,
		HookerStaged, true, true))
}
