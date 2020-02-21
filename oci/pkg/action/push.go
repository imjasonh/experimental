package action

import (
	"errors"
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/tektoncd/experimental/oci/pkg/oci"
)

// Push will perform the `push` action by recursively reading all of the
// Tekton specs passed in, bundling it into an image, and pushing the result
// to an OCI-compliant repository.
func Push(r string, filePaths []string) error {
	// Validate the parameters.
	if r == "" || len(filePaths) == 0 {
		return errors.New("must specify a valid image name and file paths")
	}

	ref, err := name.ParseReference(r)
	if err != nil {
		return err
	}

	resources, err := oci.ReadPaths(filePaths)
	if err != nil {
		return err
	}

	img, err := oci.Bundle(resources)
	if err != nil {
		return err
	}

	if err := remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return err
	}

	d, err := img.Digest()
	if err != nil {
		return err
	}
	log.Println("Pushed", ref.Context().Digest(d.String()))
	return nil
}
