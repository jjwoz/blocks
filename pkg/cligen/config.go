/*
Copyright © 2022 john wozniak jwozniak.dev@gmail.com>

*/
package cligen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "upload latest config",
	Long: `upload the poc latest config so the cli can read in the config and proceed with 
			reading in the roundStartTime and roundEndTime.`,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:
		// 1. start with reading in string args and printing it back out
		// 2. allow for user to paste in path for the config and parse the json and read the file
		// 3. once read fetch the data from polyscan
		fmt.Println("config called start reading.....")
		fmt.Println("\n")
		readArgs(args)
	},
}

// CommandLineArgs the input from the cli
type CommandLineArgs struct {
	RoundStartTime string
	RoundEndTime   string
}

// CommandLineResponse the response from polygon scan
type CommandLineResponse struct {
	StartBlock string
	EndBlock   string
}

type PolyscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func readArgs(args []string) {
	var cla CommandLineArgs

	cla.RoundStartTime = args[0]
	cla.RoundEndTime = args[1]

	res := fetchBlockData(cla)
	r, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(r))

}

func fetchBlockData(cla CommandLineArgs) CommandLineResponse {
	var clr CommandLineResponse

	apiKey := "ADD"
	url := fmt.Sprintf("https://api.polygonscan.com/api?module=block&action=getblocknobytime&timestamp=%v&closest=before&apikey=%v", cla.RoundStartTime, apiKey)
	url2 := fmt.Sprintf("https://api.polygonscan.com/api?module=block&action=getblocknobytime&timestamp=%v&closest=before&apikey=%v", cla.RoundEndTime, apiKey)
	// fetch data
	sb := get(url)
	eb := get(url2)

	clr.StartBlock = sb
	clr.EndBlock = eb

	return clr
}

func get(url string) string {
	// TODO: allow for making concurrent api calls - allowing to make both calls simultaneously
	var res PolyscanResponse

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
	}

	return res.Result

}

func init() {
	addCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
