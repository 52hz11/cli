// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package wiki

import (
	"encoding/json"
	"testing"

	"github.com/larksuite/cli/internal/citation"
	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/internal/envvars"
	"github.com/larksuite/cli/internal/httpmock"
)

func TestWikiNodeCitation(t *testing.T) {
	got := wikiNodeCitation(core.BrandFeishu, "wikcnTok", "标题", "1721996760")
	if len(got) != 1 {
		t.Fatalf("wikiNodeCitation() = %#v, want 1 entry", got)
	}
	c := got[0]
	if c.SourceType != citation.SourceWiki {
		t.Errorf("source_type = %d", c.SourceType)
	}
	if c.URL != "https://applink.feishu.cn/client/wiki/open?wikiToken=wikcnTok" {
		t.Errorf("url = %q", c.URL)
	}
	if c.Title != "标题" || c.ResourceID != "wikcnTok" {
		t.Errorf("title/resource_id = %q %q", c.Title, c.ResourceID)
	}
	if c.PublishTime != citation.Time("1721996760") {
		t.Errorf("publish_time = %q", c.PublishTime)
	}
}

func TestWikiNodeCitationLarkBrand(t *testing.T) {
	got := wikiNodeCitation(core.BrandLark, "wikcnTok", "t", "")
	if got[0].URL != "https://applink.larksuite.com/client/wiki/open?wikiToken=wikcnTok" {
		t.Errorf("lark brand url = %q", got[0].URL)
	}
	if got[0].PublishTime != "" {
		t.Errorf("empty edit time must omit publish_time, got %q", got[0].PublishTime)
	}
}

func TestWikiNodeCitationEmptyToken(t *testing.T) {
	got := wikiNodeCitation(core.BrandFeishu, "", "t", "")
	if len(got) != 1 || got[0].URL != "" {
		t.Fatalf("empty token must yield empty url (Normalize drops it): %#v", got)
	}
}

func TestWikiNodeGetMountedExecuteEmitsCitation(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())
	t.Setenv(envvars.CliCitation, "1")

	factory, stdout, _, registry := cmdutil.TestFactory(t, wikiTestConfig())
	registry.Register(&httpmock.Stub{
		Method: "GET",
		URL:    "/open-apis/wiki/v2/spaces/get_node",
		Body: map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"node": map[string]interface{}{
					"space_id":      "space_123",
					"node_token":    "wikcnABC",
					"obj_token":     "docxXYZ",
					"obj_type":      "docx",
					"title":         "Design Spec",
					"obj_edit_time": "1700000000",
				},
			},
			"msg": "success",
		},
	})

	err := mountAndRunWiki(t, WikiNodeGet, []string{
		"+node-get",
		"--node-token", testWikiNodeToken,
		"--as", "bot",
	}, factory, stdout)
	if err != nil {
		t.Fatalf("mountAndRunWiki() error = %v", err)
	}

	var envelope struct {
		OK        bool                   `json:"ok"`
		Data      map[string]interface{} `json:"data"`
		Citations []citation.Citation    `json:"citations"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal wiki envelope: %v\nstdout=%s", err, stdout.String())
	}
	if !envelope.OK {
		t.Fatalf("expected ok=true envelope, got stdout=%s", stdout.String())
	}
	if _, ok := envelope.Data["url"]; ok {
		t.Fatalf("citation must not add url to data: %#v", envelope.Data["url"])
	}
	if len(envelope.Citations) != 1 {
		t.Fatalf("citations = %#v, want 1 entry", envelope.Citations)
	}

	got := envelope.Citations[0]
	if got.SourceType != citation.SourceWiki {
		t.Errorf("source_type = %d, want %d", got.SourceType, citation.SourceWiki)
	}
	if got.URL != "https://applink.feishu.cn/client/wiki/open?wikiToken=wikcnABC" {
		t.Errorf("url = %q", got.URL)
	}
	if got.Title != "Design Spec" {
		t.Errorf("title = %q", got.Title)
	}
	if got.ResourceID != "wikcnABC" {
		t.Errorf("resource_id = %q, want node_token alone", got.ResourceID)
	}
	if got.PublishTime != citation.Time("1700000000") {
		t.Errorf("publish_time = %q", got.PublishTime)
	}
}
