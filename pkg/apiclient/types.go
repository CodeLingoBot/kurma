// Copyright 2015-2016 Apcera Inc. All rights reserved.

package apiclient

import (
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"

	"github.com/apcera/kurma/pkg/image"
	ntypes "github.com/apcera/kurma/pkg/networkmanager/types"
	kschema "github.com/apcera/kurma/schema"
)

type Pod struct {
	UUID     string              `json:"uuid"`
	Name     string              `json:"name"`
	Pod      *schema.PodManifest `json:"pod"`
	Networks []*ntypes.IPResult  `json:"networks"`
	State    State               `json:"state"`
}

type Image struct {
	Hash     string                `json:"hash"`
	Manifest *schema.ImageManifest `json:"manifest"`
	Size     int64                 `json:"size"`
}

type PodCreateRequest struct {
	Name            string              `json:"name"`
	Pod             *schema.PodManifest `json:"pod"`
	Networks        []string            `json:"networks,omitempty"`
	StagerImageHash string              `json:"stagerImageHash,omitempty"`
}

type PodListResponse struct {
	Pods []*Pod `json:"pods"`
}

type PodResponse struct {
	Pod *Pod `json:"pod"`
}

type ContainerEnterRequest struct {
	UUID    string         `json:"uuid"`
	AppName string         `json:"appName"`
	App     kschema.RunApp `json:"app"`
}

type ContainerEnterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ImageListResponse struct {
	Images []*Image `json:"images"`
}

type ImageFetchRequest struct {
	ImageURI string `json:"imageUri"`
	*image.FetchConfig
}

type ImageResponse struct {
	Image *Image `json:"image"`
}

type None struct{}

type State string

const (
	STATE_NEW      = State("NEW")
	STATE_STARTING = State("STARTING")
	STATE_RUNNING  = State("RUNNING")
	STATE_STOPPING = State("STOPPING")
	STATE_STOPPED  = State("STOPPED")
	STATE_EXITED   = State("EXITED")
)

type HostInfo struct {
	Hostname      string       `json:"hostname"`
	Cpus          int          `json:"cpus"`
	Memory        int64        `json:"memory"`
	Platform      string       `json:"platform"`
	Arch          string       `json:"arch"`
	ACVersion     types.SemVer `json:"ac_version"`
	KurmaVersion  types.SemVer `json:"kurma_version"`
	KernelVersion string       `json:"kernel_version"`
}
