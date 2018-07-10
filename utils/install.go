package utils

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	EXECUTE = 0744
	READ    = 0644
)

func InstallMasterComponents(config *Configuration) {
	PopulateCache()
	PlaceComponentsFromCache()
	ReplaceString(getKubeletServiceConf(), DEFAULT_DNS_IP, GetIPFromSubnet(config.MasterConfiguration.Networking.ServiceSubnet, 10))
	EnableAndStartService("kubelet.service")
	writeKeepAlivedServiceFiles(config.VipConfiguration)
	EnableAndStartService("keepalived.service")
}

func InstallWorkerComponents() {
	PopulateCache()
	PlaceComponentsFromCache()
	EnableAndStartService("kubelet.service")
}

func writeTemplateIntoFile(tmpl, name, file string, data interface{}) {
	err := os.MkdirAll(filepath.Dir(file), READ)
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
