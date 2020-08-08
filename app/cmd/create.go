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
	"github.com/realbmail/go-bas-mail-server/app/cmdclient"
	"github.com/realbmail/go-bas-mail-server/app/cmdcommon"
	"github.com/realbmail/go-bas-mail-server/bmailcrypt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var keypassword string

func inputpassword() (password string, err error) {
	passwd, err := gopass.GetPasswdPrompt("Please Enter Password: ", true, os.Stdin, os.Stdout)
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

//var accountofflinemode *bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create account",
	Long:  `create account`,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}

		var err error

		if bmailcrypt.KeyIsGenerated() {
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

		cmdclient.StringOpCmdSend("", cmdcommon.CMD_ACCOUNT_CREATE, keypassword)
	},
}

func init() {
	accountCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//accountofflinemode = createCmd.Flags().BoolP("offline", "o", false, "create account offline offline")

	//createCmd.Flags().BoolVarP(&accountofflinemode, "offline", "o", false, "create account offline offline")

}
