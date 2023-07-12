package cmd

import (
	"fmt"
	"os"

	"github.com/leonardobiffi/kse/package/encode"
	"github.com/leonardobiffi/kse/package/utils"
	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode Kubernetes Secrets",
	Long:  `Encode Kubernetes Secrets.`,
	RunE:  runEncode,
}

var flagsEncode struct {
	file      string
	overwrite bool
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().StringVarP(&flagsEncode.file, "file", "f", "", "file path to encode")
	encodeCmd.Flags().BoolVarP(&flagsEncode.overwrite, "overwrite", "o", false, "overwrite file")
}

func runEncode(cmd *cobra.Command, args []string) (err error) {
	var content []byte

	svc := encode.New()

	if flagsEncode.file != "" {
		files, err := utils.FindFiles(flagsEncode.file)
		if err != nil {
			return err
		}

		for _, file := range files {
			// only print filename without overwrite flag
			if flagsEncode.overwrite {
				fmt.Println("Encoding file:", file)
			}

			content, err = utils.ReadFile(file)
			if err != nil {
				return err
			}

			out, err := svc.Execute(content)
			if err != nil {
				return err
			}

			if flagsEncode.overwrite {
				err = utils.UpdateFile(file, out)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(string(out))
			}
		}

		return nil
	}

	if os.Stdin != nil {
		content = utils.ReadStdin(os.Stdin)
	}

	out, err := svc.Execute(content)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return
}
