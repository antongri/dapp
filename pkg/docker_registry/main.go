package docker_registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

type RepoImage struct {
	Repository string
	Tag        string
	v1.Image
}

var GCRUrlPatterns = []string{"^container\\.cloud\\.google\\.com", "^gcr\\.io", "^.*\\.gcr\\.io"}

func IsGCR(reference string) (bool, error) {
	u, err := url.Parse(fmt.Sprintf("scheme://%s", reference))
	if err != nil {
		return false, err
	}

	for _, pattern := range GCRUrlPatterns {
		matched, err := regexp.MatchString(pattern, u.Host)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func DimgTags(reference string) ([]string, error) {
	images, err := ImagesByDappDimgLabel(reference, "true")
	if err != nil {
		return nil, err
	}

	return imagesTags(images), nil
}

func DimgstageTags(reference string) ([]string, error) {
	images, err := ImagesByDappDimgLabel(reference, "false")
	if err != nil {
		return nil, err
	}

	return imagesTags(images), nil
}

func imagesTags(images []RepoImage) []string {
	if len(images) == 0 {
		return []string{}
	} else {
		var tags []string
		for _, image := range images {
			tags = append(tags, image.Tag)
		}

		return tags
	}
}

func ImagesByDappDimgLabel(reference, labelValue string) ([]RepoImage, error) {
	var repoImages []RepoImage

	tags, err := Tags(reference)
	if err != nil {
		if strings.Contains(err.Error(), "NAME_UNKNOWN") {
			return []RepoImage{}, nil
		}
		return nil, err
	}

	for _, tag := range tags {
		tagReference := strings.Join([]string{reference, tag}, ":")
		v1Image, _, err := image(tagReference)
		if err != nil {
			if strings.Contains(err.Error(), "BLOB_UNKNOWN") {
				fmt.Printf("Ignore broken tag '%s': %s\n", tag, err)
				continue
			}
			return nil, err
		}

		configFile, err := v1Image.ConfigFile()
		if err != nil {
			if strings.Contains(err.Error(), "MANIFEST_UNKNOWN") {
				fmt.Printf("Ignore broken tag '%s': %s\n", tag, err)
				continue
			}
			return nil, err
		}

		for k, v := range configFile.Config.Labels {
			if k == "dapp-dimg" && v == labelValue {
				repoImage := RepoImage{
					Repository: reference,
					Tag:        tag,
					Image:      v1Image,
				}

				repoImages = append(repoImages, repoImage)
				break
			}
		}
	}

	return repoImages, nil
}

func Tags(reference string) ([]string, error) {
	repo, err := name.NewRepository(reference, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("parsing repo %q: %v", reference, err)
	}

	auth, err := authn.DefaultKeychain.Resolve(repo.Registry)
	if err != nil {
		return nil, fmt.Errorf("getting creds for %q: %v", repo, err)
	}

	tags, err := remote.List(repo, auth, http.DefaultTransport)
	if err != nil {
		return nil, fmt.Errorf("reading tags for %q: %v", repo, err)
	}

	return tags, nil
}

func ImageId(reference string) (string, error) {
	i, _, err := image(reference)
	if err != nil {
		return "", err
	}

	manifest, err := i.Manifest()
	if err != nil {
		return "", err
	}

	return manifest.Config.Digest.String(), nil
}

func ImageParentId(reference string) (string, error) {
	configFile, err := ImageConfigFile(reference)
	if err != nil {
		return "", err
	}

	return configFile.ContainerConfig.Image, nil
}

func ImageConfigFile(reference string) (v1.ConfigFile, error) {
	i, _, err := image(reference)
	if err != nil {
		return v1.ConfigFile{}, err
	}

	configFile, err := i.ConfigFile()
	if err != nil {
		return v1.ConfigFile{}, err
	}

	return *configFile, nil
}

func ImageDelete(reference string) error {
	r, err := name.ParseReference(reference, name.WeakValidation)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %v", reference, err)
	}

	auth, err := authn.DefaultKeychain.Resolve(r.Context().Registry)
	if err != nil {
		return fmt.Errorf("getting creds for %q: %v", r, err)
	}

	if err := remote.Delete(r, auth, http.DefaultTransport); err != nil {
		if strings.Contains(err.Error(), "UNAUTHORIZED") {
			if gitlabRegistryDeleteErr := GitlabRegistryDelete(r, auth, http.DefaultTransport); gitlabRegistryDeleteErr != nil {
				if strings.Contains(gitlabRegistryDeleteErr.Error(), "UNAUTHORIZED") {
					return fmt.Errorf("deleting image %q: %v", r, err)
				}
				return fmt.Errorf("deleting image %q: %v", r, gitlabRegistryDeleteErr)
			}
		} else {
			return fmt.Errorf("deleting image %q: %v", r, err)
		}
	}

	return nil
}

// TODO https://gitlab.com/gitlab-org/gitlab-ce/issues/48968
func GitlabRegistryDelete(ref name.Reference, auth authn.Authenticator, t http.RoundTripper) error {
	scopes := []string{ref.Scope("*")}
	tr, err := transport.New(ref.Context().Registry, auth, t, scopes)
	if err != nil {
		return err
	}
	c := &http.Client{Transport: tr}

	u := url.URL{
		Scheme: ref.Context().Registry.Scheme(),
		Host:   ref.Context().RegistryStr(),
		Path:   fmt.Sprintf("/v2/%s/manifests/%s", ref.Context().RepositoryStr(), ref.Identifier()),
	}

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unrecognized status code during DELETE: %v; %v", resp.Status, string(b))
	}
}

func ImageDigest(reference string) (string, error) {
	i, _, err := image(reference)
	if err != nil {
		return "", err
	}

	digest, err := i.Digest()
	if err != nil {
		return "", err
	}

	return digest.String(), nil
}

func image(reference string) (v1.Image, name.Reference, error) {
	ref, err := name.ParseReference(reference, name.WeakValidation)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing reference %q: %v", reference, err)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, nil, fmt.Errorf("reading image %q: %v", ref, err)
	}

	return img, ref, nil
}
