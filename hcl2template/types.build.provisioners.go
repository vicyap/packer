package hcl2template

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// Provisioner represents a parsed provisioner
type Provisioner struct {
	// Cfg is a parsed config
	Cfg interface{}
}

type ProvisionerGroup struct {
	CommunicatorRef CommunicatorRef

	Provisioners []Provisioner
	HCL2Ref      HCL2Ref
}

// ProvisionerGroups is a slice of provision blocks; which contains
// provisioners
type ProvisionerGroups []*ProvisionerGroup

func (p *Parser) decodeProvisionerGroup(block *hcl.Block, provisionerSpecs pluginLoader) (*ProvisionerGroup, hcl.Diagnostics) {
	var b struct {
		Communicator string   `hcl:"communicator,optional"`
		Remain       hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	pg := &ProvisionerGroup{}
	pg.CommunicatorRef = communicatorRefFromString(b.Communicator)
	pg.HCL2Ref.DeclRange = block.DefRange

	buildSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
	}
	for _, k := range provisionerSpecs.List() {
		buildSchema.Blocks = append(buildSchema.Blocks, hcl.BlockHeaderSchema{
			Type: k,
		})
	}

	content, moreDiags := b.Remain.Content(buildSchema)
	diags = append(diags, moreDiags...)
	for _, block := range content.Blocks {
		provisioner, err := provisionerSpecs.Get(block.Type)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Summary: "Failed loading " + block.Type,
				Subject: &block.LabelRanges[0],
				Detail:  err.Error(),
			})
			continue
		}
		flatProvisinerCfg, moreDiags := decodeHCL2Spec(block, nil, provisioner)
		diags = append(diags, moreDiags...)
		pg.Provisioners = append(pg.Provisioners, Provisioner{flatProvisinerCfg})
	}

	return pg, diags
}

func (pgs ProvisionerGroups) FirstCommunicatorRef() CommunicatorRef {
	if len(pgs) == 0 {
		return NoCommunicator
	}
	return pgs[0].CommunicatorRef
}
