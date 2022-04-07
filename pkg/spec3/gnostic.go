/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spec3

import (
	"errors"

	openapi_v3 "github.com/google/gnostic/openapiv3"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func (k *OpenAPI) FromGnostic(g *openapi_v3.Document) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	k.Version = g.Openapi

	if info := g.Info; info != nil {
		k.Info = &spec.Info{}
		spec2VendorExtensibleFromGnostic(&k.Info.VendorExtensible, info.SpecificationExtension)

		k.Info.Description = info.Description
		k.Info.Title = info.Description
		k.Info.TermsOfService = info.TermsOfService
		if contact := info.Contact; contact != nil {
			k.Info.Contact = &spec.ContactInfo{
				Name:  contact.Name,
				URL:   contact.Url,
				Email: contact.Email,
			}

			// data loss! contact.SpecificationExtension
			if len(contact.SpecificationExtension) > 0 {
				ok = false
			}
		}
		if license := info.License; license != nil {
			k.Info.License = &spec.License{
				Name: license.Name,
				URL:  license.Url,
			}
			// data loss!: license.SpecificationExtension
			if len(license.SpecificationExtension) > 0 {
				ok = false
			}
		}
		k.Info.Version = info.Version
	}

	return ok, nil
}

func (k *Paths) FromGnostic(g *openapi_v3.Paths) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	if paths := g.Path; paths != nil {
		k.Paths = make(map[string]*Path)

		for _, path := range paths {
			if path == nil {
				continue
			}

			converted := &Path{}
			if nok, err := converted.FromGnostic(path.Value); err != nil {
				return false, err
			} else {
				ok = ok && nok
			}

			k.Paths[path.Name] = converted
		}
	}

	return ok, nil
}

func (k *Path) FromGnostic(g *openapi_v3.PathItem) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2RefableFromGnostic(&k.Refable, g.XRef); err != nil {
		return false, err
	}

	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	k.Summary = g.Summary
	k.Description = g.Description

	if g.Get != nil {
		k.Get = &Operation{}
		if nok, err := k.Get.FromGnostic(g.Get); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Put != nil {
		k.Put = &Operation{}
		if nok, err := k.Put.FromGnostic(g.Put); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Post != nil {
		k.Post = &Operation{}
		if nok, err := k.Post.FromGnostic(g.Post); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Delete != nil {
		k.Delete = &Operation{}
		if nok, err := k.Delete.FromGnostic(g.Delete); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Options != nil {
		k.Options = &Operation{}
		if nok, err := k.Options.FromGnostic(g.Options); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Head != nil {
		k.Head = &Operation{}
		if nok, err := k.Head.FromGnostic(g.Head); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Patch != nil {
		k.Patch = &Operation{}
		if nok, err := k.Patch.FromGnostic(g.Patch); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Trace != nil {
		k.Trace = &Operation{}
		if nok, err := k.Trace.FromGnostic(g.Trace); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Servers != nil {
		k.Servers = make([]*Server, len(g.Servers))
		for _, server := range g.Servers {
			if server == nil {
				continue
			}

			converted := &Server{}
			if nok, err := converted.FromGnostic(server); err != nil {
				return false, err
			} else if !nok {
				ok = false
			}

			k.Servers = append(k.Servers, converted)
		}
	}

	if g.Parameters != nil {
		k.Parameters = make([]*Parameter, len(g.Parameters))

		for _, parameter := range g.Parameters {
			if parameter == nil {
				continue
			}

			converted := &Parameter{}
			if nok, err := converted.FromGnostic(parameter); err != nil {
				return false, err
			} else if !nok {
				ok = false
			}

			k.Parameters = append(k.Parameters, converted)
		}
	}

	return ok, nil
}

func (k *Parameter) FromGnostic(g *openapi_v3.ParameterOrReference) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	switch p := g.GetOneof().(type) {
	case *openapi_v3.ParameterOrReference_Parameter:
		param := p.Parameter

		if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, param.SpecificationExtension); err != nil {
			return false, err
		}

		k.Name = param.Name
		k.In = param.In
		k.Description = param.Description
		k.Required = param.Required
		k.Deprecated = param.Deprecated
		k.AllowEmptyValue = param.AllowEmptyValue
		k.Style = param.Style
		k.Explode = param.Explode
		k.AllowReserved = param.AllowReserved

		if param.Example != nil {
			if err := param.Example.ToRawInfo().Decode(&k.Example); err != nil {
				return false, err
			}
		}

		if param.Examples != nil {
			k.Examples = make(map[string]*Example)
			for _, example := range param.Examples.GetAdditionalProperties() {
				if example == nil {
					continue
				}

				converted := &Example{}
				if nok, err := converted.FromGnostic(example.GetValue()); err != nil {
					return false, err
				} else if !nok {
					ok = false
				}

				k.Examples[example.Name] = converted
			}
		}

	case *openapi_v3.ParameterOrReference_Reference:
		if err := spec2RefableFromGnostic(&k.Refable, p.Reference.XRef); err != nil {
			return false, err
		}

		k.Description = p.Reference.Description

		// data loss! p.Reference.Summary
		if len(p.Reference.Summary) > 0 {
			ok = false
		}
	default:
		return false, errors.New("unrecognized parameter type")
	}

	return ok, nil
}

func (k *Example) FromGnostic(g *openapi_v3.ExampleOrReference) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true

	switch p := g.GetOneof().(type) {
	case *openapi_v3.ExampleOrReference_Example:
		if p.Example == nil {
			break
		}

		if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, p.Example.SpecificationExtension); err != nil {
			return false, err
		}

		k.Summary = p.Example.Summary
		k.Description = p.Example.Description
		k.ExternalValue = p.Example.ExternalValue

		if p.Example.Value != nil {
			if err := p.Example.Value.ToRawInfo().Decode(&k.Value); err != nil {
				return false, err
			}
		}

	case *openapi_v3.ExampleOrReference_Reference:
		if p.Reference == nil {
			break
		}

		if err := spec2RefableFromGnostic(&k.Refable, p.Reference.XRef); err != nil {
			return false, err
		}

		k.Description = p.Reference.Description
		k.Summary = p.Reference.Summary
	default:
		return false, errors.New("unrecognized example type")
	}

	return ok, nil
}

func (k *Operation) FromGnostic(g *openapi_v3.Operation) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	//!RFC: Should this copy?
	k.Tags = g.Tags

	k.Summary = g.Summary
	k.Description = g.Description

	if g.ExternalDocs != nil {
		k.ExternalDocs = &ExternalDocumentation{}
		if nok, err := k.ExternalDocs.FromGnostic(g.ExternalDocs); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	k.OperationId = g.OperationId

	if g.Parameters != nil {
		k.Parameters = make([]*Parameter, len(g.Parameters))
		for _, p := range g.Parameters {
			if p == nil {
				continue
			}
			converted := &Parameter{}
			if nok, err := converted.FromGnostic(p); err != nil {
				return false, err
			} else if !nok {
				ok = false
			}

			k.Parameters = append(k.Parameters, converted)
		}
	}

	if g.RequestBody != nil {
		k.RequestBody = &RequestBody{}
		if nok, err := k.RequestBody.FromGnostic(g.RequestBody); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	if g.Responses != nil {
		k.Responses = &Responses{}
		if nok, err := k.Responses.FromGnostic(g.Responses); err != nil {
			return false, err
		} else if !nok {
			ok = false
		}
	}

	k.Deprecated = g.Deprecated

	if g.Security != nil {
		k.SecurityRequirement = make([]*SecurityRequirement, len(g.Security))
		for _, v := range g.Security {
			if v == nil {
				continue
			}

			converted := &SecurityRequirement{}
			if nok, err := converted.FromGnostic(v); err != nil {
				return false, err
			} else if !nok {
				ok = false
			}

			k.SecurityRequirement = append(k.SecurityRequirement, converted)
		}
	}

	if g.Servers != nil {
		k.Servers = make([]*Server, len(g.Servers))

		for _, v := range g.Servers {
			if v == nil {
				continue
			}

			converted := &Server{}
			if nok, err := converted.FromGnostic(v); err != nil {
				return false, err
			} else if !nok {
				ok = false
			}

			k.Servers = append(k.Servers, converted)
		}
	}

	return ok, nil
}

func (k *Server) FromGnostic(g *openapi_v3.Server) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	k.URL = g.Url
	k.Description = g.Description

	if g.Variables != nil {
		k.Variables = make(map[string]*ServerVariable)
		for _, v := range g.Variables.AdditionalProperties {
			if v == nil {
				continue
			}
			converted := &ServerVariable{}
			if err := converted.FromGnostic(v.GetValue()); err != nil {
				return false, err
			}

			k.Variables[v.GetName()] = converted
		}
	}

	return ok, nil
}

func (k *ServerVariable) FromGnostic(g *openapi_v3.ServerVariable) error {
	if g == nil {
		return nil
	}

	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return err
	}

	k.Description = g.Description
	k.Default = g.Default
	k.Enum = g.Enum

	return nil
}

func (k *Components) FromGnostic(g *openapi_v3.Components) (bool, error) {
	ok := true

	// data loss! g.SpecificationExtension
	if len(g.SpecificationExtension) > 0 {
		ok = false
	}

	if g.Schemas != nil {
		k.Schemas = make(map[string]*spec.Schema)
		for _, v := range g.Schemas.AdditionalProperties {
			if v == nil {
				continue
			}

			converted := &spec.Schema{}
		}
	}

	return ok, nil
}

func (k *ExternalDocumentation) FromGnostic(g *openapi_v3.ExternalDocs) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	return ok, nil
}

func (k *RequestBody) FromGnostic(g *openapi_v3.RequestBodyOrReference) (bool, error) {

}

func (k *Responses) FromGnostic(g *openapi_v3.Responses) (bool, error) {
	if g == nil {
		return true, nil
	}

	ok := true
	if err := spec2VendorExtensibleFromGnostic(&k.VendorExtensible, g.SpecificationExtension); err != nil {
		return false, err
	}

	return ok, nil
}

func (k *SecurityRequirement) FromGnostic(g *openapi_v3.SecurityRequirement) (bool, error) {

}

////////////////////////////////////////////////////////////////////////////////
// OpenAPI V2 Types
// References to these type should be removed from the v3 structs before beta
// goes GA

func spec2VendorExtensibleFromGnostic(k *spec.VendorExtensible, g []*openapi_v3.NamedAny) error {
	if g == nil {
		return nil
	}

	k.Extensions = make(spec.Extensions)

	for _, v := range g {
		if v == nil {
			continue
		}

		var iface interface{}
		if err := v.Value.ToRawInfo().Decode(iface); err != nil {
			return err
		}
		k.Extensions[v.Name] = iface
	}

	return nil
}

func spec2RefableFromGnostic(k *spec.Refable, g string) error {

}
