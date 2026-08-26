// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/httpmock"
	"github.com/larksuite/cli/shortcuts/common"
)

func TestFormConfigGetCallsResourceEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		command string
		url     string
		run     func(*testing.T, []string, *cmdutil.Factory, *bytes.Buffer) error
	}{
		{
			name:    "submission settings",
			command: "+form-submission-settings-get",
			url:     "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/submission-settings",
			run: func(t *testing.T, args []string, factory *cmdutil.Factory, stdout *bytes.Buffer) error {
				return runShortcut(t, BaseFormSubmissionSettingsGet, args, factory, stdout)
			},
		},
		{
			name:    "notifications",
			command: "+form-notifications-get",
			url:     "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/notifications",
			run: func(t *testing.T, args []string, factory *cmdutil.Factory, stdout *bytes.Buffer) error {
				return runShortcut(t, BaseFormNotificationsGet, args, factory, stdout)
			},
		},
		{
			name:    "submit actions",
			command: "+form-submit-actions-get",
			url:     "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/submit-actions",
			run: func(t *testing.T, args []string, factory *cmdutil.Factory, stdout *bytes.Buffer) error {
				return runShortcut(t, BaseFormSubmitActionsGet, args, factory, stdout)
			},
		},
		{
			name:    "lottery",
			command: "+form-lottery-get",
			url:     "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/lottery",
			run: func(t *testing.T, args []string, factory *cmdutil.Factory, stdout *bytes.Buffer) error {
				return runShortcut(t, BaseFormLotteryGet, args, factory, stdout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, stdout, reg := newExecuteFactory(t)
			reg.Register(&httpmock.Stub{
				Method: "GET",
				URL:    tt.url,
				Body:   map[string]interface{}{"code": 0, "data": map[string]interface{}{"ok": true}},
			})
			err := tt.run(t, []string{
				tt.command,
				"--base-token", "app_x",
				"--table-id", "tbl_1",
				"--form-id", "vew_1",
			}, factory, stdout)
			if err != nil {
				t.Fatalf("run shortcut: %v", err)
			}
			if got := stdout.String(); !strings.Contains(got, `"ok": true`) {
				t.Fatalf("stdout=%s", got)
			}
		})
	}
}

func TestFormSubmissionSettingsUpdateBuildsBodies(t *testing.T) {
	tests := []struct {
		name  string
		flags []string
		want  map[string]interface{}
	}{
		{
			name: "submission period enabled",
			flags: []string{
				"--submission-period-enabled=true",
				"--start-at", "2026-08-25T10:00:00+08:00",
				"--end-at", "2026-08-31T18:00:00+08:00",
				"--timezone", "Asia/Shanghai",
			},
			want: map[string]interface{}{"submission_period": map[string]interface{}{
				"enabled": true, "start_at": "2026-08-25T10:00:00+08:00", "end_at": "2026-08-31T18:00:00+08:00", "timezone": "Asia/Shanghai",
			}},
		},
		{
			name:  "user submit limit",
			flags: []string{"--user-submit-limit-enabled=true", "--user-submit-limit", "3", "--user-submit-cycle", "day"},
			want:  map[string]interface{}{"user_submit_limit": map[string]interface{}{"enabled": true, "frequency_limit": float64(3), "frequency_cycle": "day"}},
		},
		{
			name:  "allow modify false",
			flags: []string{"--allow-modify-submission=false"},
			want:  map[string]interface{}{"allow_modify_submission": false},
		},
		{
			name:  "ai voice input false",
			flags: []string{"--ai-voice-input-enabled=false"},
			want:  map[string]interface{}{"ai_voice_input": map[string]interface{}{"enabled": false}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := runFormConfigPatch(t, BaseFormSubmissionSettingsUpdate, "+form-submission-settings-update", "submission-settings", tt.flags...)
			if got := decodeCapturedJSONBody(t, stub); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("request body=%#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFormNotificationsUpdateBuildsScheduledBody(t *testing.T) {
	stub := runFormConfigPatch(t, BaseFormNotificationsUpdate, "+form-notifications-update", "notifications",
		"--type", "scheduled",
		"--enabled=true",
		"--locale", "zh_cn",
		"--receivers-json", `[{"open_id":"ou_1"}]`,
		"--notify-time", "2026-08-25T10:00:00+08:00",
		"--repeat-type", "day",
		"--timezone", "Asia/Shanghai",
	)
	want := map[string]interface{}{
		"locale": "zh_cn",
		"scheduled": map[string]interface{}{
			"enabled":     true,
			"receivers":   []interface{}{map[string]interface{}{"open_id": "ou_1"}},
			"notify_time": "2026-08-25T10:00:00+08:00",
			"repeat_type": "day",
			"timezone":    "Asia/Shanghai",
		},
	}
	if got := decodeCapturedJSONBody(t, stub); !reflect.DeepEqual(got, want) {
		t.Fatalf("request body=%#v, want %#v", got, want)
	}
}

func TestFormSubmitActionsUpdateBuildsResultPageBody(t *testing.T) {
	stub := runFormConfigPatch(t, BaseFormSubmitActionsUpdate, "+form-submit-actions-update", "submit-actions",
		"--type", "result-page",
		"--enabled=true",
		"--revision", "123",
		"--title", "submitted",
		"--description-json", `[{"type":"text","text":"thanks"},{"type":"url","text":"detail","url":"https://example.com"},{"type":"mention","open_id":"ou_1"}]`,
	)
	want := map[string]interface{}{
		"revision": float64(123),
		"result_page": map[string]interface{}{
			"enabled": true,
			"title":   "submitted",
			"description": []interface{}{
				map[string]interface{}{"type": "text", "text": "thanks"},
				map[string]interface{}{"type": "url", "text": "detail", "url": "https://example.com"},
				map[string]interface{}{"type": "mention", "open_id": "ou_1"},
			},
		},
	}
	if got := decodeCapturedJSONBody(t, stub); !reflect.DeepEqual(got, want) {
		t.Fatalf("request body=%#v, want %#v", got, want)
	}
}

func TestFormLotteryActionBuildsUpdateBody(t *testing.T) {
	stub := runFormConfigPost(t, BaseFormLotteryAction, "+form-lottery-action", "lottery/actions",
		"--action", "update",
		"--lottery-json", `{"version":1,"probability":10000,"awards":[{"name":"gift","quantity":10}]}`,
	)
	want := map[string]interface{}{
		"action": "update",
		"lottery": map[string]interface{}{
			"version":     float64(1),
			"probability": float64(10000),
			"awards":      []interface{}{map[string]interface{}{"name": "gift", "quantity": float64(10)}},
		},
	}
	if got := decodeCapturedJSONBody(t, stub); !reflect.DeepEqual(got, want) {
		t.Fatalf("request body=%#v, want %#v", got, want)
	}
}

func TestFormConfigRejectsUnsupportedWriteFields(t *testing.T) {
	tests := []struct {
		name     string
		shortcut common.Shortcut
		command  string
		args     []string
		param    string
		want     string
	}{
		{
			name:     "notification open chat",
			shortcut: BaseFormNotificationsUpdate,
			command:  "+form-notifications-update",
			args:     []string{"--type", "scheduled", "--enabled=true", "--receivers-json", `[{"open_chat_id":"oc_1"}]`, "--notify-time", "2026-08-25T10:00:00+08:00", "--repeat-type", "day", "--timezone", "Asia/Shanghai"},
			param:    "--receivers-json",
			want:     "open_chat_id",
		},
		{
			name:     "lottery icon token",
			shortcut: BaseFormLotteryAction,
			command:  "+form-lottery-action",
			args:     []string{"--action", "enable", "--lottery-json", `{"icon_token":"img_1","awards":[]}`},
			param:    "--lottery-json",
			want:     "icon_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, stdout, _ := newExecuteFactory(t)
			args := []string{tt.command, "--base-token", "app_x", "--table-id", "tbl_1", "--form-id", "vew_1"}
			args = append(args, tt.args...)
			err := runShortcut(t, tt.shortcut, args, factory, stdout)
			assertInvalidArgumentValidation(t, err, tt.param, nil, tt.want)
		})
	}
}

func runFormConfigPatch(t *testing.T, shortcut common.Shortcut, command, segment string, flags ...string) *httpmock.Stub {
	t.Helper()
	factory, stdout, reg := newExecuteFactory(t)
	stub := &httpmock.Stub{
		Method: "PATCH",
		URL:    "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/" + segment,
		Body:   map[string]interface{}{"code": 0, "data": map[string]interface{}{"ok": true}},
	}
	reg.Register(stub)
	args := []string{command, "--base-token", "app_x", "--table-id", "tbl_1", "--form-id", "vew_1"}
	args = append(args, flags...)
	if err := runShortcut(t, shortcut, args, factory, stdout); err != nil {
		t.Fatalf("run shortcut: %v", err)
	}
	return stub
}

func runFormConfigPost(t *testing.T, shortcut common.Shortcut, command, segment string, flags ...string) *httpmock.Stub {
	t.Helper()
	factory, stdout, reg := newExecuteFactory(t)
	stub := &httpmock.Stub{
		Method: "POST",
		URL:    "/open-apis/base/v3/bases/app_x/tables/tbl_1/forms/vew_1/" + segment,
		Body:   map[string]interface{}{"code": 0, "data": map[string]interface{}{"ok": true}},
	}
	reg.Register(stub)
	args := []string{command, "--base-token", "app_x", "--table-id", "tbl_1", "--form-id", "vew_1"}
	args = append(args, flags...)
	if err := runShortcut(t, shortcut, args, factory, stdout); err != nil {
		t.Fatalf("run shortcut: %v", err)
	}
	return stub
}
