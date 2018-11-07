package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/platform9/nodeadm/pkg/kubeadm"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	log "github.com/platform9/nodeadm/pkg/logrus"
	"github.com/platform9/nodeadm/utils"
	executil "github.com/platform9/nodeadm/utils/exec"
	"github.com/spf13/cobra"
)

var (
	initCfgPath string
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := &apis.InitConfiguration{}
		if initCfgPath != "" {
			config, err = utils.InitConfigurationFromFile(initCfgPath)
			if err != nil {
				log.Fatalf("Failed to read configuration from file %q: %v", initCfgPath, err)
			}
		}
		apis.SetInitDefaults(config)
		if err := apis.SetInitDynamicDefaults(config); err != nil {
			log.Fatalf("Failed to set dynamic defaults: %v", err)
		}
		if errors := apis.ValidateInit(config); len(errors) > 0 {
			log.Error("Failed to validate configuration:")
			for i, err := range errors {
				log.Errorf("%v: %v", i, err)
			}
			os.Exit(1)
		}

		if err := kubeadm.WriteConfiguration(constants.KubeadmConfig, config.MasterConfiguration); err != nil {
			log.Fatalf("Unable to write kubeadm configuration to %s: %s", constants.KubeadmConfig, err)
		}

		utils.InstallMasterComponents(config)

		kubeadmInit(constants.KubeadmConfig)

		log.Infoln("Applying workaround for https://github.com/kubernetes/kubeadm/issues/857")
		if err := ensureKubeProxyRespectsHostoverride(); err != nil {
			log.Fatalf("Failed to apply workaround: %v", err)
		}

		log.Println("Configuring pod network")
		if err := networkInit(config); err != nil {
			log.Fatalf("Unable to configure pod network: %s", err)
		}
	},
}

func networkInit(config *apis.InitConfiguration) error {
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse masterConfiguration: %s", err)
	}
	if !p.ExistsP("networking.podSubnet") {
		return fmt.Errorf("masterConfiguration.networking.podSubnet must be defined")
	}
	podSubnet, ok := p.Path("networking.podSubnet").Data().(string)
	if !ok {
		return fmt.Errorf("masterConfiguration.networking.podSubnet must be a string")
	}

	log.Println("Setting net.bridge.bridge-nf-call-iptables=1")
	cmd := exec.Command(constants.Sysctl, "net.bridge.bridge-nf-call-iptables=1")
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}

	log.Println("Applying flannel daemonset from disk")
	manifest, err := ioutil.ReadFile(filepath.Join(constants.CacheDir, constants.FlannelDirName, constants.FlannelManifestFilename))
	if err != nil {
		log.Fatalf("failed to open network backend manifest (%q): %s", strings.Join(cmd.Args, " "), err)
	}

	log.Printf("Setting flannel net-conf.json Network to %s", podSubnet)
	manifestWithPodSubnet := strings.Replace(string(manifest), constants.DefaultPodNetwork, podSubnet, -1)
	cmd = exec.Command(filepath.Join(constants.BaseInstallDir, "kubectl"), fmt.Sprintf("--kubeconfig=%s", constants.AdminKubeconfigFile), "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifestWithPodSubnet)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func kubeadmInit(config string) {
	cmd := exec.Command(filepath.Join(constants.BaseInstallDir, "kubeadm"), "init", "--ignore-preflight-errors=all", "--config="+config)
	err := executil.LogRun(cmd)
	if err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().StringVar(&initCfgPath, "cfg", "", "Location of configuration file")
}
