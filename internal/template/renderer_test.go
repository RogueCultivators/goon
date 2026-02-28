package template

import (
	"strings"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	if renderer == nil {
		t.Fatal("NewRenderer() returned nil renderer")
	}

	if renderer.templates == nil {
		t.Error("NewRenderer() templates is nil")
	}
}

func TestRenderProjectData(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	tests := []struct {
		name         string
		templateName string
		data         ProjectData
		wantContains []string
		wantErr      bool
	}{
		{
			name:         "render go.mod template",
			templateName: "go.mod.tmpl",
			data: ProjectData{
				ProjectName: "testapp",
				ModuleName:  "github.com/user/testapp",
			},
			wantContains: []string{
				"module github.com/user/testapp",
				"go 1.22",
			},
			wantErr: false,
		},
		{
			name:         "render main.go template",
			templateName: "main.go.tmpl",
			data: ProjectData{
				ProjectName: "testapp",
				ModuleName:  "github.com/user/testapp",
			},
			wantContains: []string{
				"package main",
				"github.com/user/testapp",
			},
			wantErr: false,
		},
		{
			name:         "render config.yaml template",
			templateName: "config.yaml.tmpl",
			data: ProjectData{
				ProjectName: "testapp",
				ModuleName:  "github.com/user/testapp",
			},
			wantContains: []string{
				"port:",
				"log:",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderer.Render(tt.templateName, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Render() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Render() unexpected error = %v", err)
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Render() result does not contain %q", want)
				}
			}
		})
	}
}

func TestRenderModuleData(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	tests := []struct {
		name         string
		templateName string
		data         ModuleData
		wantContains []string
		wantErr      bool
	}{
		{
			name:         "render handler template",
			templateName: "handler.go.tmpl",
			data: ModuleData{
				ModuleName:      "user",
				CapitalizedName: "User",
				ProjectModule:   "github.com/user/testapp",
			},
			wantContains: []string{
				"package user",
				"type Handler struct",
				"NewHandler",
			},
			wantErr: false,
		},
		{
			name:         "render service template",
			templateName: "service.go.tmpl",
			data: ModuleData{
				ModuleName:      "product",
				CapitalizedName: "Product",
				ProjectModule:   "github.com/user/testapp",
			},
			wantContains: []string{
				"package product",
				"type Service struct",
				"NewService",
			},
			wantErr: false,
		},
		{
			name:         "render model template",
			templateName: "model.go.tmpl",
			data: ModuleData{
				ModuleName:      "order",
				CapitalizedName: "Order",
				ProjectModule:   "github.com/user/testapp",
			},
			wantContains: []string{
				"package order",
				"Order",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderer.Render(tt.templateName, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Render() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Render() unexpected error = %v", err)
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Render() result does not contain %q", want)
				}
			}
		})
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	_, err = renderer.Render("nonexistent.tmpl", ProjectData{})
	if err == nil {
		t.Error("Render() should fail with nonexistent template")
	}
}
