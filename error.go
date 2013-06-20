package sphero

import (
	"errors"
)

var (
	NotImplementedError       = errors.New("This feature is not yet implemented")
	GeneralError              = errors.New("General, non-specific error")
	ChecksumError             = errors.New("Checksum failure")
	CommandFragmentError      = errors.New("Received command fragment")
	UnknownCommandError       = errors.New("Unknown command ID")
	UnsupportedCommandError   = errors.New("Command currently unsupported")
	BadMessageFormatError     = errors.New("Bad message format")
	InvalidParametersError    = errors.New("Parameter value(s) invalid")
	FailedExecuteCommandError = errors.New("Failed to execute command")
	UnknownDeviceError        = errors.New("Unknown device ID")
	PowerTooLowError          = errors.New("Voltage too low for reflash operation")
	IllegalPageError          = errors.New("Illegal page number provided")
	FlashFailError            = errors.New("Page did not reprogram correctly")
	ApplicationCorruptError   = errors.New("Main application corrupt")
	MessageTimeoutError       = errors.New("Message state machine timed out")
	UnknownError              = errors.New("Unkown error")
)
