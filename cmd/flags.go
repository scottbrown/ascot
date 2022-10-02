package cmd

var Profile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "The AWS profile to use")
}
