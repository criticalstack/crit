package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// owner: @ChrisRx
	// alpha: v1.0
	AuthProxyCA featuregate.Feature = "AuthProxyCA"

	// owner: @ChrisRx
	// alpha: v1.0
	BootstrapServer featuregate.Feature = "BootstrapServer"

	// owner: @ChrisRx
	// alpha: v1.0
	UploadETCDSecrets featuregate.Feature = "UploadETCDSecrets"
)

func init() {
	// Every feature should be initiated here:
	runtime.Must(MutableGates.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		AuthProxyCA:       {Default: false, PreRelease: featuregate.Alpha},
		BootstrapServer:   {Default: false, PreRelease: featuregate.Alpha},
		UploadETCDSecrets: {Default: true, PreRelease: featuregate.Deprecated},
	}))
}
