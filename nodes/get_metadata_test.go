package nodes_test

import (
	"context"
	"reflect"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetMetadata_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetMetadata(ctx, ax, &gen.GetMetadataInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantDesc := "A high performance web server and reverse proxy server"
	if got.Description != wantDesc {
		t.Errorf("Description = %q, want %q", got.Description, wantDesc)
	}
	wantDocs := []string{"man:nginx(8)", "https://nginx.org/en/docs/"}
	if !reflect.DeepEqual(got.Documentation, wantDocs) {
		t.Errorf("Documentation = %v, want %v", got.Documentation, wantDocs)
	}
}

func TestGetMetadata_Absent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetMetadata(ctx, ax, &gen.GetMetadataInput{Text: "[Unit]\nAfter=network.target\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != "" || len(got.Documentation) != 0 {
		t.Errorf("expected empty metadata, got %+v", got)
	}
}
