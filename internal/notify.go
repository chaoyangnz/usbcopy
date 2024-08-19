package internal

import (
	"fmt"
	"github.com/go-toast/toast"
)

func notify(context *Context, format string, args ...any) {
	if context.UiMode {
		notification := toast.Notification{
			AppID:   "usbcopy",
			Title:   "🛈 usbcopy has something happening",
			Message: fmt.Sprintf(format, args...),
		}
		err := notification.Push()
		if err != nil {
			log("%v", err)
		}
	} else {
		log(format, args...)
	}
}
