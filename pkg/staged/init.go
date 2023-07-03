package staged

import (
	"context"
	"github.com/uberate/gf"
)

var StageFactory *gf.Factory[Staged]

var CommonKind = "common"
var ScopeKind = "scope"

var V1Version = "v1"

func init() {
	StageFactory = gf.NewFactor[Staged](&CommonStaged{
		BaseEntity: gf.NewBaseEntityGenerator("", "")(""),
	})

	StageFactory.Registry(CommonKind, V1Version, GenerateCommonStaged)
	StageFactory.Registry(ScopeKind, V1Version, AbsStaged(ScopeKind, V1Version, func(ctx context.Context) (context.Context, error) {
		ctx = NewShared(ctx)
		return ctx, nil
	}, true, true))
}
