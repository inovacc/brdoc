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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/inovacc/brdoc"
	"github.com/spf13/cobra"
)

const maxLine = 1024 * 1024

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

var (
	buf          = make([]byte, 0, 64*1024)
	cpfGenerate  bool
	cpfValidate  string
	cpfFrom      string
	cpfCount     int
	cnpjGenerate bool
	cnpjValidate string
	cnpjFrom     string
	cnpjCount    int
	cnpjLegacy   bool
)

var rootCmd = &cobra.Command{
	Use:   "brdoc",
	Short: "Brazilian documents utilities (CPF/CNPJ)",
	Long:  "brdoc is a small CLI to generate and validate Brazilian documents like CPF and CNPJ.",
}

// Flags for cnpj
func init() {
	cnpjCmd.Flags().BoolVarP(&cnpjGenerate, "generate", "g", false, "Generate a valid CNPJ")
	cnpjCmd.Flags().StringVarP(&cnpjValidate, "validate", "v", "", "Validate a CNPJ value")
	cnpjCmd.Flags().StringVarP(&cnpjFrom, "from", "f", "", "Validate many CNPJs from file or '-' for stdin")
	cnpjCmd.Flags().IntVarP(&cnpjCount, "count", "n", 0, "When generating, how many CNPJs to output")
	cnpjCmd.Flags().BoolVar(&cnpjLegacy, "legacy", false, "When generating, output legacy numeric-only CNPJ (12 digits base + 2 numeric check digits)")

	cpfCmd.Flags().BoolVarP(&cpfGenerate, "generate", "g", false, "Generate a valid CPF")
	cpfCmd.Flags().StringVarP(&cpfValidate, "validate", "v", "", "Validate a CPF value")
	cpfCmd.Flags().StringVarP(&cpfFrom, "from", "f", "", "Validate many CPFs from file or '-' for stdin")
	cpfCmd.Flags().IntVarP(&cpfCount, "count", "n", 0, "When generating, how many CPFs to output")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Avoid duplicate help/usage or error printing when returning errors from RunE
	// We handle error printing in main().
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	rootCmd.AddCommand(cpfCmd)
	rootCmd.AddCommand(cnpjCmd)
}

var cpfCmd = &cobra.Command{
	Use:   "cpf",
	Short: "Generate or validate CPF",
	Example: strings.Join([]string{
		"brdoc cpf --generate",
		"brdoc cpf --generate --count 10",
		"brdoc cpf --validate 123.456.789-09",
		"brdoc cpf --validate --from cpfs.txt",
		"type cpfs.txt | brdoc cpf --validate --from -",
	}, "\n"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags combination
		if cpfGenerate && (cpfValidate != "" || cpfFrom != "") {
			return errors.New("--generate cannot be used with --validate or --from")
		}

		if cpfFrom != "" && cpfValidate != "" {
			return errors.New("--from and --validate are mutually exclusive for CPF")
		}

		if !cpfGenerate && cpfValidate == "" && cpfFrom == "" {
			return errors.New("either --generate, --validate, or --from must be provided")
		}

		c := sdk.NewCPF()
		if cpfGenerate {
			if cpfCount <= 0 {
				cpfCount = 1
			}

			w := bufio.NewWriter(cmd.OutOrStdout())
			defer func(w *bufio.Writer) {
				if err := w.Flush(); err != nil {
					panic(err)
				}
			}(w)

			for i := 0; i < cpfCount; i++ {
				_, _ = fmt.Fprintln(w, c.Generate())
			}

			return nil
		}

		// validate single or bulk
		if cpfFrom != "" { // bulk from file or stdin
			r, closeFn, err := openReader(cpfFrom)
			if err != nil {
				return err
			}

			if closeFn != nil {
				defer closeFn()
			}

			scanner := bufio.NewScanner(r)
			// Increase buf in case of long lines
			scanner.Buffer(buf, maxLine)

			w := bufio.NewWriter(cmd.OutOrStdout())
			defer func(w *bufio.Writer) {
				if err := w.Flush(); err != nil {
					panic(err)
				}
			}(w)

			anyInvalid := false
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				if c.Validate(line) {
					if formatted, err := c.Format(line); err == nil {
						_, _ = fmt.Fprintf(w, "valid\t%s\n", formatted)
					} else {
						_, _ = fmt.Fprintln(w, "valid")
					}
				} else {
					anyInvalid = true
					_, _ = fmt.Fprintf(w, "invalid\t%s\n", line)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			if anyInvalid {
				cmd.SilenceUsage = true
			}

			return nil
		}

		// single validate value
		valid := c.Validate(cpfValidate)
		if valid {
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

var cnpjCmd = &cobra.Command{
	Use:   "cnpj",
	Short: "Generate or validate CNPJ",
	Example: strings.Join([]string{
		"brdoc cnpj --generate",
		"brdoc cnpj --generate --legacy",
		"brdoc cnpj --generate --count 10",
		"brdoc cnpj --validate 12.345.678/0001-95",
		"brdoc cnpj --validate --from cnpjs.txt",
		"type cnpjs.txt | brdoc cnpj --validate --from -",
	}, "\n"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags combination
		if cnpjGenerate && (cnpjValidate != "" || cnpjFrom != "") {
			return errors.New("--generate cannot be used with --validate or --from")
		}

		if cnpjFrom != "" && cnpjValidate != "" {
			return errors.New("--from and --validate are mutually exclusive for CNPJ")
		}

		if !cnpjGenerate && cnpjValidate == "" && cnpjFrom == "" {
			return errors.New("either --generate, --validate, or --from must be provided")
		}

		c := sdk.NewCNPJ()
		if cnpjGenerate {
			if cnpjCount <= 0 {
				cnpjCount = 1
			}

			w := bufio.NewWriter(cmd.OutOrStdout())
			defer func(w *bufio.Writer) {
				if err := w.Flush(); err != nil {
					panic(err)
				}
			}(w)

			for i := 0; i < cnpjCount; i++ {
				if cnpjLegacy {
					result, _ := c.Format(c.GenerateLegacy())
					_, _ = fmt.Fprintln(w, result)
				} else {
					result, _ := c.Format(c.Generate())
					_, _ = fmt.Fprintln(w, result)
				}
			}

			return nil
		}

		// validate single or bulk
		if cnpjFrom != "" { // bulk from file or stdin
			r, closeFn, err := openReader(cnpjFrom)
			if err != nil {
				return err
			}

			if closeFn != nil {
				defer closeFn()
			}

			scanner := bufio.NewScanner(r)
			scanner.Buffer(buf, maxLine)

			w := bufio.NewWriter(cmd.OutOrStdout())
			defer func(w *bufio.Writer) {
				if err := w.Flush(); err != nil {
					panic(err)
				}
			}(w)

			anyInvalid := false
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				if c.Validate(line) {
					if formatted, err := c.Format(line); err == nil {
						_, _ = fmt.Fprintf(w, "valid\t%s\n", formatted)
					} else {
						_, _ = fmt.Fprintln(w, "valid")
					}
				} else {
					anyInvalid = true
					_, _ = fmt.Fprintf(w, "invalid\t%s\n", line)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			if anyInvalid {
				cmd.SilenceUsage = true
			}

			return nil
		}

		// single validate value
		if c.Validate(cnpjValidate) {
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

// openReader returns an io.Reader for the given path. If a path is "-", it returns stdin.
// The second return value is a close function for file readers (nil for stdin).
func openReader(path string) (io.Reader, func(), error) {
	if path == "-" {
		return os.Stdin, nil, nil
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}

	closeFn := func() { _ = f.Close() }

	return f, closeFn, nil
}
