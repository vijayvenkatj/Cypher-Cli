package cmd
import (
	"fmt"
	"github.com/spf13/cobra"
)


var Show = &cobra.Command{
	Use: "show",
	Short: "Show all passwords in the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		var vault InputJSON
		vault,err := GetVault()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		for _, pass := range vault.Passwords {
			fmt.Println("ID: ", pass.Id)
			fmt.Println("Name: ", pass.Name)
			fmt.Println("Username: ", pass.Username)
		}
	},
}
