package k8s

import (
	"ethbaas/internal/model"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type KV map[string]interface{}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(proj *model.Project) error {
	if err := p.genProjHome(proj); err != nil {
		return err
	}

	if err := p.genNameSpace(proj); err != nil {
		return err
	}

	if err := p.genConfigMap(proj); err != nil {
		return err
	}
	if err := p.genPv(proj); err != nil {
		return err
	}
	if err := p.genPvc(proj); err != nil {
		return err
	}
	if err := p.genDeploy(proj); err != nil {
		return err
	}
	if err := p.genSvc(proj); err != nil {
		return err
	}
	return nil
}

// generate proj home dir
func (p *Parser) genProjHome(proj *model.Project) error {
	return os.MkdirAll(proj.Home(), os.ModePerm)
}

// generate pv yaml
func (p *Parser) genPv(proj *model.Project) error {
	for i := 0; i < proj.NodeCount; i++ {
		pv := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolume",
			"metadata": map[string]interface{}{
				"name": fmt.Sprintf("ethbaas-node%d-pv", i),
				"labels": map[string]interface{}{
					"app": fmt.Sprintf("ethbaas-node%d-pv", i),
				},
			},
			"spec": map[string]interface{}{
				"capacity": map[string]interface{}{
					"storage": "10Gi",
				},
				"accessModes": []string{"ReadWriteMany"},
				"hostPath": map[string]interface{}{
					"path": fmt.Sprintf("/media/ethbaas/%s/node%d", proj.Name, i),
					"type": "DirectoryOrCreate",
				},
			},
		}
		composeBytes, err := yaml.Marshal(&pv)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(proj.PvFile(i), composeBytes, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// generate pvc yaml
func (p *Parser) genPvc(proj *model.Project) error {
	for i := 0; i < proj.NodeCount; i++ {
		pvc := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"namespace": proj.NS(),
				"name":      fmt.Sprintf("node%d-pv", i),
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"storage": "10Gi",
					},
				},
				"accessModes": []string{"ReadWriteMany"},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": fmt.Sprintf("ethbaas-node%d-pv", i),
					},
				},
			},
		}
		composeBytes, err := yaml.Marshal(&pvc)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(proj.PvcFile(i), composeBytes, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// generate deployment yaml
func (p *Parser) genDeploy(proj *model.Project) error {
	for i := 0; i < proj.NodeCount; i++ {
		nodeLabel := fmt.Sprintf("ethbaas-node%d", i)
		deploy := map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"namespace": proj.NS(),
				"name":      fmt.Sprintf("node%d", i),
				"labels": map[string]interface{}{
					"node": nodeLabel,
				},
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"node": nodeLabel,
					},
				},
				"replicas":        1,
				"minReadySeconds": 10,
				"strategy": map[string]interface{}{
					"type": "Recreate",
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"node": nodeLabel,
						},
					},
					"spec": KV{
						"affinity": KV{
							"nodeAffinity": KV{
								"requiredDuringSchedulingIgnoredDuringExecution": KV{
									"nodeSelectorTerms": []KV{
										{
											"matchExpressions": []KV{
												{
													"key":      "ethbaas_node",
													"operator": "In",
													"values":   []string{fmt.Sprintf("node%d", i)},
												},
											},
										},
									},
								},
							},
						},
						"initContainers": []KV{
							{
								"name":    "gen-genesis",
								"image":   "docker.io/ethereum/client-go:v1.10.18",
								"command": []string{"geth", "init", "/genesis.json"},
								"args":    []string{"--datadir=/chaindata"},
								"volumeMounts": []KV{
									{
										"name":      "chaindata",
										"mountPath": "/chaindata",
									},
									{
										"name":      "configmap",
										"mountPath": "/genesis.json",
										"subPath":   "genesis.json",
									},
								},
							},
						},
						"containers": []KV{
							{
								"name":  fmt.Sprintf("node%d", i),
								"image": "docker.io/ethereum/client-go:v1.10.18",
								"ports": []KV{
									{"containerPort": 8545},
									{"containerPort": 8546},
									{"containerPort": 30303, "protocol": "TCP"},
									{"containerPort": 30303, "protocol": "UDP"},
								},
								"args": []string{
									"--datadir=/chaindata",
									"--networkid=1874",
									"--mine",
									"--miner.threads=1",
									"--miner.etherbase=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
									"--http",
									"--http.api=admin,web3,eth,net,debug,personal",
									"--http.corsdomain=*",
									"--http.addr=0.0.0.0",
									"--rpc.allow-unprotected-txs",
									"--gcmode=archive",
									"--nodiscover",
								},
								"volumeMounts": []KV{
									{
										"name":      "chaindata",
										"mountPath": "/chaindata",
									},
								},
							},
						},
						"volumes": []KV{
							{
								"name": "chaindata",
								"persistentVolumeClaim": KV{
									"claimName": fmt.Sprintf("node%d-pv", i),
								},
							},
							{
								"name": "configmap",
								"configMap": KV{
									"name": "configmap",
								},
							},
							{
								"name": "run",
								"hostPath": KV{
									"path": "/var/run",
								},
							},
						},
					},
				},
			},
		}
		composeBytes, err := yaml.Marshal(&deploy)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(proj.DeployFile(i), composeBytes, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// generate service yaml
func (p *Parser) genSvc(proj *model.Project) error {
	for i := 0; i < proj.NodeCount; i++ {
		svc := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"namespace": proj.NS(),
				"name":      fmt.Sprintf("node%d", i),
				"labels": map[string]interface{}{
					"node": fmt.Sprintf("node%d", i),
				},
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"node": fmt.Sprintf("ethbaas-node%d", i),
				},
				"type": "NodePort",
				"ports": []interface{}{
					map[string]interface{}{
						"name":       "http",
						"targetPort": 8545,
						"port":       8545,
						"protocol":   "TCP",
						"nodePort":   proj.FirstNodePort + int32(i),
					},
					map[string]interface{}{
						"name":       "websoket",
						"targetPort": 8546,
						"port":       8546,
						"protocol":   "TCP",
					},
					map[string]interface{}{
						"name":       "p2p1",
						"targetPort": 30303,
						"port":       30303,
						"protocol":   "TCP",
					},
					map[string]interface{}{
						"name":       "p2p2",
						"targetPort": 30303,
						"port":       30303,
						"protocol":   "UDP",
					},
				},
			},
		}
		composeBytes, err := yaml.Marshal(&svc)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(proj.SvcFile(i), composeBytes, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// generate namespace yaml
func (p *Parser) genNameSpace(proj *model.Project) error {
	ns := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": proj.NS(),
		},
	}
	composeBytes, err := yaml.Marshal(&ns)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(proj.NsFile(), composeBytes, os.ModePerm)
}

// generate configmap yaml
func (p *Parser) genConfigMap(proj *model.Project) error {
	genesis, err := ioutil.ReadFile("genesis/genesis.json")
	if err != nil {
		return err
	}
	cm := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"namespace": proj.NS(),
			"name":      "configmap",
		},
		"data": map[string]interface{}{
			"genesis.json": string(genesis),
		},
	}
	composeBytes, err := yaml.Marshal(&cm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(proj.CmFile(), composeBytes, os.ModePerm)
}
