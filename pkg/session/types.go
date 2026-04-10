package session

type LocalAuthMode uint8

const (
	LocalAuthModeUnknown LocalAuthMode = iota
	LocalAuthModeNoPassword
	LocalAuthModePassword
)

type TLSMode string

const (
	TLSModeUnknown    TLSMode = ""
	TLSModeDisabled   TLSMode = "disabled"
	TLSModeSelfSigned TLSMode = "self-signed"
)

type DisplayRotation string

const (
	DisplayRotationUnknown  DisplayRotation = ""
	DisplayRotationNormal   DisplayRotation = "270"
	DisplayRotationInverted DisplayRotation = "90"
)

type CloudState struct {
	Connected bool
	URL       string
	AppURL    string
}

type AccessState struct {
	LocalAuthMode LocalAuthMode
	LoopbackOnly  bool
	Cloud         CloudState
	TLS           TLSMode
}

type USBConfig struct {
	VendorID     string
	ProductID    string
	SerialNumber string
	Manufacturer string
	Product      string
}

type USBDevices struct {
	AbsoluteMouse bool
	RelativeMouse bool
	Keyboard      bool
	MassStorage   bool
	SerialConsole bool
	Network       bool
}

type HardwareState struct {
	USBEmulation    *bool
	USBConfig       USBConfig
	USBDevices      USBDevices
	USBDeviceCount  int
	DisplayRotation DisplayRotation
}

type NetworkSettings struct {
	Hostname string
	IP       string
}

type NetworkState struct {
	Hostname string
	IP       string
	DHCP     *bool
}

type JigglerConfig struct {
	InactivityLimitSeconds int
	JitterPercentage       int
	ScheduleCronTab        string
	Timezone               string
}

type VersionInfo struct {
	AppVersion    string
	SystemVersion string
}

type MQTTSettings struct {
	Enabled           bool
	Broker            string
	Port              int
	Username          string
	Password          string
	BaseTopic         string
	UseTLS            bool
	TLSInsecure       bool
	EnableHADiscovery bool
	EnableActions     bool
	DebounceMs        int
}

type MQTTStatus struct {
	Connected bool
	Error     string
}

type MQTTTestResult struct {
	Success bool
	Error   string
}

type AdvancedState struct {
	DevMode      *bool
	USBEmulation *bool
	Version      VersionInfo
}
