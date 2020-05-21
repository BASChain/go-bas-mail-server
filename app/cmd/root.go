// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/BASChain/go-bas-mail-server/app/cmdcommon"
	"github.com/BASChain/go-bas-mail-server/config"

	"github.com/BASChain/go-bas-mail-server/app/cmdservice"
	"github.com/spf13/cobra"
	"log"
)

////var cfgFile string
//
var (
	cmdconfigfilename string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bms",
	Short: "start bms in current shell",
	Long:  `start bms in current shell`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()
		cfg := config.GetBMSCfg()
		cfg.Save()


		//go bmtpserver.GetBMTPServer().StartTCPServer()

		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitCfg() {
	if cmdconfigfilename != "" {
		cfg := config.LoadFromCfgFile(cmdconfigfilename)
		if cfg == nil {
			return
		}
	} else {
		config.LoadFromCmd(cfginit)
	}
	//Set2SmartContract()
}

//
//func Set2SmartContract() {
//	cfg := config.GetBasDisCfg()
//
//	//fmt.Println(*cfg)
//
//	//if cfg.RopstenNAP != "" {
//	//	BAS_Ethereum.RopstenNetworkAccessPoint = cfg.RopstenNAP
//	//}
//	//
//	//if cfg.TokenAddr != "" {
//	//	BAS_Ethereum.BASTokenAddress = cfg.TokenAddr
//	//}
//	//
//	//if cfg.MgrAddr != "" {
//	//	BAS_Ethereum.BASManagerSimpleAddress = cfg.MgrAddr
//	//}
//
//}

func cfginit(bc *config.BMSConfig) *config.BMSConfig {
	cfg := bc
	//if cmdrootudpport > 0 && cmdrootudpport < 65535 {
	//	cfg.UpdPort = cmdrootudpport
	//}
	//if cmdroottcpport > 0 && cmdroottcpport < 65535 {
	//	cfg.TcpPort = cmdroottcpport
	//}
	//if cmdropstennap != "" {
	//	cfg.RopstenNAP = cmdropstennap
	//}
	//if cmdbastokenaddr != "" {
	//	cfg.TokenAddr = cmdbastokenaddr
	//}
	//if cmdbasmgraddr != "" {
	//	cfg.MgrAddr = cmdbasmgraddr
	//}
	//if cmddohserverport > 0 && cmddohserverport < 65535 {
	//	cfg.DohServerPort = cmddohserverport
	//}
	//
	//if cmdcertfile != "" {
	//	cfg.CertFile = cmdcertfile
	//}
	//if cmdkeyfile != "" {
	//	cfg.KeyFile = cmdkeyfile
	//}
	//if cmdquerydnstimeout != 0 {
	//	cfg.TimeOut = cmdquerydnstimeout
	//}
	//if cmdquerydnstrytimes != 0 {
	//	cfg.TryTimes = cmdquerydnstrytimes
	//}

	return cfg

}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//rootCmd.Flags().IntVarP(&cmdroottcpport, "tcp-listen-port", "t", 65566, "local tcp listen port")
	//rootCmd.Flags().IntVarP(&cmdrootudpport, "udp-listen-port", "u", 65566, "local udp listen port")
	//rootCmd.Flags().StringVarP(&cmdropstennap, "ropsten-network-access-point", "r", "", "ropsten network access point")
	//rootCmd.Flags().StringVarP(&cmdbastokenaddr, "bas-token-address", "a", "", "bas token address")
	//rootCmd.Flags().StringVarP(&cmdbasmgraddr, "bas-mgr-address", "m", "", "bas manager address")
	rootCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")
	rootCmd.Flags().StringVarP(&keypassword, "password", "p", "", "password for key encrypt")
	//rootCmd.Flags().IntVarP(&cmddohserverport, "doh-listen-port", "p", 65566, "local doh server listen port")
	//rootCmd.Flags().StringVarP(&cmdcertfile, "cert-file", "f", "", "certificate file for tls")
	//rootCmd.Flags().StringVarP(&cmdkeyfile, "key-file", "k", "", "private key file for tls")
	//rootCmd.Flags().StringVarP(&cmddnspath, "dns-query-path", "q", "", "path for dns query")
	//rootCmd.Flags().IntVarP(&cmdquerydnstimeout, "dns-query-time", "o", 0, "max time for wait remote dns server reply")
	//rootCmd.Flags().IntVarP(&cmdquerydnstrytimes, "dns-query-times", "s", 0, "max times for sending dns to remote dns server ")
}

//
//// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".app" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".app")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
