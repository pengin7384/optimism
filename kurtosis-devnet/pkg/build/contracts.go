package build

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"text/template"
)

// ContractBuilder handles building smart contracts using just commands
type ContractBuilder struct {
	// Base directory where the build commands should be executed
	baseDir string
	// Template for the build command
	cmdTemplate *template.Template

	// Dry run mode
	dryRun bool
}

const (
	contractsCmdTemplateStr = "just _contracts-build {{.BundlePath}}"
)

var defaultContractTemplate *template.Template

func init() {
	defaultContractTemplate = template.Must(template.New("contract_build_cmd").Parse(contractsCmdTemplateStr))
}

type ContractBuilderOptions func(*ContractBuilder)

func WithContractBaseDir(baseDir string) ContractBuilderOptions {
	return func(b *ContractBuilder) {
		b.baseDir = baseDir
	}
}

func WithContractTemplate(cmdTemplate *template.Template) ContractBuilderOptions {
	return func(b *ContractBuilder) {
		b.cmdTemplate = cmdTemplate
	}
}

func WithContractDryRun(dryRun bool) ContractBuilderOptions {
	return func(b *ContractBuilder) {
		b.dryRun = dryRun
	}
}

// NewContractBuilder creates a new ContractBuilder instance
func NewContractBuilder(opts ...ContractBuilderOptions) *ContractBuilder {
	b := &ContractBuilder{
		baseDir:     ".",
		cmdTemplate: defaultContractTemplate,
		dryRun:      false,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// templateData holds the data for the command template
type contractTemplateData struct {
	Layer      string
	BundlePath string
}

// Build executes the contract build command
func (b *ContractBuilder) Build(layer string, bundlePath string) error {
	log.Printf("Building contracts bundle: %s", bundlePath)

	// Prepare template data
	data := contractTemplateData{
		Layer:      layer,
		BundlePath: bundlePath,
	}

	// Execute template to get command string
	var cmdBuf bytes.Buffer
	if err := b.cmdTemplate.Execute(&cmdBuf, data); err != nil {
		return fmt.Errorf("failed to execute command template: %w", err)
	}

	// Create command
	cmd := exec.Command("sh", "-c", cmdBuf.String())
	cmd.Dir = b.baseDir

	if !b.dryRun {
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("contract build command failed: %w\nOutput: %s", err, string(output))
		}
	}

	return nil
}