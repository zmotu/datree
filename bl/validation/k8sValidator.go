package validation

import (
	"fmt"
	"io"
	"os"

	"github.com/datreeio/datree/pkg/extractor"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type ValidationClient interface {
	Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result
}

type K8sValidator struct {
	validationClient ValidationClient
}

func New() *K8sValidator {
	return &K8sValidator{}
}

type ValidationStatus int

const (
	InvalidYamlFile ValidationStatus = iota + 1
	InvalidK8sFile
)

func (validationStatus ValidationStatus) String() string {
	return [...]string{"InvalidYamlFile", "InvalidK8sFile"}[validationStatus]
}

type InvalidFile struct {
	Path             string
	ValidationStatus ValidationStatus
	ValidationErrors []error
}

func (val *K8sValidator) InitClient(k8sVersion string) {
	val.validationClient = newKubconformValidator(k8sVersion)
}

func (val *K8sValidator) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *InvalidFile) {
	validK8sFilesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidK8sFilesChan := make(chan *InvalidFile, concurrency)

	go func() {
		defer func() {
			close(invalidK8sFilesChan)
			close(validK8sFilesConfigurationsChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {

			isValid, validationErrors, err := val.validateResource(fileConfigurations.FileName)
			if err != nil {
				invalidK8sFilesChan <- &InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationStatus: InvalidK8sFile,
					ValidationErrors: []error{err},
				}
				continue
			}
			if isValid {
				validK8sFilesConfigurationsChan <- fileConfigurations
			} else {
				invalidK8sFilesChan <- &InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationStatus: InvalidK8sFile,
					ValidationErrors: validationErrors,
				}
			}
		}
	}()
	return validK8sFilesConfigurationsChan, invalidK8sFilesChan
}

func (val *K8sValidator) validateResource(filepath string) (bool, []error, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, []error{}, fmt.Errorf("failed opening %s: %s", filepath, err)
	}

	results := val.validationClient.Validate(filepath, f)
	isValid := true
	var validationErrors []error
	for _, res := range results {

		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			isValid = false
			validationErrors = append(validationErrors, res.Err)
		}
	}

	return isValid, validationErrors, nil
}

func newKubconformValidator(k8sVersion string) ValidationClient {
	v, _ := kubeconformValidator.New(nil, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion})
	return v
}