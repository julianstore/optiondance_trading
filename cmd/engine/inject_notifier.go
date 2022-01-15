package main

import (
	"option-dance/core"
	"option-dance/service/notifier"
)

func provideNotifier(
	messageStore core.MessageStore,
	positionStore core.PositionStore,
	messageBuilder core.MessageBuilder,
) core.Notifier {
	if notify {
		return notifier.NewNotifier(
			messageStore, positionStore, messageBuilder,
		)
	}
	return notifier.NewMuteNotifier()
}
