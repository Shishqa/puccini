package csar

import (
	"fmt"
	"strconv"

	"github.com/tliron/exturl"
)

func NewURL(csarUrl exturl.URL, format string, path string) (exturl.URL, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	if exturl.IsValidTarballArchiveFormat(format) {
		return exturl.NewTarballURL(path, csarUrl, format), nil
	}

	switch format {
	case "zip", "csar":
		return exturl.NewZipURL(path, csarUrl), nil
	default:
		return nil, fmt.Errorf("unsupported CSAR archive format: %q", format)
	}
}

func GetDefaultServiceTemplateURL(csarUrl exturl.URL, format string) (exturl.URL, error) {
	return GetServiceTemplateURL(csarUrl, format, "")
}

func GetServiceTemplateURL(csarUrl exturl.URL, format string, serviceTemplateName string) (exturl.URL, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	meta, err := ReadMetaFromURL(csarUrl, format)
	if err != nil {
		if exturl.IsNotFound(err) {
			if meta, err = NewMetaFor(csarUrl, format); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if (serviceTemplateName == "") || (serviceTemplateName == "0") {
		// Default entry point

		// Attempt to use Entry-Definitions in TOSCA.meta
		if meta.EntryDefinitions != "" {
			return NewURL(csarUrl, format, meta.EntryDefinitions)
		}

		// Attempt to find it in root of archive
		if path, err := GetRootPath(csarUrl, format); err == nil {
			return NewURL(csarUrl, format, path)
		} else {
			return nil, err
		}
	} else {
		// Alternative entry points

		// Try as integer
		if serviceTemplateNumber, err := strconv.ParseUint(serviceTemplateName, 10, 64); err == nil {
			if otherDefinitionIndex := int(serviceTemplateNumber) - 1; (otherDefinitionIndex >= 0) && (otherDefinitionIndex < len(meta.OtherDefinitions)) {
				return NewURL(csarUrl, format, meta.OtherDefinitions[otherDefinitionIndex])
			}
		}

		// Try as string
		for _, otherDefinition := range meta.OtherDefinitions {
			if otherDefinition == serviceTemplateName {
				return NewURL(csarUrl, format, serviceTemplateName)
			}
		}
	}

	return nil, fmt.Errorf("CSAR does not have service template %q: %s", serviceTemplateName, csarUrl.String())
}
