package v1alpha2

import (
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	kubeproxyconfigv1alpha1 "k8s.io/kube-proxy/config/v1alpha1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"k8s.io/utils/pointer"

	"github.com/criticalstack/crit/pkg/config/constants"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ControlPlaneConfiguration(obj *ControlPlaneConfiguration) {
	if obj == nil {
		obj = &ControlPlaneConfiguration{}
	}
	if obj.ClusterName == "" {
		obj.ClusterName = constants.DefaultClusterName
	}
	if obj.PodSubnet == "" {
		obj.PodSubnet = constants.DefaultPodSubnet
	}
	if obj.ServiceSubnet == "" {
		obj.ServiceSubnet = constants.DefaultServiceSubnet
	}
	if obj.KubeAPIServerConfiguration.BindPort == 0 {
		obj.KubeAPIServerConfiguration.BindPort = constants.DefaultKubeAPIServerBindPort
	}
	if obj.CoreDNSVersion == "" {
		obj.CoreDNSVersion = constants.DefaultCoreDNSVersion
	}
	obj.CritBootstrapServerConfiguration.Version = strings.TrimPrefix(obj.CritBootstrapServerConfiguration.Version, "v")
	if obj.CritBootstrapServerConfiguration.Version == "" {
		obj.CritBootstrapServerConfiguration.Version = constants.DefaultBootstrapServerVersion
	}
	if obj.CritBootstrapServerConfiguration.BindPort == 0 {
		obj.CritBootstrapServerConfiguration.BindPort = constants.DefaultBootstrapServerBindPort
	}
	obj.KubeAPIServerConfiguration.HealthcheckProxyVersion = strings.TrimPrefix(obj.KubeAPIServerConfiguration.HealthcheckProxyVersion, "v")
	if obj.KubeAPIServerConfiguration.HealthcheckProxyVersion == "" {
		obj.KubeAPIServerConfiguration.HealthcheckProxyVersion = constants.DefaultHealthcheckProxyVersion
	}
	if obj.KubeAPIServerConfiguration.HealthcheckProxyBindPort == 0 {
		obj.KubeAPIServerConfiguration.HealthcheckProxyBindPort = constants.DefaultHealthcheckProxyBindPort
	}
	if obj.KubeProxyConfiguration.Config == nil {
		obj.KubeProxyConfiguration.Config = &kubeproxyconfigv1alpha1.KubeProxyConfiguration{}
		SetDefaults_KubeProxyConfiguration(obj.KubeProxyConfiguration.Config)
	}
	SetDefaults_NodeConfiguration(&obj.NodeConfiguration)
}

func SetDefaults_WorkerConfiguration(obj *WorkerConfiguration) {
	if obj == nil {
		obj = &WorkerConfiguration{}
	}
	if obj.ClusterName == "" {
		obj.ClusterName = constants.DefaultClusterName
	}
	SetDefaults_NodeConfiguration(&obj.NodeConfiguration)
}

func SetDefaults_NodeConfiguration(obj *NodeConfiguration) {
	obj.KubernetesVersion = strings.TrimPrefix(obj.KubernetesVersion, "v")
	if obj.KubeDir == "" {
		obj.KubeDir = constants.DefaultKubeDir
	}
	if obj.ContainerRuntime == "" {
		obj.ContainerRuntime = constants.Containerd
	}
	if obj.KubeletConfiguration == nil {
		obj.KubeletConfiguration = &kubeletconfigv1beta1.KubeletConfiguration{}
		SetDefaults_KubeletConfiguration(obj.KubeletConfiguration)
	}
}

const (
	DefaultIPTablesMasqueradeBit = 14
	DefaultIPTablesDropBit       = 15
)

var (
	// DefaultEvictionHard includes default options for hard eviction.
	DefaultEvictionHard = map[string]string{
		"memory.available":  "100Mi",
		"nodefs.available":  "10%",
		"nodefs.inodesFree": "5%",
		"imagefs.available": "15%",
	}

	// Refer to [Node Allocatable](https://git.k8s.io/community/contributors/design-proposals/node/node-allocatable.md) doc for more information.
	DefaultNodeAllocatableEnforcement = []string{"pods"}

	zeroDuration = metav1.Duration{}
)

func init() {
	_ = kubeletconfigv1beta1.AddToScheme(clientsetscheme.Scheme)
}

func SetDefaults_KubeletConfiguration(obj *kubeletconfigv1beta1.KubeletConfiguration) {
	if obj.ClusterDomain == "" {
		obj.ClusterDomain = constants.DefaultClusterDomain
	}
	if obj.SyncFrequency == zeroDuration {
		obj.SyncFrequency = metav1.Duration{Duration: 1 * time.Minute}
	}
	if obj.FileCheckFrequency == zeroDuration {
		obj.FileCheckFrequency = metav1.Duration{Duration: 20 * time.Second}
	}
	if obj.HTTPCheckFrequency == zeroDuration {
		obj.HTTPCheckFrequency = metav1.Duration{Duration: 20 * time.Second}
	}
	if obj.Address == "" {
		obj.Address = "0.0.0.0"
	}
	if obj.Port == 0 {
		obj.Port = 10250
	}
	if obj.Authentication.Anonymous.Enabled == nil {
		obj.Authentication.Anonymous.Enabled = pointer.BoolPtr(false)
	}
	if obj.Authentication.Webhook.Enabled == nil {
		obj.Authentication.Webhook.Enabled = pointer.BoolPtr(true)
	}
	if obj.Authentication.Webhook.CacheTTL == zeroDuration {
		obj.Authentication.Webhook.CacheTTL = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.Authorization.Mode == "" {
		obj.Authorization.Mode = kubeletconfigv1beta1.KubeletAuthorizationModeWebhook
	}
	if obj.Authorization.Webhook.CacheAuthorizedTTL == zeroDuration {
		obj.Authorization.Webhook.CacheAuthorizedTTL = metav1.Duration{Duration: 5 * time.Minute}
	}
	if obj.Authorization.Webhook.CacheUnauthorizedTTL == zeroDuration {
		obj.Authorization.Webhook.CacheUnauthorizedTTL = metav1.Duration{Duration: 30 * time.Second}
	}
	if obj.RegistryPullQPS == nil {
		obj.RegistryPullQPS = pointer.Int32Ptr(5)
	}
	if obj.RegistryBurst == 0 {
		obj.RegistryBurst = 10
	}
	if obj.EventRecordQPS == nil {
		obj.EventRecordQPS = pointer.Int32Ptr(5)
	}
	if obj.EventBurst == 0 {
		obj.EventBurst = 10
	}
	if obj.EnableDebuggingHandlers == nil {
		obj.EnableDebuggingHandlers = pointer.BoolPtr(true)
	}
	if obj.HealthzPort == nil {
		obj.HealthzPort = pointer.Int32Ptr(10248)
	}
	if obj.HealthzBindAddress == "" {
		obj.HealthzBindAddress = "127.0.0.1"
	}
	if obj.OOMScoreAdj == nil {
		obj.OOMScoreAdj = pointer.Int32Ptr(-999)
	}
	if obj.StreamingConnectionIdleTimeout == zeroDuration {
		obj.StreamingConnectionIdleTimeout = metav1.Duration{Duration: 4 * time.Hour}
	}
	if obj.NodeStatusReportFrequency == zeroDuration {
		// For backward compatibility, NodeStatusReportFrequency's default value is
		// set to NodeStatusUpdateFrequency if NodeStatusUpdateFrequency is set
		// explicitly.
		if obj.NodeStatusUpdateFrequency == zeroDuration {
			obj.NodeStatusReportFrequency = metav1.Duration{Duration: 5 * time.Minute}
		} else {
			obj.NodeStatusReportFrequency = obj.NodeStatusUpdateFrequency
		}
	}
	if obj.NodeStatusUpdateFrequency == zeroDuration {
		obj.NodeStatusUpdateFrequency = metav1.Duration{Duration: 10 * time.Second}
	}
	if obj.NodeLeaseDurationSeconds == 0 {
		obj.NodeLeaseDurationSeconds = 40
	}
	if obj.ImageMinimumGCAge == zeroDuration {
		obj.ImageMinimumGCAge = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.ImageGCHighThresholdPercent == nil {
		// default is below docker's default dm.min_free_space of 90%
		obj.ImageGCHighThresholdPercent = pointer.Int32Ptr(85)
	}
	if obj.ImageGCLowThresholdPercent == nil {
		obj.ImageGCLowThresholdPercent = pointer.Int32Ptr(80)
	}
	if obj.VolumeStatsAggPeriod == zeroDuration {
		obj.VolumeStatsAggPeriod = metav1.Duration{Duration: time.Minute}
	}
	if obj.CgroupsPerQOS == nil {
		obj.CgroupsPerQOS = pointer.BoolPtr(true)
	}
	if obj.CgroupDriver == "" {
		obj.CgroupDriver = "cgroupfs"
	}
	if obj.CPUManagerPolicy == "" {
		obj.CPUManagerPolicy = "none"
	}
	if obj.CPUManagerReconcilePeriod == zeroDuration {
		// Keep the same as default NodeStatusUpdateFrequency
		obj.CPUManagerReconcilePeriod = metav1.Duration{Duration: 10 * time.Second}
	}
	if obj.TopologyManagerPolicy == "" {
		obj.TopologyManagerPolicy = kubeletconfigv1beta1.NoneTopologyManagerPolicy
	}
	if obj.RuntimeRequestTimeout == zeroDuration {
		obj.RuntimeRequestTimeout = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.HairpinMode == "" {
		obj.HairpinMode = kubeletconfigv1beta1.PromiscuousBridge
	}
	if obj.MaxPods == 0 {
		obj.MaxPods = 110
	}
	// default nil or negative value to -1 (implies node allocatable pid limit)
	if obj.PodPidsLimit == nil || *obj.PodPidsLimit < int64(0) {
		obj.PodPidsLimit = pointer.Int64Ptr(-1)
	}
	if obj.ResolverConfig == "" {
		obj.ResolverConfig = "/etc/resolv.conf"
	}
	if obj.CPUCFSQuota == nil {
		obj.CPUCFSQuota = pointer.BoolPtr(true)
	}
	if obj.CPUCFSQuotaPeriod == nil {
		obj.CPUCFSQuotaPeriod = &metav1.Duration{Duration: 100 * time.Millisecond}
	}
	if obj.MaxOpenFiles == 0 {
		obj.MaxOpenFiles = 1000000
	}
	if obj.ContentType == "" {
		obj.ContentType = "application/vnd.kubernetes.protobuf"
	}
	if obj.KubeAPIQPS == nil {
		obj.KubeAPIQPS = pointer.Int32Ptr(5)
	}
	if obj.KubeAPIBurst == 0 {
		obj.KubeAPIBurst = 10
	}
	if obj.SerializeImagePulls == nil {
		obj.SerializeImagePulls = pointer.BoolPtr(true)
	}
	if obj.EvictionHard == nil {
		obj.EvictionHard = DefaultEvictionHard
	}
	if obj.EvictionPressureTransitionPeriod == zeroDuration {
		obj.EvictionPressureTransitionPeriod = metav1.Duration{Duration: 5 * time.Minute}
	}
	if obj.EnableControllerAttachDetach == nil {
		obj.EnableControllerAttachDetach = pointer.BoolPtr(true)
	}
	if obj.MakeIPTablesUtilChains == nil {
		obj.MakeIPTablesUtilChains = pointer.BoolPtr(true)
	}
	if obj.IPTablesMasqueradeBit == nil {
		obj.IPTablesMasqueradeBit = pointer.Int32Ptr(DefaultIPTablesMasqueradeBit)
	}
	if obj.IPTablesDropBit == nil {
		obj.IPTablesDropBit = pointer.Int32Ptr(DefaultIPTablesDropBit)
	}
	if obj.FailSwapOn == nil {
		obj.FailSwapOn = pointer.BoolPtr(true)
	}
	if obj.ContainerLogMaxSize == "" {
		obj.ContainerLogMaxSize = "10Mi"
	}
	if obj.ContainerLogMaxFiles == nil {
		obj.ContainerLogMaxFiles = pointer.Int32Ptr(5)
	}
	if obj.ConfigMapAndSecretChangeDetectionStrategy == "" {
		obj.ConfigMapAndSecretChangeDetectionStrategy = kubeletconfigv1beta1.WatchChangeDetectionStrategy
	}
	if obj.EnforceNodeAllocatable == nil {
		obj.EnforceNodeAllocatable = DefaultNodeAllocatableEnforcement
	}
}

func init() {
	_ = kubeproxyconfigv1alpha1.AddToScheme(clientsetscheme.Scheme)
}

func SetDefaults_KubeProxyConfiguration(obj *kubeproxyconfigv1alpha1.KubeProxyConfiguration) {
	if obj.BindAddress == "" {
		obj.BindAddress = "0.0.0.0"
	}
	if obj.ClientConnection.Kubeconfig == "" {
		obj.ClientConnection.Kubeconfig = "/var/lib/kube-proxy/kubeconfig.conf"
	}
	if obj.ClientConnection.ContentType == "" {
		obj.ClientConnection.ContentType = "application/vnd.kubernetes.protobuf"
	}
	if obj.ClientConnection.QPS == 0 {
		obj.ClientConnection.QPS = 5.0
	}
	if obj.ClientConnection.Burst == 0 {
		obj.ClientConnection.Burst = 10
	}
	if obj.ConfigSyncPeriod == zeroDuration {
		obj.ConfigSyncPeriod = metav1.Duration{Duration: 15 * time.Minute}
	}
	if obj.Conntrack.MaxPerCore == nil {
		obj.Conntrack.MaxPerCore = pointer.Int32Ptr(32 * 1024)
	}
	if obj.Conntrack.Min == nil {
		obj.Conntrack.Min = pointer.Int32Ptr(128 * 1024)
	}
	if obj.Conntrack.TCPEstablishedTimeout == nil {
		obj.Conntrack.TCPEstablishedTimeout = &metav1.Duration{Duration: 24 * time.Hour}
	}
	if obj.Conntrack.TCPCloseWaitTimeout == nil {
		obj.Conntrack.TCPCloseWaitTimeout = &metav1.Duration{Duration: 1 * time.Hour}
	}
	if obj.FeatureGates == nil {
		obj.FeatureGates = make(map[string]bool)
	}
	if obj.HealthzBindAddress == "" {
		obj.HealthzBindAddress = "0.0.0.0:10256"
	}
	if obj.IPTables.MasqueradeBit == nil {
		obj.IPTables.MasqueradeBit = pointer.Int32Ptr(14)
	}
	if obj.IPTables.SyncPeriod == zeroDuration {
		obj.IPTables.SyncPeriod = metav1.Duration{Duration: 30 * time.Second}
	}
	if obj.IPVS.SyncPeriod == zeroDuration {
		obj.IPVS.SyncPeriod = metav1.Duration{Duration: 30 * time.Second}
	}
	if obj.MetricsBindAddress == "" {
		obj.MetricsBindAddress = "127.0.0.1:10249"
	}
	if obj.Mode == "" {
		obj.Mode = "iptables"
	}
	if obj.OOMScoreAdj == nil {
		obj.OOMScoreAdj = pointer.Int32Ptr(-999)
	}
	if obj.UDPIdleTimeout == zeroDuration {
		obj.UDPIdleTimeout = metav1.Duration{Duration: 250 * time.Millisecond}
	}
}
