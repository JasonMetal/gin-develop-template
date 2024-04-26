package cmd

import (
	"context"
	"develop-template/app/constant"
	"develop-template/app/logic/crawlerLogic"
	"fmt"
	time2 "time"

	"github.com/spf13/cobra"
)

var savePageDataCronCmd = &cobra.Command{
	Use:   "savePageDataCron",
	Short: "爬数据入库",
	Long:  "把test数据抓取入库",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("savePageDataCron called")

		savePageDataCron()
	},
}

func init() {
	rootCmd.AddCommand(savePageDataCronCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncTemplateTypeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncTemplateTypeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func savePageDataCron() {
	//GCtx := gin.Context{
	//	Request: &http.Request{},
	//}
	ctx := context.Background()
	logic := crawlerLogic.NewLogic(ctx)
	logic.CreateVodDetailDataCron()

	fmt.Printf("\nRun savePageDataCron..., %s\n%s\n", time2.Now().Format(constant.DefaultDateFormat))
}
