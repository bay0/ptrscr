package gist

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/github"
)

// Create pushes file to github gist
func Create(client *github.Client, bytes []byte, filename string) (*github.Gist, error) {
	ctx := context.Background()
	f := make(map[github.GistFilename]github.GistFile)

	base64Img := "data:image/png;base64,"
	base64Img += base64.StdEncoding.EncodeToString(bytes)

	f[github.GistFilename(filename)] = github.GistFile{
		Content: github.String(base64Img),
		Size:    github.Int(len(bytes)),
	}
	gist := &github.Gist{
		Description: github.String("Enemy Eater Snap!"),
		Public:      github.Bool(false),
		Files:       f,
	}
	gistResponse, _, err := client.Gists.Create(ctx, gist)
	return gistResponse, err
}
