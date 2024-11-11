package keyvaluefactory

import (
	"errors"

	"github.com/randyardiansyah25/libpkg/util/env"
)

func GetStore() (Store, error) {
	driverName := env.GetString("app.keyvalue_driver")
	if driverName == DRIVER_REDIS {
		return newRedisImpl(), nil
	} else {
		return nil, errors.New("unimplemented keyvalue store driver")
	}
}
