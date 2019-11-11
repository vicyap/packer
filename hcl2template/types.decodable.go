package hcl2template

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Decodable interface {
	HCL2Spec() map[string]hcldec.Spec
}

func decodeHCL2Spec(block *hcl.Block, ctx *hcl.EvalContext, dec Decodable) (cty.Value, hcl.Diagnostics) {
	spec := dec.HCL2Spec()
	return hcldec.Decode(block.Body, hcldec.ObjectSpec(spec), ctx)
}

type SelfFlattened interface {
	FlatMapstructure() interface{}
}

func unmarshalCty(block *hcl.Block, val cty.Value, dst SelfFlattened) hcl.Diagnostics {
	v := dst.FlatMapstructure()
	err := gocty.FromCtyValue(val, v)

	diags := hcl.Diagnostics{}
	if err != nil {
		switch err := err.(type) {
		case cty.PathError:
			diags = append(diags, &hcl.Diagnostic{
				Summary: "gocty.FromCtyValue: " + err.Error(),
				Subject: &block.DefRange,
				Detail:  fmt.Sprintf("%v", err.Path),
			})
		default:
			diags = append(diags, &hcl.Diagnostic{
				Summary: "gocty.FromCtyValue: " + err.Error(),
				Subject: &block.DefRange,
				Detail:  fmt.Sprintf("%v", err),
			})
		}
	}
	return diags
}
