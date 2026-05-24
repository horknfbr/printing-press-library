package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildPersonalizationAuditReportsDocumentedFields(t *testing.T) {
	product := sampleProduct()

	rows := buildPersonalizationAudit(product)

	if len(rows) != 1 {
		t.Fatalf("expected 1 audit row, got %d", len(rows))
	}
	if !rows[0].Supported {
		t.Fatalf("expected documented image/text fields to be supported")
	}
	if rows[0].InputText != "Name" || rows[0].FontFamily != "Inter" {
		t.Fatalf("unexpected personalization row: %#v", rows[0])
	}
}

func TestBuildPlacementMatrixJoinsUploadNames(t *testing.T) {
	product := sampleProduct()
	uploads := []ppJSONObj{{"id": "img_1", "file_name": "front.png"}}

	rows := buildPlacementMatrix(product, uploads)

	if len(rows) != 2 {
		t.Fatalf("expected one row per variant, got %d", len(rows))
	}
	if rows[0].UploadFileName != "front.png" || rows[0].Scale != 1.2 {
		t.Fatalf("unexpected placement row: %#v", rows[0])
	}
}

func TestBuildProductDriftDetectsChangedTitle(t *testing.T) {
	expected := ppJSONObj{"title": "Original", "blueprint_id": float64(384)}
	actual := ppJSONObj{"title": "Changed", "blueprint_id": float64(384)}

	rows := buildProductDrift(expected, actual)

	foundTitleDrift := false
	for _, row := range rows {
		if row.Path == "title" && row.Status == "drift" {
			foundTitleDrift = true
		}
	}
	if !foundTitleDrift {
		t.Fatalf("expected title drift in %#v", rows)
	}
}

func TestBuildPersonalizationBatchWritesExpandedManifest(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "template.json")
	csvPath := filepath.Join(dir, "rows.csv")
	outDir := filepath.Join(dir, "out")
	if err := os.WriteFile(templatePath, []byte(`{"title":"{{title}}","print_areas":[{"placeholders":[{"images":[{"id":"{{image_id}}","input_text":"{{text}}"}]}]}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(csvPath, []byte("title,image_id,text\nMug,img_1,Sam\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	rows, err := buildPersonalizationBatch(templatePath, csvPath, outDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Title != "Mug" || !rows[0].TextUsed {
		t.Fatalf("unexpected batch rows: %#v", rows)
	}
	data, err := os.ReadFile(rows[0].Output)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest["title"] != "Mug" {
		t.Fatalf("template token was not replaced: %#v", manifest)
	}
}

func TestBuildCatalogMarginMatrixComputesMargin(t *testing.T) {
	variants := []ppJSONObj{{"id": "1", "title": "S", "cost": float64(1200)}}
	shipping := []ppJSONObj{{"first_item": float64(500)}}

	rows := buildCatalogMarginMatrix(variants, shipping, 24.99)

	if len(rows) != 1 {
		t.Fatalf("expected one margin row, got %d", len(rows))
	}
	if rows[0].Cost != 12 || rows[0].Shipping != 5 || rows[0].EstimatedMargin != 7.99 {
		t.Fatalf("unexpected margin row: %#v", rows[0])
	}
}

func TestBuildAssetReuseFlagsUnusedAndSharedUploads(t *testing.T) {
	products := []ppJSONObj{sampleProduct(), sampleProductWithID("prod_2")}
	uploads := []ppJSONObj{{"id": "img_1", "file_name": "front.png"}, {"id": "unused", "file_name": "unused.png"}}

	rows := buildAssetReuse(products, uploads)

	if len(rows) != 2 {
		t.Fatalf("expected two upload rows, got %d", len(rows))
	}
	if !rows[0].SharedArtwork || rows[0].UseCount != 2 {
		t.Fatalf("expected shared artwork row, got %#v", rows[0])
	}
	if !rows[1].Unused {
		t.Fatalf("expected unused upload row, got %#v", rows[1])
	}
}

func TestBuildFulfillmentRiskFlagsMissingShipment(t *testing.T) {
	orders := []ppJSONObj{{
		"id":     "order_1",
		"status": "pending",
		"line_items": []any{
			map[string]any{"product_id": "prod_1", "variant_id": "101"},
		},
	}}
	products := []ppJSONObj{{"id": "prod_1", "status": "visible"}}

	rows := buildFulfillmentRisk(orders, products)

	if len(rows) != 1 {
		t.Fatalf("expected one risk row, got %d", len(rows))
	}
	if rows[0].Risks[0] != "no shipment records" {
		t.Fatalf("unexpected risks: %#v", rows[0].Risks)
	}
}

func sampleProduct() ppJSONObj {
	return sampleProductWithID("prod_1")
}

func sampleProductWithID(id string) ppJSONObj {
	return ppJSONObj{
		"id": id,
		"print_areas": []any{
			map[string]any{
				"variant_ids": []any{"101", "102"},
				"placeholders": []any{
					map[string]any{
						"position": "front",
						"images": []any{
							map[string]any{
								"id":          "img_1",
								"name":        "front art",
								"type":        "text",
								"input_text":  "Name",
								"font_family": "Inter",
								"x":           float64(0.5),
								"y":           float64(0.5),
								"scale":       float64(1.2),
							},
						},
					},
				},
			},
		},
	}
}
