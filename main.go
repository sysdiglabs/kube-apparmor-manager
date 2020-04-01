package main

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	"github.com/sysdiglabs/kube-apparmor-manager/aa"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	appArmor *aa.AppArmor
	err      error
)

func init() {
	appArmor, err = aa.NewAppArmor()

	if err != nil {
		panic(err)
	}
}

func main() {
	var logLevel string

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	var rootCmd = &cobra.Command{
		Use:   "kube-apparmor-manager",
		Short: "kube-apparmor-manager manages AppArmor service and profiles enforcement on worker nodes",
		Long:  "kube-apparmor-manager manages AppArmor service and profiles enforcement on worker nodes through syncing with AppArmor CRD in Kubernetes cluster",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			lvl, err := log.ParseLevel(logLevel)
			if err != nil {
				log.Fatal(err)
			}

			log.SetLevel(lvl)
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "level", "info", "Log level")

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Install CRD in the cluster and AppArmor services on worker nodes",
		Long:  "Install CRD in the Kubernetes cluster database and AppArmor services on worker nodes",
		Run: func(cmd *cobra.Command, args []string) {
			err := appArmor.InstallCRD()
			if err != nil {
				log.Fatalf("failed to install CRD: %v", err)
			}

			err = appArmor.InstallAppArmor()
			if err != nil {
				log.Fatalf("failed to install AppArmor service: %v", err)
			}
		},
	}

	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize the AppArmor profiles from the Kubernetes database (etcd) to worker nodes",
		Long:  "Synchronize the AppArmor profiles from the Kubernetes database (etcd) to worker nodes",
		Run: func(cmd *cobra.Command, args []string) {
			err := appArmor.Sync()
			if err != nil {
				log.Fatalf("sync error: %v", err)
			}
		},
	}

	var enforcedCmd = &cobra.Command{
		Use:   "enforced",
		Short: "Check AppArmor profile enforcement status on worker nodes",
		Long:  "Check AppArmor profile enforcement status on worker nodes",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := appArmor.AppArmorStatus()
			if err != nil {
				log.Fatalf("check enforcement status error: %v", err)
			}

			list.PrintEnforcementStatus()
		},
	}

	var enabledCmd = &cobra.Command{
		Use:   "enabled",
		Short: "Check AppArmor status on worker nodes",
		Long:  "Check AppArmor status on worker nodes",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := appArmor.AppArmorEnabled()
			if err != nil {
				log.Fatalf("check enabled status error: %v", err)
			}

			list.PrintEnabledStatus()
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(enforcedCmd)
	rootCmd.AddCommand(enabledCmd)

	rootCmd.Execute()
}
