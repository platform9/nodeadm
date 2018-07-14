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

	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	netutil "k8s.io/apimachinery/pkg/util/net"
)

func InstallMasterComponents(config *apis.InitConfiguration) {
	PopulateCache()
	PlaceComponentsFromCache(config.Networking)
	// write 20-nodeadm.conf drop-in
	// cluster-dns, cluster-domain, max-pods
	EnableAndStartService("kubelet.service")
	writeKeepAlivedServiceFiles(config.VIPConfiguration)
	EnableAndStartService("keepalived.service")
}

func InstallWorkerComponents(config *apis.JoinConfiguration) {
	PopulateCache()
	PlaceComponentsFromCache(config.Networking)
	EnableAndStartService("kubelet.service")
}

func PlaceComponentsFromCache(netConfig apis.Networking) {
	placeKubeComponents()
	placeCNIPlugin()
	placeAndModifyKubeletServiceFile()
	placeAndModifyKubeadmKubeletSystemdDropin(netConfig)
	placeNetworkConfig()
}

func placeAndModifyKubeletServiceFile() {
	serviceFile := filepath.Join(constants.SYSTEMD_DIR, "kubelet.service")
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubelet.service"), serviceFile)
	ReplaceString(serviceFile, "/usr/bin", constants.BASE_INSTALL_DIR)
}

func placeAndModifyKubeadmKubeletSystemdDropin(netConfig apis.Networking) {
	err := os.MkdirAll(filepath.Join(constants.SYSTEMD_DIR, "kubelet.service.d"), constants.EXECUTE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(constants.SYSTEMD_DIR, "kubelet.service.d", constants.KubeadmKubeletSystemdDropinFilename)
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, constants.KubeadmKubeletSystemdDropinFilename), confFile)
	ReplaceString(confFile, "/usr/bin", constants.BASE_INSTALL_DIR)

	dnsIP, err := kubeadmconstants.GetDNSIP(netConfig.ServiceSubnet)
	if err != nil {
		log.Fatalf("Failed to derive DNS IP from service subnet %q: %v", netConfig.ServiceSubnet, err)
	}
	ReplaceString(confFile, constants.DEFAULT_DNS_IP, dnsIP.String())
}

func placeKubeComponents() {
	err := os.MkdirAll(constants.KUBE_VERSION_INSTALL_DIR, constants.EXECUTE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", constants.KUBE_VERSION_INSTALL_DIR, err)
	}
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubectl"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubectl"))
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubeadm"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubeadm"))
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubelet"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubelet"))
	CreateSymLinks(constants.KUBE_VERSION_INSTALL_DIR, constants.BASE_INSTALL_DIR, true)
}

func placeCNIPlugin() {
	tmpFile := fmt.Sprintf("cni-plugins-amd64-%s.tgz", constants.CNI_VERSION)
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.CNI_DIR_NAME, tmpFile), filepath.Join("/tmp", tmpFile))
	if _, err := os.Stat(constants.CNI_VERSION_INSTALL_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(constants.CNI_VERSION_INSTALL_DIR, constants.EXECUTE)
		if err != nil {
			log.Fatalf("Failed to create dir %s with error %v\n", constants.CNI_VERSION_INSTALL_DIR, err)
		}
		Run("", "tar", "-xvf", filepath.Join("/tmp", tmpFile), "-C", constants.CNI_VERSION_INSTALL_DIR)
		CreateSymLinks(constants.CNI_VERSION_INSTALL_DIR, constants.CNI_BASE_DIR, true)
	}

}

func placeNetworkConfig() {
	os.MkdirAll(constants.CONF_INSTALL_DIR, constants.EXECUTE)
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.FLANNEL_DIR_NAME, "kube-flannel.yml"), filepath.Join(constants.CONF_INSTALL_DIR, "flannel.yaml"))
}

func writeTemplateIntoFile(tmpl, name, file string, data interface{}) {
	err := os.MkdirAll(filepath.Dir(file), constants.READ)
	if err != nil {
		log.Fatalf("Failed to create dirs for path %s with error %v", filepath.Dir(file), err)
	}
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Failed to create file %q: %v", file, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	t := template.Must(template.New(name).Parse(tmpl))
	t.Execute(w, data)
	err = w.Flush()
	if err != nil {
		log.Fatalf("Failed to write to file %q: %v", file, err)
	}
}

func writeKeepAlivedServiceFiles(config apis.VIPConfiguration) {
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
		config.RouterID = constants.DEFAULT_ROUTER_ID
	}
	kaConfFileTemplate :=
		`vrrp_instance K8S_APISERVER {
	interface {{.NetworkInterface}}
	state BACKUP
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
	confFile := filepath.Join(constants.SYSTEMD_DIR, "keepalived.conf")
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
	kaServiceData := KaServiceData{confFile, constants.KEEPALIVED_IMG}
	writeTemplateIntoFile(kaSvcFileTemplate, "kaSvcFileTemplate", filepath.Join(constants.SYSTEMD_DIR, "keepalived.service"), kaServiceData)
}
