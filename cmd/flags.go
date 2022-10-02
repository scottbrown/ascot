package cmd

var Profile string
var ShowRequiredPermissions bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "The AWS profile to use")
	rootCmd.PersistentFlags().BoolVarP(&ShowRequiredPermissions, "show-required-permissions", "", false, "Shows the IAM permissions required by the command and then exits")
}
