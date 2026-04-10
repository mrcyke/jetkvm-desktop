package client

import (
	"encoding/json"

	"github.com/lkarlslund/jetkvm-desktop/pkg/virtualmedia"
)

type LocalVersion struct {
	AppVersion    string `json:"appVersion"`
	SystemVersion string `json:"systemVersion"`
}

type UpdateStatus struct {
	AppUpdateAvailable    bool `json:"appUpdateAvailable"`
	SystemUpdateAvailable bool `json:"systemUpdateAvailable"`
}

type NetworkSettings struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type NetworkState struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	DHCP     bool   `json:"dhcp"`
}

type CloudState struct {
	Connected bool   `json:"connected"`
	URL       string `json:"url"`
	AppURL    string `json:"appUrl"`
}

type TLSState struct {
	Mode string `json:"mode"`
}

type USBConfig struct {
	VendorID  string `json:"vendor_id"`
	ProductID string `json:"product_id"`
}

type DisplayRotationState struct {
	Rotation string `json:"rotation"`
}

type DeveloperModeState struct {
	Enabled bool `json:"enabled"`
}

type KeysDownState struct {
	Modifier byte   `json:"modifier"`
	Keys     []byte `json:"keys"`
}

type signalingMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type offerSignalData struct {
	SD string `json:"sd"`
}

type storageUploadRequest struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type networkSettingsRequest struct {
	Settings NetworkSettings `json:"settings"`
}

type setTLSStateRequest struct {
	State TLSState `json:"state"`
}

type setDisplayRotationRequest struct {
	Params DisplayRotationState `json:"params"`
}

type setQualityRequest struct {
	Factor float64 `json:"factor"`
}

type rebootRequest struct {
	Force bool `json:"force"`
}

type wheelReportRequest struct {
	WheelY int `json:"wheelY"`
}

type storageFilesResponse struct {
	Files []virtualmedia.StorageFile `json:"files"`
}

type rawList[T any] []T

func (r *rawList[T]) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	out := make([]T, len(raw))
	for i := range raw {
		if len(raw[i]) == 0 || string(raw[i]) == "null" {
			continue
		}
		if err := json.Unmarshal(raw[i], &out[i]); err != nil {
			return err
		}
	}
	*r = out
	return nil
}
