/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"github.com/BASChain/go-bas-mail-server/app/cmdcommon"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bas-mail-server/rsakey"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var keypassword string
var remoteserver string

func inputpassword() (password string, err error) {
	passwd, err := gopass.GetPasswdPrompt("Please Enter Password:", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	if len(passwd) < 1 {
		return "", errors.New("Please input valid password")
	}

	return string(passwd), nil
}

func inputChoose() (choose string, err error) {
	c, err := gopass.GetPasswdPrompt("Do you reinit config[yes/no]: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	return string(c), nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init bms",
	Long:  `init bms`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		_, err = cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()

		if rsakey.KeyGenerated(config.GetBMSCfg().GetKeyPath()) {
			var choose string
			if choose, err = inputChoose(); err != nil {
				log.Println(err)
			}

			if choose != "yes" {
				log.Println("init break, use old configuration")
				return
			}
		}

		if keypassword == "" {
			if keypassword, err = inputpassword(); err != nil {
				log.Println(err)
				return
			}
		}

		if keypassword == "" {
			log.Println("Please input password")
			return
		}

		if remoteserver == "" {
			log.Println("Please input remote host ip address")
			return
		}

		if err = genRsaKey(keypassword); err != nil {
			panic("Generate rsa key pair failed")
		}

		cfg := config.GetBMSCfg()
		cfg.RemoteServer = remoteserver

		cfg.Save()

		s58 := rsakey.PubKey2Addr(config.GetBMSCfg().PubKey)

		log.Println("Init success!, Public key: ", s58)
	},
}

func genRsaKey(password string) error {
	priv, _ := rsakey.GenerateKeyPair(2048)

	cfg := config.GetBMSCfg()

	err := rsakey.Save2File(cfg.GetKeyPath(), priv, password)
	if err != nil {
		return err
	}

	cfg.SetPrivKey(priv)
	cfg.SetPubKey(&priv.PublicKey)

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.Flags().StringVarP(&keypassword, "password", "p", "", "password for key encrypt")
	initCmd.Flags().StringVarP(&remoteserver, "host", "r", "", "remote server ip address")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
