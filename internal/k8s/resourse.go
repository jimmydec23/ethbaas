package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Apply(file string) error {
	clientset, dd, err := getClient()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func Delete(file string) error {
	clientset, dd, err := getClient()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}
		if err := dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func GetPods(namespace string) ([]Pod, error) {
	clientset, _, err := getClient()
	if err != nil {
		return nil, err
	}
	podList, err := clientset.CoreV1().
		Pods(namespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	svcList, err := clientset.CoreV1().
		Services(namespace).
		List(context.Background(), metav1.ListOptions{})

	svcs := map[string]Pod{}
	for _, svc := range svcList.Items {
		ports := []int32{}
		for _, p := range svc.Spec.Ports {
			ports = append(ports, p.NodePort)
		}
		s := Pod{
			Service:    svc.Name,
			NodePorts:  ports,
			PodSvcBind: svc.Spec.Selector["node"],
		}
		svcs[s.PodSvcBind] = s
	}

	pods := []Pod{}
	for _, p := range podList.Items {
		ports := []int32{}
		for _, port := range p.Spec.Containers[0].Ports {
			ports = append(ports, port.ContainerPort)
		}
		podSvcBind := p.ObjectMeta.Labels["node"]
		po, ok := svcs[podSvcBind]
		if !ok {
			return nil, fmt.Errorf("Pod service binding failed: %s", p.ObjectMeta.Name)
		}
		po.Name = p.ObjectMeta.Name
		po.Status = string(p.Status.Phase)
		po.Ports = ports
		pods = append(pods, po)
	}
	return pods, nil
}

func GetSvcs(namespace string) ([]v1.Service, error) {
	clientset, _, err := getClient()
	if err != nil {
		return nil, err
	}

	svcList, err := clientset.CoreV1().
		Services(namespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return svcList.Items, nil
}

func getClient() (*kubernetes.Clientset, dynamic.Interface, error) {
	home := homedir.HomeDir()
	configPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	di, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientSet, di, err
}
