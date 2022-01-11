/*
Copyright © 2022 john wozniak jwozniak.dev@gmail.com>

*/
package cligen

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	polyConfig string
	apiKey     string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "upload latest config",
	Long: `upload the poc latest config so the cli can read in the config and proceed with 
			reading in the roundStartTime and roundEndTime.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO:
		// 1. start with reading in string args and printing it back out
		// 2. allow for user to paste in path for the config and parse the json and read the file
		// 3. once read fetch the data from polyscan
		if err := readArgs(args); err != nil {
			return err
		}
		return nil
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

func readArgs(args []string) error {
	var cla CommandLineArgs

	cla.RoundStartTime = args[0]
	cla.RoundEndTime = args[1]

	res, err := fetchBlockData(cla)
	if err != nil {
		return err
	}

	r, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(r))
	return nil
}

func fetchBlockData(cla CommandLineArgs) (CommandLineResponse, error) {
	var clr CommandLineResponse
	var err error

	// TODO: move to viper
	url := fmt.Sprintf("https://api.polygonscan.com/api?module=block&action=getblocknobytime&timestamp=%v&closest=before&apikey=%v", cla.RoundStartTime, apiKey)
	url2 := fmt.Sprintf("https://api.polygonscan.com/api?module=block&action=getblocknobytime&timestamp=%v&closest=before&apikey=%v", cla.RoundEndTime, apiKey)
	// fetch data
	sb, err := get(url)
	if err != nil {
		return clr, err
	}
	eb, err := get(url2)

	clr.StartBlock = sb
	clr.EndBlock = eb

	return clr, nil
}

func get(url string) (string, error) {
	// TODO: allow for making concurrent api calls - allowing to make both calls simultaneously
	var res PolyscanResponse
	if apiKey == "" {
		err := fmt.Errorf("error reading in apikey")
		return "", err
	}

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

	return res.Result, nil

}

func init() {
	addCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	addCmd.Flags().StringVar(&polyConfig, "poly-config", "", "pass the poc-polygon-latest.json")
	addCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "provide the api key for polyscan")
	err := viper.BindPFlag("apikey", addCmd.PersistentFlags().Lookup("apikey"))
	if err != nil {
		fmt.Println("ERROR WITH BIND")
		return
	}
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
