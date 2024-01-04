// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/web-platform-tests/wpt.fyi/shared (interfaces: WebFeaturesManifestDownloader,WebFeatureManifestParser)

// Package sharedtest is a generated GoMock package.
package sharedtest

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	shared "github.com/web-platform-tests/wpt.fyi/shared"
)

// MockWebFeaturesManifestDownloader is a mock of WebFeaturesManifestDownloader interface.
type MockWebFeaturesManifestDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockWebFeaturesManifestDownloaderMockRecorder
}

// MockWebFeaturesManifestDownloaderMockRecorder is the mock recorder for MockWebFeaturesManifestDownloader.
type MockWebFeaturesManifestDownloaderMockRecorder struct {
	mock *MockWebFeaturesManifestDownloader
}

// NewMockWebFeaturesManifestDownloader creates a new mock instance.
func NewMockWebFeaturesManifestDownloader(ctrl *gomock.Controller) *MockWebFeaturesManifestDownloader {
	mock := &MockWebFeaturesManifestDownloader{ctrl: ctrl}
	mock.recorder = &MockWebFeaturesManifestDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebFeaturesManifestDownloader) EXPECT() *MockWebFeaturesManifestDownloaderMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockWebFeaturesManifestDownloader) Download(arg0 context.Context) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockWebFeaturesManifestDownloaderMockRecorder) Download(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockWebFeaturesManifestDownloader)(nil).Download), arg0)
}

// MockWebFeatureManifestParser is a mock of WebFeatureManifestParser interface.
type MockWebFeatureManifestParser struct {
	ctrl     *gomock.Controller
	recorder *MockWebFeatureManifestParserMockRecorder
}

// MockWebFeatureManifestParserMockRecorder is the mock recorder for MockWebFeatureManifestParser.
type MockWebFeatureManifestParserMockRecorder struct {
	mock *MockWebFeatureManifestParser
}

// NewMockWebFeatureManifestParser creates a new mock instance.
func NewMockWebFeatureManifestParser(ctrl *gomock.Controller) *MockWebFeatureManifestParser {
	mock := &MockWebFeatureManifestParser{ctrl: ctrl}
	mock.recorder = &MockWebFeatureManifestParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebFeatureManifestParser) EXPECT() *MockWebFeatureManifestParserMockRecorder {
	return m.recorder
}

// Parse mocks base method.
func (m *MockWebFeatureManifestParser) Parse(arg0 context.Context, arg1 io.ReadCloser) (shared.WebFeaturesData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0, arg1)
	ret0, _ := ret[0].(shared.WebFeaturesData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *MockWebFeatureManifestParserMockRecorder) Parse(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockWebFeatureManifestParser)(nil).Parse), arg0, arg1)
}
