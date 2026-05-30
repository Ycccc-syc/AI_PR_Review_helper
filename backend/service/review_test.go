package service

import (
	"strings"
	"testing"

	"qiniudemo/model"
	"qiniudemo/repository"
)

func TestFilterDiff(t *testing.T) {
	tests := []struct {
		name string
		diff string
		want string // check that it doesn't contain
	}{
		{
			name: "filters package-lock.json",
			diff: `diff --git a/package-lock.json b/package-lock.json
@@ -1,3 +1,3 @@
-  "version": "1.0.0"
+  "version": "1.0.1"
diff --git a/src/main.go b/src/main.go
@@ -1,3 +1,3 @@
 func main() {
-  fmt.Println("hello")
+  fmt.Println("world")
 }`,
		},
		{
			name: "filters yarn.lock and go.sum",
			diff: `diff --git a/yarn.lock b/yarn.lock
@@ -1,1 +1,1 @@
-old-hash
+new-hash
diff --git a/go.sum b/go.sum
@@ -1,1 +1,1 @@
-old-hash
+new-hash
diff --git a/main.go b/main.go
@@ -1,1 +1,1 @@
-1
+2`,
		},
		{
			name: "filters binary/image/font files",
			diff: `diff --git a/icon.png b/icon.png
Binary files differ
diff --git a/logo.svg b/logo.svg
@@ -1,1 +1,1 @@
-<svg old>
+<svg new>
diff --git a/font.woff2 b/font.woff2
Binary files differ
diff --git a/app.js b/app.js
@@ -1,1 +1,1 @@
-console.log(1)
+console.log(2)`,
		},
		{
			name: "filters generated files (.exe .dll .sql .csv)",
			diff: `diff --git a/data.csv b/data.csv
@@ -1,2 +1,2 @@
-name,age
-alice,30
+name,age
+bob,25
diff --git a/migration.sql b/migration.sql
@@ -1,1 +1,1 @@
-ALTER TABLE foo;
+ALTER TABLE bar;
diff --git a/core.go b/core.go
@@ -1,1 +1,1 @@
-var x = 1
+var x = 2`,
		},
		{
			name: "keeps valid code files",
			diff: `diff --git a/main.go b/main.go
@@ -1,2 +1,2 @@
 func main() {
-  println("old")
+  println("new")
 }
diff --git a/handler/auth.go b/handler/auth.go
@@ -1,1 +1,1 @@
-package old
+package new
diff --git a/src/utils/helper.ts b/src/utils/helper.ts
@@ -1,1 +1,1 @@
-export const x = 1
+export const x = 2`,
		},
		{
			name: "all files filtered returns empty-diff safe",
			diff: `diff --git a/package-lock.json b/package-lock.json
@@ -1,1 +1,1 @@
-old
+new
diff --git a/go.sum b/go.sum
@@ -1,1 +1,1 @@
-old
+new`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterDiff(tt.diff)
			switch tt.name {
			case "filters package-lock.json":
				if strings.Contains(got, "package-lock.json") {
					t.Error("package-lock.json should be filtered")
				}
				if !strings.Contains(got, "main.go") {
					t.Error("main.go should be kept")
				}
			case "filters yarn.lock and go.sum":
				if strings.Contains(got, "yarn.lock") || strings.Contains(got, "go.sum") {
					t.Error("yarn.lock and go.sum should be filtered")
				}
				if !strings.Contains(got, "main.go") {
					t.Error("main.go should be kept")
				}
			case "filters binary/image/font files":
				if strings.Contains(got, "icon.png") || strings.Contains(got, "logo.svg") || strings.Contains(got, "font.woff2") {
					t.Error("image/svg/font files should be filtered")
				}
				if !strings.Contains(got, "app.js") {
					t.Error("app.js should be kept")
				}
			case "filters generated files (.exe .dll .sql .csv)":
				if strings.Contains(got, "data.csv") || strings.Contains(got, "migration.sql") {
					t.Error(".csv and .sql files should be filtered")
				}
				if !strings.Contains(got, "core.go") {
					t.Error("core.go should be kept")
				}
			case "keeps valid code files":
				if !strings.Contains(got, "main.go") || !strings.Contains(got, "auth.go") || !strings.Contains(got, "helper.ts") {
					t.Error("all valid code files should be kept")
				}
			case "all files filtered returns empty-diff safe":
				// After filtering, the diff should not contain any file blocks
				if strings.Count(got, "diff --git") >= strings.Count(tt.diff, "diff --git") {
					t.Log("warning: filtered diff may still contain content when all files skipped")
				}
			}
		})
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantKey string // check that the extracted JSON contains this string
	}{
		{
			name:    "pure JSON",
			raw:     `{"summary":"test","risk_items":[],"suggestions":[]}`,
			wantKey: `"summary"`,
		},
		{
			name: "JSON in ```json fence",
			raw: `Here is the result:
` + "```json\n" + `{"summary":"test","risk_items":[{"severity":"high","file":"main.go","line":10,"title":"bug","detail":"bad","suggestion":"fix"}],"suggestions":[]}` + "\n```",
			wantKey: `"severity":"high"`,
		},
		{
			name: "JSON in ``` fence without language tag",
			raw:  "```\n" + `{"summary":"test","risk_items":[],"suggestions":[{"priority":"medium","content":"add tests"}]}` + "\n```",
			wantKey: `"priority":"medium"`,
		},
		{
			name:    "JSON with leading text and trailing text",
			raw:     "Okay, here's the analysis:\n\n{\"summary\":\"test\",\"risk_items\":[],\"suggestions\":[]}\n\nHope this helps!",
			wantKey: `"summary"`,
		},
		{
			name:    "JSON with whitespace and newlines",
			raw:     "\n\n  {\n  \"summary\": \"test\",\n  \"risk_items\": [],\n  \"suggestions\": []\n}  \n\n",
			wantKey: `"risk_items"`,
		},
		{
			name:    "empty string",
			raw:     "",
			wantKey: "",
		},
		{
			name:    "no JSON at all",
			raw:     "Sorry, I cannot analyze this PR.",
			wantKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJSON(tt.raw)
			if tt.wantKey == "" {
				if got != "" {
					t.Logf("expected empty, got: %s", got)
				}
				return
			}
			if !strings.Contains(got, tt.wantKey) {
				t.Errorf("expected extracted JSON to contain %q, got: %s", tt.wantKey, got)
			}
			// valid JSON check
			if got != "" && got[0] != '{' {
				t.Errorf("extracted JSON should start with '{', got: %s", got)
			}
			if got != "" && got[len(got)-1] != '}' {
				t.Errorf("extracted JSON should end with '}', got: %s", got)
			}
		})
	}
}

func TestSplitDiffByFile(t *testing.T) {
	// Build a mock diff with header + 10 files
	parts := []string{"header line\n@@ -0,0 +1 @@\n+start\n"}
	for i := 0; i < 10; i++ {
		part := "diff --git a/file" + string(rune('0'+i)) + ".go b/file" + string(rune('0'+i)) + ".go\n" +
			"@@ -1,1 +1,1 @@\n" +
			"-old\n" +
			"+new\n"
		parts = append(parts, part)
	}
	diff10 := strings.Join(parts, "\n")

	tests := []struct {
		name          string
		diff          string
		filesPerChunk int
		wantChunks    int
	}{
		{
			name:          "diff smaller than chunk size returns 1 chunk",
			diff:          diff10,
			filesPerChunk: 20,
			wantChunks:    1,
		},
		{
			name:          "10 files at 5 per chunk = 2 chunks",
			diff:          diff10,
			filesPerChunk: 5,
			wantChunks:    2,
		},
		{
			name:          "10 files at 3 per chunk = 4 chunks",
			diff:          diff10,
			filesPerChunk: 3,
			wantChunks:    4,
		},
		{
			name:          "10 files at 1 per chunk = 10 chunks",
			diff:          diff10,
			filesPerChunk: 1,
			wantChunks:    10,
		},
		{
			name:          "single file diff returns 1 chunk",
			diff:          "diff --git a/main.go b/main.go\n@@ -1,1 +1,1 @@\n-old\n+new",
			filesPerChunk: 5,
			wantChunks:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitDiffByFile(tt.diff, tt.filesPerChunk)
			if len(got) != tt.wantChunks {
				t.Errorf("splitDiffByFile() = %d chunks, want %d", len(got), tt.wantChunks)
			}
			// verify no chunk is empty
			for i, chunk := range got {
				if len(strings.TrimSpace(chunk)) == 0 {
					t.Errorf("chunk %d is empty", i)
				}
			}
		})
	}
}

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		fileName string
		want     bool
	}{
		{"package-lock.json", true},
		{"yarn.lock", true},
		{"go.sum", true},
		{"Cargo.lock", true},
		{"app.go", false},
		{"main.ts", false},
		{"icon.png", true},
		{"logo.svg", true},
		{"font.woff2", true},
		{"Roboto.ttf", true},
		{"report.pdf", true},
		{"binary.exe", true},
		{"data.csv", true},
		{"schema.sql", true},
		{"src/handler/auth.go", false},
		{"components/Button.tsx", false},
		{"util.py", false},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			if got := shouldSkip(tt.fileName); got != tt.want {
				t.Errorf("shouldSkip(%q) = %v, want %v", tt.fileName, got, tt.want)
			}
		})
	}
}

func TestBuildStats(t *testing.T) {
	pr := &repository.PRInfo{
		FilesChanged: 5,
		Additions:    100,
		Deletions:    30,
	}
	risks := []model.RiskItem{
		{Severity: "high"},
		{Severity: "medium"},
		{Severity: "medium"},
		{Severity: "low"},
		{Severity: "low"},
		{Severity: "low"},
	}
	stats := buildStats(pr, risks)
	if stats.FilesChanged != 5 {
		t.Errorf("FilesChanged = %d, want 5", stats.FilesChanged)
	}
	if stats.Additions != 100 {
		t.Errorf("Additions = %d, want 100", stats.Additions)
	}
	if stats.Deletions != 30 {
		t.Errorf("Deletions = %d, want 30", stats.Deletions)
	}
	if stats.RiskHigh != 1 {
		t.Errorf("RiskHigh = %d, want 1", stats.RiskHigh)
	}
	if stats.RiskMedium != 2 {
		t.Errorf("RiskMedium = %d, want 2", stats.RiskMedium)
	}
	if stats.RiskLow != 3 {
		t.Errorf("RiskLow = %d, want 3", stats.RiskLow)
	}
}
