/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"goldfish/internal/tapd"
	"log"
	"os"
)

var force *bool

// tapdCmd represents the tapd command
var tapdCmd = &cobra.Command{
	Use:   "tapd",
	Short: "Sync tapd stories and bugs into feishu bitable",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 设置日志输出到标准输出
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		tapd.New().Run(*force)
	},
}

func init() {
	force = tapdCmd.PersistentFlags().BoolP("force", "f", false, "clean table before sync")
	rootCmd.AddCommand(tapdCmd)
}
