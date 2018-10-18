package cmd

import (
	log "github.com/platform9/nodeadm/logs"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/deprecated"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := &apis.JoinConfiguration{}
		configPath := cmd.Flag("cfg").Value.String()
		if len(configPath) != 0 {
			config, err = utils.JoinConfigurationFromFile(configPath)
			if err != nil {
				log.Fatalf("Failed to read configuration from file %q: %v", configPath, err)
			}
		}
		apis.SetJoinDefaults(config)
		utils.InstallNodeComponents(config)
		kubeadmJoin(cmd.Flag("token").Value.String(),
			cmd.Flag("master").Value.String(),
			cmd.Flag("cahash").Value.String())
	},
}

func kubeadmJoin(token, master, cahash string) {
	deprecated.Run(constants.BaseInstallDir, "kubeadm", "join", "--ignore-preflight-errors=all", "--token", token, master, "--discovery-token-ca-cert-hash", cahash)
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().String("cfg", "", "Location of configuration file")
	nodeCmdJoin.Flags().String("token", "", "kubeadm token to be used for kubeadm join")
	nodeCmdJoin.Flags().String("master", "", "masterIP:masterPort for the master to join")
	nodeCmdJoin.Flags().String("cahash", "", "CA hash")
}
