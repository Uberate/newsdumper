package hooks

import "news/pkg/factory"

var HookFactory *factory.Factory[Hook]

func init() {
	HookFactory = factory.NewFactor[Hook](&EmptyHooker{})

	HookFactory.Registry(SMTPHookKind, V1Str, GeneratorSMTPHook)
}
