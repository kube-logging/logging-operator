// Copyright (c) 2020 Banzai Cloud Zrt. All Rights Reserved.

package annotation

import (
	"regexp"
	"strings"
	"unicode"
)

// Config .
type Config struct {
	TailerDescriptorDelimiter             rune
	TailerAnnotationDelimiter             rune
	TailerDescriptorValidatorRegexPattern string
}

var defaults = Config{
	TailerDescriptorDelimiter:             ':',
	TailerAnnotationDelimiter:             ',',
	TailerDescriptorValidatorRegexPattern: "^([a-zA-Z0-9-_]+:)?(/[a-zA-Z-0-9._-]+)+$",
}

// TailerDescriptor alias to string with the format of "containername:absolutepath" or "absolutepath"
type TailerDescriptor = string

// TailerAnnotation is a set of Descriptors separated by commas
type TailerAnnotation = TailerDescriptor

// FilePaths .
type FilePaths = []string

// ContainerPaths .
type ContainerPaths map[string]FilePaths

// Handler .
type Handler struct {
	Config
	containerPaths       ContainerPaths
	defaultContainerName string
}

// NewHandler is a custom constructor which receives the available container names
func NewHandler(containerNames []string) *Handler {
	h := &Handler{
		containerPaths: make(ContainerPaths, len(containerNames)),
		Config:         defaults,
	}
	if len(containerNames) > 0 {
		h.defaultContainerName = containerNames[0]
	}
	for _, containerName := range containerNames {
		h.containerPaths[containerName] = []string{}
	}
	return h
}

func (h *Handler) addTailerDescriptor(tailerDescriptor TailerDescriptor) {
	if h == nil || h.defaultContainerName == "" || len(h.containerPaths) == 0 {
		return
	}

	tailerDescriptor = strings.TrimFunc(tailerDescriptor, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if !h.validateTailerDescriptor(tailerDescriptor) {
		return
	}

	elements := strings.FieldsFunc(tailerDescriptor, func(r rune) bool {
		return r == h.TailerDescriptorDelimiter
	})

	var containerName = h.defaultContainerName
	switch len(elements) {
	case 2:
		_, ok := h.containerPaths[elements[0]]
		if !ok {
			return
		}
		containerName = elements[0]
		fallthrough
	case 1:
		h.containerPaths[containerName] = append(h.containerPaths[containerName], elements[len(elements)-1])
	case 0:
		fallthrough
	default:
		return
	}
}

// AddTailerAnnotation .
func (h *Handler) AddTailerAnnotation(tailerAnnotation TailerAnnotation) {
	descriptorStrings := strings.FieldsFunc(tailerAnnotation, func(r rune) bool {
		return r == h.TailerAnnotationDelimiter
	})
	for _, descriptorString := range descriptorStrings {
		h.addTailerDescriptor(descriptorString)
	}
}

// FilePathsForContainer returns FilePaths for given container
func (h *Handler) FilePathsForContainer(containerName string) FilePaths {
	if h == nil {
		return nil
	}
	if containerName == "" {
		containerName = h.defaultContainerName
	}
	if paths, ok := h.containerPaths[containerName]; ok {
		return paths
	}
	return FilePaths{}
}

// AllFilePaths returns FilePaths for all containers
func (h *Handler) AllFilePaths() FilePaths {
	if h == nil {
		return nil
	}
	result := FilePaths{}
	for _, v := range h.containerPaths {
		result = append(result, v...)
	}
	return result
}

func (h *Handler) validateTailerDescriptor(tailerDescriptor TailerDescriptor) bool {
	res, _ := regexp.MatchString(h.TailerDescriptorValidatorRegexPattern, tailerDescriptor)
	// res defaults to false on error
	return res
}
