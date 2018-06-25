package cmd

import (
	"path/filepath"

	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

var overwriteSymlink bool

// nodeCmd represents the cluster command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download components",
	Run: func(cmd *cobra.Command, args []string) {
		var kubeRootDir = filepath.Join(utils.BASE_DIR)
		var cniRootDir = filepath.Join(utils.CNI_BASE_DIR)

		var kubeDir = filepath.Join(kubeRootDir, "kubernetes-"+utils.KUBERNETES_VERSION)
		var cniDir = filepath.Join(cniRootDir, "cni-"+utils.CNI_VERSION)
		utils.DownloadKubeComponents(kubeDir, utils.KUBERNETES_VERSION)
		utils.CreateSymLinks(kubeDir, kubeRootDir, overwriteSymlink)

		utils.DownloadCNIPlugin(cniDir, utils.CNI_VERSION)
		utils.CreateSymLinks(cniDir, cniRootDir, overwriteSymlink)

		utils.DownloadKubeletServiceFiles(kubeRootDir, utils.KUBERNETES_VERSION)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolVar(&overwriteSymlink, "overwriteSymlink", false, "Overwrite the symlinks")
}
