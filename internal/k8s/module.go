package k8s

type Pod struct {
	Name       string
	Status     string
	Ports      []int32
	Service    string
	NodePorts  []int32
	PodSvcBind string
}
