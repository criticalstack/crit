package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// owner: @ChrisRx
	// alpha: v1.0
	Krustlet featuregate.Feature = "Krustlet"

	// owner: @ChrisRx
	// alpha: v1.0
	LocalRegistry featuregate.Feature = "LocalRegistry"

	// owner: @ChrisRx
	// alpha: v1.0
	MachineAPI featuregate.Feature = "MachineAPI"

	// owner: @ChrisRx
	// alpha: v1.0
	FixCgroupMounts featuregate.Feature = "FixCgroupMounts"
)

func init() {
	// Every feature should be initiated here:
	runtime.Must(MutableGates.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		Krustlet:        {Default: false, PreRelease: featuregate.Alpha},
		LocalRegistry:   {Default: false, PreRelease: featuregate.Alpha},
		MachineAPI:      {Default: false, PreRelease: featuregate.Alpha},
		FixCgroupMounts: {Default: false, PreRelease: featuregate.Alpha},
	}))
}
