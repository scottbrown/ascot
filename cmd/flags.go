package cmd

var Profile string
var ShowRequiredPermissions bool
var HowItWorks bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "The AWS profile to use")
	rootCmd.PersistentFlags().BoolVarP(&ShowRequiredPermissions, "show-required-permissions", "", false, "Shows the IAM permissions required by the command and then exits")
	rootCmd.PersistentFlags().BoolVarP(&HowItWorks, "how-it-works", "", false, "Shows the business logic for the command")
}
