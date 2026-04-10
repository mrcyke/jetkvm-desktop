package session

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
	Cloud CloudState
	TLS   TLSMode
}

type USBConfig struct {
	VendorID  string
	ProductID string
}

type HardwareState struct {
	USBEmulation    *bool
	USBConfig       USBConfig
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

type AdvancedState struct {
	DevMode      *bool
	USBEmulation *bool
	Version      VersionInfo
}
