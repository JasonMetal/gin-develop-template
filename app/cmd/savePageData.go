package cmd

import (
	"develop-template/app/constant"
	"develop-template/app/logic/crawlerLogic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
	time2 "time"
)

var savePageDataCmd = &cobra.Command{
	Use:   "savePageData",
	Short: "爬数据入库",
	Long:  "把test数据抓取入库",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("savePageData called")

		savePageData()
	},
}

func init() {
	rootCmd.AddCommand(savePageDataCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncTemplateTypeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncTemplateTypeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func savePageData() {
	GCtx := gin.Context{
		Request: &http.Request{},
	}
	logic := crawlerLogic.NewLogic(&GCtx)
	getAllTagsData := logic.GetAllTagsData()
	fmt.Printf("\nRun getAllTagsData..., %s\n%s\n", getAllTagsData, time2.Now().Format(constant.DefaultDateFormat))
}
