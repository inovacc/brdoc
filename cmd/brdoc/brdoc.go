/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"os"

	sdk "github.com/inovacc/brdoc"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "brdoc",
	Short: "Brazilian documents utilities (CPF/CNPJ)",
	Long:  "brdoc is a small CLI to generate and validate Brazilian documents like CPF and CNPJ.",
}

func init() {
	rootCmd.AddCommand(cpfCmd)
	rootCmd.AddCommand(cnpjCmd)
}

// -----------------------------
// CPF command
// -----------------------------
var (
	cpfGenerate bool
	cpfValidate string
)

var cpfCmd = &cobra.Command{
	Use:     "cpf",
	Short:   "Generate or validate CPF",
	Example: "brdoc cpf --generate\nbrdoc cpf --validate 123.456.789-09",
	RunE: func(cmd *cobra.Command, args []string) error {
		if (cpfGenerate && cpfValidate != "") || (!cpfGenerate && cpfValidate == "") {
			// If both provided or none provided, show help with the error
			_ = cmd.Help()
			if cpfGenerate && cpfValidate != "" {
				return errors.New("flags --generate and --validate are mutually exclusive")
			}

			return errors.New("either --generate or --validate must be provided")
		}

		c := sdk.NewCPF()
		if cpfGenerate {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), c.Generate())
			return nil
		}

		// validate path
		valid := c.Validate(cpfValidate)
		if valid {
			// also print formatted variant
			if formatted, err := c.Format(cpfValidate); err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "valid\t%s\n", formatted)
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "valid")
			}

			return nil
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "invalid")
		cmd.SilenceUsage = true

		return nil
	},
}

// Flags for cpf
func init() {
	cpfCmd.Flags().BoolVarP(&cpfGenerate, "generate", "g", false, "Generate a valid CPF")
	cpfCmd.Flags().StringVarP(&cpfValidate, "validate", "v", "", "Validate a CPF value")
}

// -----------------------------
// CNPJ command
// -----------------------------
var (
	cnpjGenerate bool
	cnpjValidate string
)

var cnpjCmd = &cobra.Command{
	Use:     "cnpj",
	Short:   "Generate or validate CNPJ",
	Example: "brdoc cnpj --generate\nbrdoc cnpj --validate 12.345.678/0001-95",
	RunE: func(cmd *cobra.Command, args []string) error {
		if (cnpjGenerate && cnpjValidate != "") || (!cnpjGenerate && cnpjValidate == "") {
			_ = cmd.Help()
			if cnpjGenerate && cnpjValidate != "" {
				return errors.New("flags --generate and --validate are mutually exclusive")
			}

			return errors.New("either --generate or --validate must be provided")
		}

		c := sdk.NewCNPJ()
		if cnpjGenerate {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), c.Generate())
			return nil
		}

		valid := c.Validate(cnpjValidate)
		if valid {
			if formatted, err := c.Format(cnpjValidate); err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "valid\t%s\n", formatted)
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "valid")
			}

			return nil
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "invalid")
		cmd.SilenceUsage = true

		return nil
	},
}

// Flags for cnpj
func init() {
	cnpjCmd.Flags().BoolVarP(&cnpjGenerate, "generate", "g", false, "Generate a valid CNPJ")
	cnpjCmd.Flags().StringVarP(&cnpjValidate, "validate", "v", "", "Validate a CNPJ value")
}
