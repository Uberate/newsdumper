package staged

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestBase(t *testing.T) {
	l := logrus.New()
	a := NewBuilder(l)
	list := a.Next("set func", func(ctx context.Context) (context.Context, error) {
		ctx = SetFromContext(ctx, "1", "1")
		return ctx, nil
	}, true, true).
		Next("get func", func(ctx context.Context) (context.Context, error) {
			var res string
			var ok bool
			res, ok = GetFromContext(ctx, "1", "2")
			fmt.Println(res, ok)
			return ctx, nil
		}, true, true).Build()

	list.RunE(context.Background())
}
