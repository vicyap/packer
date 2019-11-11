package packer

import "github.com/hashicorp/hcl2/hcldec"

type HCL2Speccer interface {
	// HCL2Spec should give the hcl spec used to configure the builder. It will
	// be used to tell the HCL reading library how to validate/configure the
	// builder. After the HCL config file is parsed Prepare will be called with
	// the loaded configuration.
	HCL2Spec() map[string]hcldec.Spec
}
