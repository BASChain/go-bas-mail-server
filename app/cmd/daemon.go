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
	//"github.com/BASChain/go-bas-discover-wanip/dns/server"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"

	"github.com/BASChain/go-bas-mail-server/app/cmdcommon"
	"github.com/BASChain/go-bas-mail-server/app/cmdservice"
	"github.com/BASChain/go-bas-mail-server/config"
	//"github.com/BASChain/go-bas-discover-wanip/dns/dohserver"
	//"github.com/BASChain/go-bas-discover-wanip/dns/mem"
	"github.com/BASChain/go-bas-mail-server/rsakey"
	"path"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "ssc start in backend",
	Long:  `ssc start in backend`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()
		cfg := config.GetBMCCfg()
		cfg.Save()
		//if keypassword == ""{
		//	if keypassword, err = inputpassword(); err != nil {
		//		log.Println(err)
		//		return
		//	}
		//}

		if keypassword == "" {
			log.Println("Please input password")
			return
		}

		if priv, pub, err := rsakey.LoadRSAKey(cfg.GetKeyPath(), []byte(keypassword)); err != nil {
			log.Println("Recover RSA Key Failed")
			return
		} else {
			cfg.SetPrivKey(priv)
			cfg.SetPubKey(pub)
			cfg.PKAddr = rsakey.PubKey2Addr(pub)
		}

		daemondir := config.GetBMCHomeDir()
		cntxt := daemon.Context{
			PidFileName: path.Join(daemondir, "ssc.pid"),
			PidFilePerm: 0644,
			LogFileName: path.Join(daemondir, "ssc.log"),
			LogFilePerm: 0640,
			WorkDir:     daemondir,
			Umask:       027,
			Args:        []string{},
		}
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			log.Println("ssc starting, please check log at:", path.Join(daemondir, "ssc.log"))
			return
		}
		defer cntxt.Release()

		//BAS_Ethereum.RecoverContract()
		//go service.StartService()
		//go mem.MemStateStart()
		//go server.DNSServerDaemon()
		//go dohserver.GetDohDaemonServer().StartDaemon()
		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//daemonCmd.Flags().IntVarP(&cmdroottcpport, "tcp-listen-port", "t", 65566, "local tcp listen port")
	//daemonCmd.Flags().IntVarP(&cmdrootudpport, "udp-listen-port", "u", 65566, "local udp listen port")
	//daemonCmd.Flags().StringVarP(&cmdropstennap, "ropsten-network-access-point", "r", "", "ropsten network access point")
	//daemonCmd.Flags().StringVarP(&cmdbastokenaddr, "bas-token-address", "a", "", "bas token address")
	//daemonCmd.Flags().StringVarP(&cmdbasmgraddr, "bas-mgr-address", "m", "", "bas manager address")
	daemonCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")
	daemonCmd.Flags().StringVarP(&keypassword, "password", "p", "", "password for key encrypt")
	//daemonCmd.Flags().IntVarP(&cmddohserverport, "doh-listen-port", "p", 65566, "local doh server listen port")
	//daemonCmd.Flags().StringVarP(&cmdcertfile, "cert-file", "f", "", "certificate file for tls")
	//daemonCmd.Flags().StringVarP(&cmdkeyfile, "key-file", "k", "", "private key file for tls")
	//daemonCmd.Flags().StringVarP(&cmddnspath, "dns-query-path", "q", "", "path for dns query")
	//daemonCmd.Flags().IntVarP(&cmdquerydnstimeout, "dns-query-time", "o", 0, "max time for wait remote dns server reply")
	//daemonCmd.Flags().IntVarP(&cmdquerydnstrytimes, "dns-query-times", "s", 0, "max times for sending dns to remote dns server ")

}
