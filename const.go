package sphero

// Start of Packet values
const (
	SOP1                     = 0xff
	SOP2_ANSWER              = 0xff
	SOP2_ASYNC               = 0xfe
	SOP2_RESET_TIMEOUT       = 0xfd
	SOP2_ASYNC_RESET_TIMEOUT = 0xfc
)

// Device IDs
const (
	DID_CORE       = 0x00
	DID_BOOTLOADER = 0x01
	DID_SPHERO     = 0x02
)

// Core Commands, DID = 00h
const (
	CMD_PING               = 0x01
	CMD_VERSION            = 0x02
	CMD_CONTROL_UART_TX    = 0x03
	CMD_SET_BT_NAME        = 0x10
	CMD_GET_BT_NAME        = 0x11
	CMD_SET_AUTO_RECONNECT = 0x12
	CMD_GET_AUTO_RECONNECT = 0x13
	CMD_GET_PWR_STATE      = 0x20
	CMD_SET_PWR_NOTIFY     = 0x21
	CMD_SLEEP              = 0x22
	GET_POWER_TRIPS        = 0x23
	SET_POWER_TRIPS        = 0x24
	SET_INACTIVE_TIMER     = 0x25
	CMD_GOTO_BL            = 0x30
	CMD_RUN_L1_DIAGS       = 0x40
	CMD_RUN_L2_DIAGS       = 0x41
	CMD_CLEAR_COUNTERS     = 0x42
	CMD_ASSIGN_TIME        = 0x50
	CMD_POLL_TIMES         = 0x51
)

// Bootloader Commands, DID = 01h
const (
	BEGIN_REFLASH         = 0x02
	HERE_IS_PAGE          = 0x03
	LEAVE_BOOTLOADER      = 0x04
	IS_PAGE_BLANK         = 0x05
	CMD_ERASE_USER_CONFIG = 0x06
)

// Sphero Commands, DID = 02h
const (
	CMD_SET_CAL                 = 0x01
	CMD_SET_STABILIZ            = 0x02
	CMD_SET_ROTATION_RATE       = 0x03
	CMD_SET_BALL_REG_WEBSITE    = 0x04
	CMD_GET_BALL_REG_WEBSITE    = 0x05
	CMD_REENABLE_DEMO           = 0x06
	CMD_GET_CHASSIS_ID          = 0x07
	CMD_SET_CHASSIS_ID          = 0x08
	CMD_SELF_LEVEL              = 0x09
	CMD_SET_VDL                 = 0x0a
	CMD_SET_DATA_STREAMING      = 0x11
	CMD_SET_COLLISION_DET       = 0x12
	CMD_LOCATOR                 = 0x13
	CMD_SET_ACCELERO            = 0x14
	CMD_READ_LOCATOR            = 0x15
	CMD_SET_RGB_LED             = 0x20
	CMD_SET_BACK_LED            = 0x21
	CMD_GET_RGB_LED             = 0x22
	CMD_ROLL                    = 0x30
	CMD_BOOST                   = 0x31
	CMD_MOVE                    = 0x32
	CMD_SET_RAW_MOTORS          = 0x33
	CMD_SET_MOTION_TO           = 0x34
	CMD_SET_OPTIONS_FLAG        = 0x35
	CMD_GET_OPTIONS_FLAG        = 0x36
	CMD_SET_TEMP_OPTIONS_FLAG   = 0x37
	CMD_GET_TEMP_OPTIONS_FLAG   = 0x38
	CMD_GET_CONFIG_BLK          = 0x40
	CMD_SET_DEVICE_MODE         = 0x42
	CMD_SET_CFG_BLOCK           = 0x43
	CMD_GET_DEVICE_MODE         = 0x44
	CMD_RUN_MACRO               = 0x50
	CMD_SAVE_TEMP_MACRO         = 0x51
	CMD_SAVE_MACRO              = 0x52
	CMD_INIT_MACRO_EXECUTIVE    = 0x54
	CMD_ABORT_MACRO             = 0x55
	CMD_MACRO_STATUS            = 0x56
	CMD_SET_MACRO_PARAM         = 0x57
	CMD_APPEND_TEMP_MACRO_CHUNK = 0x58
	CMD_ERASE_ORBBAS            = 0x60
	CMD_APPEND_FRAG             = 0x61
	CMD_EXEC_ORBBAS             = 0x62
	CMD_ABORT_ORBBAS            = 0x63
	CMD_ANSWER_INPUT            = 0x64
)

// Message Response Codes
const (
	ORBOTIX_RSP_CODE_OK           = 0x00 // Command succeeded
	ORBOTIX_RSP_CODE_EGEN         = 0x01 // General, non-specific error
	ORBOTIX_RSP_CODE_ECHKSUM      = 0x02 // Received checksum failure
	ORBOTIX_RSP_CODE_EFRAG        = 0x03 // Received command fragment
	ORBOTIX_RSP_CODE_EBAD_CMD     = 0x04 // Unknown command ID
	ORBOTIX_RSP_CODE_EUNSUPP      = 0x05 // Command currently unsupported
	ORBOTIX_RSP_CODE_EBAD_MSG     = 0x06 // Bad message format
	ORBOTIX_RSP_CODE_EPARAM       = 0x07 // Parameter value(s) invalid
	ORBOTIX_RSP_CODE_EEXEC        = 0x08 // Failed to execute command
	ORBOTIX_RSP_CODE_EBAD_DID     = 0x09 // Unknown Device ID
	ORBOTIX_RSP_CODE_POWER_NOGOOD = 0x31 // Voltage too low for reflash operation
	ORBOTIX_RSP_CODE_PAGE_ILLEGAL = 0x32 // Illegal page number provided
	ORBOTIX_RSP_CODE_FLASH_FAIL   = 0x33 // Page did not reprogram correctly
	ORBOTIX_RSP_CODE_MA_CORRUPT   = 0x34 // Main Application corrupt
	ORBOTIX_RSP_CODE_MSG_TIMEOUT  = 0x35 // Msg state machine timed out
)

// Async Message Id Code
const (
	ID_POWER_NOTIFICATIONS         = 0x01 // Power notifications
	ID_LEVEL_1_DIAGNOSTIC_RESPONSE = 0x02 // Level 1 Diagnostic response
	ID_SENSOR_DATA_STREAMING       = 0x03 // Sensor data streaming
	ID_CONFIG_BLOCK_CONTENTS       = 0x04 // Config block contents
	ID_PRE_SLEEP_WARNING           = 0x05 // Pre-sleep warning (10 sec)
	ID_MACRO_MARKERS               = 0x06 // Macro markers
	ID_COLLISION_DETECTED          = 0x07 // Collision detected
	ID_ORBBAS_PRINT                = 0x08 // orbBasic PRINT message
	ID_ORBBAS_ERROR_ASCII          = 0x09 // orbBasic error message, ASCII
	ID_ORBBAS_ERROR_BINARY         = 0x0a // orbBasic error message, binary
	ID_SELF_LEVEL_RESULT           = 0x0b // Self Level Result
	ID_GYRO_AXIS_LIMIT_EXCEEDED    = 0x0c // Gyro axis limit exceeded (FW ver 3.10 and later)
)

// Battery
const (
	BATTERY_CHARGING = 0x01
	BATTERY_OK       = 0x02
	BATTERY_LOW      = 0x03
	BATTERY_CRITICAL = 0x04
)

// Power State Masks
const (
	POWER_MASK_RECVER       = 0x00 // Record version code – the following definition is for 01h
	POWER_MASK_STATE        = 0x01 // High-level state of the power system as concluded by the power manager: 01h = Battery Charging, 02h = Battery OK, 03h = Battery Low, 04h = Battery Critical
	POWER_MASK_BATT_VOLTAGE = 0x02 // Current battery voltage scaled in 100ths of a volt; 02EFh would be 7.51 volts (unsigned 16-bit value)
	POWER_MASK_NUM_CHARGES  = 0x04 // Number of battery recharges in the life of this Sphero (unsigned 16-bit value)
	POWER_MASK_TIMESINCECHG = 0x06 // Seconds awake since last recharge (unsigned 16-bit value)
)

// Data Streaming Masks

// MASK1
const (
	ACCEL_AXIS_X_RAW = 0x80000000
	ACCEL_AXIS_Y_RAW = 0x40000000
	ACCEL_AXIS_Z_RAW = 0x20000000

	GYRO_AXIS_X_RAW = 0x10000000
	GYRO_AXIS_Y_RAW = 0x08000000
	GYRO_AXIS_Z_RAW = 0x04000000

	MOTOR_RIGHT_EMF_RAW = 0x00400000
	MOTOR_LEFT_EMF_RAW  = 0x00200000

	MOTOR_LEFT_PWM_RAW  = 0x00100000
	MOTOR_RIGHT_PWM_RAW = 0x00080000

	IMU_PITCH_ANGLE_FILTERED = 0x00040000
	IMU_ROLL_ANGLE_FILTERED  = 0x00020000
	IMU_YAW_ANGLE_FILTERED   = 0x00010000

	ACCEL_AXIS_X_FILTERED = 0x00008000
	ACCEL_AXIS_Y_FILTERED = 0x00004000
	ACCEL_AXIS_Z_FILTERED = 0x00002000

	GYRO_AXIS_X_FILTERED = 0x00001000
	GYRO_AXIS_Y_FILTERED = 0x00000800
	GYRO_AXIS_Z_FILTERED = 0x00000400

	MOTOR_RIGHT_EMF_FILTERED = 0x00000040
	MOTOR_LEFT_EMF_FILTERED  = 0x00000020
)

// MASK2
const (
	QUATERNION_Q0 = 0x80000000
	QUATERNION_Q1 = 0x40000000
	QUATERNION_Q2 = 0x20000000
	QUATERNION_Q3 = 0x10000000
	ODOMETER_X    = 0x08000000
	ODOMETER_Y    = 0x04000000
	ACCEL_ONE     = 0x02000000
	VELOCITY_X    = 0x01000000
	VELOCITY_Y    = 0x00800000
)
