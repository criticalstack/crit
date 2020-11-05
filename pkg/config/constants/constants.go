package constants

const (
	DefaultClusterDomain            = "cluster.local"
	DefaultClusterName              = "crit"
	DefaultKubeAPIServerBindPort    = 6443
	DefaultKubeDir                  = "/etc/kubernetes"
	DefaultPodSubnet                = "10.253.0.0/16"
	DefaultServiceSubnet            = "10.254.0.0/16"
	DefaultBootstrapServerBindPort  = 8080
	DefaultHealthcheckProxyBindPort = 6444

	DefaultBootstrapServerVersion  = "0.3.0"
	DefaultCoreDNSVersion          = "1.6.9"
	DefaultHealthcheckProxyVersion = "0.1.0"
	DefaultPauseImageVersion       = "3.3"

	KubeAPIServerImage         = "k8s.gcr.io/kube-apiserver"
	KubeControllerManagerImage = "k8s.gcr.io/kube-controller-manager"
	KubeSchedulerImage         = "k8s.gcr.io/kube-scheduler"
	KubeProxyImage             = "k8s.gcr.io/kube-proxy"
	PauseImage                 = "k8s.gcr.io/pause"

	CoreDNSImage              = "docker.io/coredns/coredns"
	CritBootstrapServerImage  = "docker.io/criticalstack/bootstrap-server"
	CritHealthCheckProxyImage = "docker.io/criticalstack/healthcheck-proxy"
)

type ContainerRuntime string

const (
	Containerd ContainerRuntime = "containerd"
	Docker     ContainerRuntime = "docker"
	CRIO       ContainerRuntime = "crio"
)

func (cr ContainerRuntime) CRISocket() string {
	switch cr {
	case Containerd:
		return "unix:///var/run/containerd/containerd.sock"
	case Docker:
		return "unix:///var/run/docker.sock"
	case CRIO:
		return "unix:///var/run/crio/crio.sock"
	default:
		panic("unknown container runtime")
	}
}
