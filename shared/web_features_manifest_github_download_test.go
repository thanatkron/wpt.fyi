// Copyright 2024 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build small

package shared

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var compressedWebFeaturesManifestFilePath = filepath.Join("web_features_manifest_testdata", "WEB_FEATURES_MANIFEST.json.gz")

func createWebFeaturesTestdata() {
	v1Manifest := struct {
		Version int                 `json:"version,omitempty"`
		Data    map[string][]string `json:"data,omitempty"`
	}{
		Version: 1,
		Data: map[string][]string{
			"grid":    {"test1.js", "test2.js"},
			"subgrid": {"test3.js", "test4.js"},
		},
	}
	jsonData, err := json.Marshal(v1Manifest)
	if err != nil {
		panic(err)
	}

	// Create a buffer for compressing the JSON
	var buf bytes.Buffer

	// Create a gzip writer and write the JSON to it
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonData); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	// Write the compressed data to a file
	if err := os.WriteFile(compressedWebFeaturesManifestFilePath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func TestResponseBodyTransformer_Success(t *testing.T) {
	updateGolden := false // Switch this when we want to update the golden file
	if updateGolden {
		createWebFeaturesTestdata()
	}
	f, err := os.Open(compressedWebFeaturesManifestFilePath)
	defer f.Close()
	require.NoError(t, err)

	transformer := gzipBodyTransformer{}
	reader, err := transformer.Transform(f)
	defer reader.Close()
	require.NoError(t, err)

	rawBytes, err := io.ReadAll(reader)
	require.NoError(t, err)

	assert.Equal(t, `{"version":1,"data":{"grid":["test1.js","test2.js"],"subgrid":["test3.js","test4.js"]}}`, string(rawBytes))
}

type RoundTripFunc struct {
	function func(req *http.Request) *http.Response
	err      error
}

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.function(req), f.err
}

type mockBodyTransformerInput struct {
	expectedBody string
	output       io.ReadCloser
	err          error
}

type mockBodyTransformer struct {
	t *testing.T
	mockBodyTransformerInput
}

func (tr mockBodyTransformer) Transform(body io.Reader) (io.ReadCloser, error) {
	bodyBytes, err := io.ReadAll(body)
	require.NoError(tr.t, err)
	assert.Equal(tr.t, tr.expectedBody, string(bodyBytes))
	return tr.output, tr.err
}

func TestGitHubWebFeaturesManifestDownloader_Download(t *testing.T) {
	// Test cases for Download
	tests := []struct {
		name             string
		getLatestRelease func(http.ResponseWriter, *http.Request)
		roundTrip        RoundTripFunc
		transformer      mockBodyTransformerInput
		expectedBody     []byte
		expectedError    error
	}{
		{
			name: "successful download",
			getLatestRelease: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/repos/web-platform-tests/wpt/releases/latest", r.URL.Path)
				release := &github.RepositoryRelease{
					Assets: []*github.ReleaseAsset{
						{
							Name:               github.String("WEB_FEATURES_MANIFEST.json.gz"),
							BrowserDownloadURL: github.String("https://example.com/WEB_FEATURES_MANIFEST.json.gz"),
						},
					},
				}
				w.Write(mock.MustMarshal(release))
			},
			roundTrip: RoundTripFunc{function: func(req *http.Request) *http.Response {
				assert.Equal(t, "https://example.com/WEB_FEATURES_MANIFEST.json.gz", req.URL.String())
				return &http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: int64(len("raw data")),
					Body:          io.NopCloser(bytes.NewBufferString("raw data")),
				}
			}, err: nil},
			transformer: mockBodyTransformerInput{
				expectedBody: "raw data",
				output:       io.NopCloser(bytes.NewBufferString("transformed data")),
				err:          nil,
			},
			expectedBody:  []byte("transformed data"),
			expectedError: nil,
		},
		{
			name: "error getting latest release",
			getLatestRelease: func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					"failed to get release",
				)
			},
			expectedBody:  nil,
			expectedError: ErrUnableToRetrieveGitHubRelease,
		},
		{
			name: "manifest file not found",
			getLatestRelease: func(w http.ResponseWriter, r *http.Request) {
				release := &github.RepositoryRelease{
					Assets: []*github.ReleaseAsset{},
				}
				w.Write(mock.MustMarshal(release))
			},
			expectedBody:  nil,
			expectedError: ErrNoWebFeaturesManifestFileFound,
		},
		{
			name: "error downloading asset",
			getLatestRelease: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/repos/web-platform-tests/wpt/releases/latest", r.URL.Path)
				release := &github.RepositoryRelease{
					Assets: []*github.ReleaseAsset{
						{
							Name:               github.String("WEB_FEATURES_MANIFEST.json.gz"),
							BrowserDownloadURL: github.String("https://example.com/WEB_FEATURES_MANIFEST.json.gz"),
						},
					},
				}
				w.Write(mock.MustMarshal(release))
			},
			roundTrip: RoundTripFunc{function: func(req *http.Request) *http.Response {
				assert.Equal(t, "https://example.com/WEB_FEATURES_MANIFEST.json.gz", req.URL.String())
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}
			}, err: errors.New("simulated network error")},
			expectedBody:  nil,
			expectedError: ErrGitHubAssetDownloadFailedToComplete,
		},
		{
			name: "empty response body",
			getLatestRelease: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/repos/web-platform-tests/wpt/releases/latest", r.URL.Path)
				release := &github.RepositoryRelease{
					Assets: []*github.ReleaseAsset{
						{
							Name:               github.String("WEB_FEATURES_MANIFEST.json.gz"),
							BrowserDownloadURL: github.String("https://example.com/WEB_FEATURES_MANIFEST.json.gz"),
						},
					},
				}
				w.Write(mock.MustMarshal(release))
			},
			roundTrip: RoundTripFunc{function: func(req *http.Request) *http.Response {
				assert.Equal(t, "https://example.com/WEB_FEATURES_MANIFEST.json.gz", req.URL.String())
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       nil,
				}
			}, err: nil},
			expectedBody:  nil,
			expectedError: ErrMissingBodyDuringWebFeaturesManifestDownload,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockedGitHubHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetReposReleasesLatestByOwnerByRepo,
					http.HandlerFunc(tc.getLatestRelease),
				),
			)
			httpClient := &http.Client{
				Transport: tc.roundTrip,
			}
			c := github.NewClient(mockedGitHubHTTPClient)
			downloader := NewGitHubWebFeaturesManifestDownloader(httpClient, c)
			downloader.bodyTransformer = mockBodyTransformer{t, tc.transformer}
			body, err := downloader.Download(context.Background())
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Download() returned unexpected error: (%v). expected error: (%v).", err, tc.expectedError)
			}

			// No need to compare the body if there's an error
			if err != nil {
				return
			}

			bodyBytes, err := io.ReadAll(body)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, bodyBytes)
		})
	}
}
