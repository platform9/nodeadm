package utils

import (
	"bufio"
	"fmt"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	FILE_MODE = 0744
)

func InstallMasterComponents(rootDir, routerID, intf, vip string, masterConfig *kubeadm.MasterConfiguration) {
	DownloadArtifacts(rootDir, KUBERNETES_VERSION, CNI_VERSION)
	DownloadKubeletServiceFiles(rootDir, KUBERNETES_VERSION)
	EnableAndStartService("kubelet.service")
	ReplaceString(getKubeletServiceConf(), DEFAULT_DNS_IP, GetIPFromSubnet(masterConfig.Networking.ServiceSubnet, 10))
	writeKeepAlivedServiceFiles(routerID, intf, vip)
	EnableAndStartService("keepalived.service")
}

func InstallWorkerComponents(rootDir string) {
	DownloadArtifacts(rootDir, KUBERNETES_VERSION, CNI_VERSION)
	DownloadKubeletServiceFiles(rootDir, KUBERNETES_VERSION)
	EnableAndStartService("kubelet.service")
}

func DownloadKubeletServiceFiles(rootDir string, kuberneteVersion string) {
	baseURL := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", kuberneteVersion)
	//kubelet service
	serviceFile := filepath.Join(SYSTEMD_DIR, "kubelet.service")
	Download(serviceFile, baseURL+"kubelet.service", FILE_MODE)
	ReplaceString(serviceFile, "/usr/bin", rootDir)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(SYSTEMD_DIR, "kubelet.service.d"), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(SYSTEMD_DIR, "kubelet.service.d", "10-kubeadm.conf")
	Download(confFile, baseURL+"10-kubeadm.conf", FILE_MODE)
	ReplaceString(confFile, "/usr/bin", rootDir)
}

func writeTemplateIntoFile(tmpl, name, file string, data interface{}) {
	err := os.MkdirAll(filepath.Dir(file), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dirs for path %s with error %v\n", filepath.Dir(file), err)
	}
	f, err := os.Create(file)
	defer f.Close()
	w := bufio.NewWriter(f)
	if err != nil {
		log.Fatalf("Failed to create file %s\n", file)
	}
	t := template.Must(template.New(name).Parse(tmpl))
	t.Execute(w, data)
	w.Flush()
}

func writeKeepAlivedServiceFiles(routerID, intf, vip string) {
	kaConfFileTemplate := `
	vrrp_instance K8S_APISERVER {
		interface {{.Intf}}
		state BACKUP
		virtual_router_id {{.RouteId}}
		nopreempt
	
		authentication {
			auth_type AH
			auth_pass ourownpassword
		}
	
		virtual_ipaddress {
			{{.VIP}}
		}
	}`
	type KaConfData struct {
		RouterID, Intf, VIP string
	}
	kaConfData := KaConfData{routerID, intf, vip}
	confFile := filepath.Join(SYSTEMD_DIR, "keepalive.service.d", "keepalived.conf")
	writeTemplateIntoFile(kaConfFileTemplate, "vipConfFileTemplate", confFile, kaConfData)

	kaSvcFileTemplate := `
	[Unit]
	Description= Keepalived service
	After=network.target docker.service
	Requires=docker.service
	[Service]
	Type=simple
	ExecStart=/usr/bin/docker run --cap-add=NET_ADMIN \\
			--net=host --name vip \\
			-v {{.ConfFile}}:/usr/local/etc/keepalived/keepalived.conf \\
			-d {{.KeepAlivedImg}}
	ExecStartPre=-/usr/bin/docker kill vip
	ExecStartPre=-/usr/bin/docker rm vip
	ExecStop=/usr/bin/docker stop vip
	Restart=on-failure
	[Install]
	WantedBy=multi-user.target
	`
	type KaServiceData struct {
		ConfigFile, KeepAlivedImg string
	}
	kaServiceData := KaServiceData{confFile, KEEPALIVED_IMG}
	writeTemplateIntoFile(kaSvcFileTemplate, "kaSvcFileTemplate", filepath.Join(SYSTEMD_DIR, "keepalived.service"), kaServiceData)
}

func getKubeletServiceConf() string {
	return filepath.Join(SYSTEMD_DIR, "kubelet.service.d", "10-kubeadm.conf")
}

func DownloadKubeComponents(rootDir, version string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	//Download kubectl, kubeadm, kubelet if needed
	baseURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", version)
	Download(filepath.Join(rootDir, "kubectl"), baseURL+"kubectl", FILE_MODE)
	Download(filepath.Join(rootDir, "kubeadm"), baseURL+"kubeadm", FILE_MODE)
	Download(filepath.Join(rootDir, "kubelet"), baseURL+"kubelet", FILE_MODE)
}

func DownloadCNIPlugin(rootDir, version string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	baseURL := fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-amd64-%s.tgz", version, version)
	tmpFile := fmt.Sprintf("/tmp/cni-plugins-amd64-%s.tgz", version)
	Download(tmpFile, baseURL, FILE_MODE)
	Run(rootDir, "tar", "-xvf", tmpFile, "-C", rootDir)
}

func DownloadArtifacts(rootDir, kubernetesVersion, cniVersion string) {
	DownloadKubeComponents(rootDir, kubernetesVersion)
	DownloadCNIPlugin(CNI_BASE_DIR, cniVersion)
	//keepalived
	Run(rootDir, "docker", "pull", KEEPALIVED_IMG)
}
