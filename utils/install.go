package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	FILE_MODE = 0744
)

func InstallMasterComponents(config *Configuration) {
	DownloadArtifacts()
	EnableAndStartService("kubelet.service")
	ReplaceString(getKubeletServiceConf(), DEFAULT_DNS_IP, GetIPFromSubnet(config.MasterConfiguration.Networking.ServiceSubnet, 10))
	writeKeepAlivedServiceFiles(config.VipConfiguration)
	EnableAndStartService("keepalived.service")
}

func InstallWorkerComponents() {
	DownloadArtifacts()
	EnableAndStartService("kubelet.service")
}

func DownloadKubeletServiceFiles(kuberneteVersion string) {
	baseURL := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", kuberneteVersion)
	//kubelet service
	serviceFile := filepath.Join(SYSTEMD_DIR, "kubelet.service")
	Download(serviceFile, baseURL+"kubelet.service", FILE_MODE)
	ReplaceString(serviceFile, "/usr/bin", BASE_DIR)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(SYSTEMD_DIR, "kubelet.service.d"), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(SYSTEMD_DIR, "kubelet.service.d", "10-kubeadm.conf")
	Download(confFile, baseURL+"10-kubeadm.conf", FILE_MODE)
	ReplaceString(confFile, "/usr/bin", BASE_DIR)
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

func writeKeepAlivedServiceFiles(config VIPConfiguration) {
	log.Printf("Vip configuration as parsed from the file %v\n", config)
	if len(config.IP) == 0 {
		ip, err := netutil.ChooseHostInterface()
		if err != nil {
			log.Fatalf("Failed to get default interface with err %v", err)
		}
		config.IP = ip.String()
	}

	if len(config.NetworkInterface) == 0 {
		cmdStr := "route | grep '^default' | grep -o '[^ ]*$'"
		cmd := exec.Command("bash", "-c", cmdStr)
		bytes, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to get default interface with err %v", err)
		}
		config.NetworkInterface = strings.Trim(string(bytes), "\n ")
	}

	if config.RouterID == 0 {
		config.RouterID = DEFAULT_ROUTER_ID
	}
	kaConfFileTemplate :=
		`vrrp_instance K8S_APISERVER {
	interface {{.NetworkInterface}}
	state MASTER
	virtual_router_id {{.RouterID}}
	nopreempt
	authentication {
		auth_type AH
		auth_pass ourownpassword
	}
	virtual_ipaddress {
		{{.IP}}
	}
}
`
	confFile := filepath.Join(SYSTEMD_DIR, "keepalived.conf")
	writeTemplateIntoFile(kaConfFileTemplate, "vipConfFileTemplate", confFile, config)

	kaSvcFileTemplate := `
[Unit]
Description= Keepalived service
After=network.target docker.service
Requires=docker.service
[Service]
Type=simple
ExecStart=/usr/bin/docker run --cap-add=NET_ADMIN \
		--net=host --name vip \
		-v {{.ConfigFile}}:/usr/local/etc/keepalived/keepalived.conf \
		{{.KeepAlivedImg}}
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
	CreateSymLinks(KUBE_DIR, BASE_DIR, true)

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
	CreateSymLinks(CNI_DIR, CNI_BASE_DIR, true)
}

func DownloadNetworkConfig() {
	os.MkdirAll(CONF_DIR, FILE_MODE)
	url := fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/kube-flannel.yml", FLANNEL_VERSION)
	file := filepath.Join(CONF_DIR, "flannel.yaml")
	Download(file, url, FILE_MODE)
}

func DownloadArtifacts() {
	DownloadKubeComponents(KUBE_DIR, KUBERNETES_VERSION)
	DownloadCNIPlugin(CNI_DIR, CNI_VERSION)
	DownloadKubeletServiceFiles(KUBERNETES_VERSION)
	DownloadNetworkConfig()
	//keepalived
	Run("", "docker", "pull", KEEPALIVED_IMG)
}
