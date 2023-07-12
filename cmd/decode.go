package cmd

import (
	"fmt"
	"os"

	"github.com/leonardobiffi/kse/package/decode"
	"github.com/leonardobiffi/kse/package/utils"
	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode Kubernetes Secrets",
	Long:  `Decode Kubernetes Secrets.`,
	RunE:  runDecode,
}

var flagsDecode struct {
	file      string
	overwrite bool
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	decodeCmd.Flags().StringVarP(&flagsDecode.file, "file", "f", "", "file path to decode")
	decodeCmd.Flags().BoolVarP(&flagsDecode.overwrite, "overwrite", "o", false, "overwrite file")
}

func runDecode(cmd *cobra.Command, args []string) (err error) {
	var content []byte

	svc := decode.New()

	if flagsDecode.file != "" {
		files, err := utils.FindFiles(flagsDecode.file)
		if err != nil {
			return err
		}

		for _, file := range files {
			// only print filename without overwrite flag
			if flagsDecode.overwrite {
				fmt.Println("Decoding file:", file)
			}

			content, err = utils.ReadFile(file)
			if err != nil {
				return err
			}

			out, err := svc.Execute(content)
			if err != nil {
				return err
			}

			if flagsDecode.overwrite {
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
