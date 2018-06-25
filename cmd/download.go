package cmd

import (
	"path/filepath"

	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
	"log"
)

var kube, cni, overwriteSymlink bool
var downloadDir string

// nodeCmd represents the cluster command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download components",
	Run: func(cmd *cobra.Command, args []string) {
		var kubeRootDir = filepath.Join(downloadDir, "kubernetes")
		var cniRootDir = filepath.Join(downloadDir, "cni")

		var kubeDir = filepath.Join(kubeRootDir, "kubernetes-"+utils.KUBERNETES_VERSION)
		var cniDir = filepath.Join(cniRootDir, "cni-"+utils.CNI_VERSION)

		if !kube && !cni {
			log.Print("Please set which component to download")
		}

		if kube {
			utils.DownloadKubeComponents(kubeDir, utils.KUBERNETES_VERSION)
			utils.CreateSymLinks(kubeDir, kubeRootDir, overwriteSymlink)
		}

		if cni {
			utils.DownloadCNIPlugin(cniDir, utils.CNI_VERSION)
			utils.CreateSymLinks(cniDir, cniRootDir, overwriteSymlink)
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolVar(&kube, "kube", true, "Download Kubernetes components")
	downloadCmd.Flags().BoolVar(&cni, "cni", true, "Download CNI plugin")
	downloadCmd.Flags().BoolVar(&overwriteSymlink, "overwriteSymlink", false, "Overwrite the symlinks")
	downloadCmd.Flags().StringVar(&downloadDir, "downloadDir", "/opt", "Destination directory")
}
