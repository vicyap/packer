package hcl2template

import (
	"github.com/mitchellh/mapstructure"
	"encoding/json"

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
	Configure(...interface{}) error
}

func unmarshalCty(val cty.Value, dst SelfFlattened) error {
	flat := dst.FlatMapstructure()
	err := gocty.FromCtyValue(val, flat)
	if err != nil {
		return err
	}

	// Currently we json encode the "flat structure" to later mapstructure-load
	// it onto dst.
	//
	// May be another solution to this would be to set dst from v by simply
	// typechecking v and dst.
	b, err := json.Marshal(flat)
	if err != nil {
		return err
	}
	mapstructure.Decode
	dst.Configure(...interface{})
}
