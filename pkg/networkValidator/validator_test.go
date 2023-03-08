package networkValidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkValidatorNetworkError(t *testing.T) {
	validator := NewNetworkValidator()

	err := test_identifyNetworkError_network_error(validator, "fail")
	isLocalMode := validator.IsLocalMode()
	assert.NotNil(t, err)
	assert.Equal(t, false, isLocalMode)

	err = test_identifyNetworkError_network_error(validator, "local")
	isLocalMode = validator.IsLocalMode()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isLocalMode)

	err = test_identifyNoInternet_noSuchHost_network_error(validator, "fail")
	isLocalMode = validator.IsLocalMode()
	assert.NotNil(t, err)
	assert.Equal(t, false, isLocalMode)

	err = test_identifyNoInternet_noSuchHost_network_error(validator, "local")
	isLocalMode = validator.IsLocalMode()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isLocalMode)

	err = test_identifyNoInternet_connectionRefused_network_error(validator, "fail")
	isLocalMode = validator.IsLocalMode()
	assert.NotNil(t, err)
	assert.Equal(t, false, isLocalMode)

	err = test_identifyNoInternet_connectionRefused_network_error(validator, "local")
	isLocalMode = validator.IsLocalMode()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isLocalMode)

}

func TestNetworkValidatorOtherError(t *testing.T) {
	validator := NewNetworkValidator()

	err := test_identifyNetworkError_other_error(validator, "fail")
	isLocalMode := validator.IsLocalMode()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, isLocalMode)

	err = test_identifyNetworkError_other_error(validator, "local")
	isLocalMode = validator.IsLocalMode()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isLocalMode)

}

func test_identifyNetworkError_network_error(validator *NetworkValidator, offlineMode string) error {
	validator.SetOfflineMode(offlineMode)
	return validator.IdentifyNetworkError(errors.New("connection refused"))
}

func test_identifyNoInternet_noSuchHost_network_error(validator *NetworkValidator, offlineMode string) error {
	validator.SetOfflineMode(offlineMode)
	return validator.IdentifyNetworkError(errors.New("no such host"))
}

func test_identifyNoInternet_connectionRefused_network_error(validator *NetworkValidator, offlineMode string) error {
	validator.SetOfflineMode(offlineMode)
	return validator.IdentifyNetworkError(errors.New("tcp dial connection refused"))
}

func test_identifyNetworkError_other_error(validator *NetworkValidator, offlineMode string) error {
	validator.SetOfflineMode(offlineMode)
	return validator.IdentifyNetworkError(errors.New("mysql server is away"))
}
