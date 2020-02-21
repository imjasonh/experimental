package oci

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// Bundle bundles resources into an Image.
func Bundle(resources []ParsedTektonResource) (v1.Image, error) {
	img := empty.Image
	for _, r := range resources {
		l, err := tarball.LayerFromReader(strings.NewReader(r.Contents))
		if err != nil {
			return nil, fmt.Errorf("Error creating layer for resource %s/%s: %w", r.Kind, r.Name, err)
		}
		img, err = mutate.Append(img, mutate.Addendum{
			// TODO: Specify custom layer media type ("application/vnd.cdf.tekton.catalog.v1alpha1+json")
			Layer: l,
			Annotations: map[string]string{
				"org.opencontainers.image.title": getLayerName(r.Kind.Kind, r.Name),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("Error appending resource %q: %w", r.Name, err)
		}
	}
	return img, nil
}

// Extract pulls out config bytes from the image for an object with the given
// kind and name.
func Extract(img v1.Image, kind, name string) ([]byte, error) {
	m, err := img.Manifest()
	if err != nil {
		return nil, err
	}
	ls, err := img.Layers()
	if err != nil {
		return nil, err
	}
	var layer v1.Layer
	for idx, l := range m.Layers {
		// TODO: Check for custom media type.
		if l.Annotations["org.opencontainers.image.title"] == getLayerName(kind, name) {
			layer = ls[idx]
			break
		}
	}
	if layer == nil {
		return nil, fmt.Errorf("Resource %s/%s not found", kind, name)
	}
	rc, err := layer.Uncompressed()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func getLayerName(kind string, name string) string {
	return fmt.Sprintf("%s/%s", strings.ToLower(kind), name)
}
