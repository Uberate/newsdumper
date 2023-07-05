package hooks

import (
	"github.com/uberate/gf"
)

var HookFactory *gf.Factory[Hook]

func init() {
	HookFactory = gf.NewFactor[Hook](&EmptyHooker{})

	HookFactory.Registry(SMTPHookKind, V1Str, GeneratorSMTPHook)
	HookFactory.Registry(LarkHookKind, V1Str, GeneratorLarkHookInstance)
}
