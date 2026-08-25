// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/larksuite/cli/shortcuts/common"
)

var formConfigCommonFlags = []common.Flag{
	baseTokenFlag(true),
	{Name: "table-id", Desc: "table ID", Required: true},
	{Name: "form-id", Desc: "form ID", Required: true},
}

var BaseFormSubmissionSettingsGet = formConfigGetShortcut(
	"+form-submission-settings-get",
	"Get form submission settings",
	"submission-settings",
)

var BaseFormSubmissionSettingsUpdate = common.Shortcut{
	Service:     "base",
	Command:     "+form-submission-settings-update",
	Description: "Update one form submission setting group",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   authTypes(),
	Flags: appendFormConfigFlags(
		common.Flag{Name: "submission-period-enabled", Type: "bool", Desc: "enable or disable the form submission period"},
		common.Flag{Name: "start-at", Desc: "submission start time in RFC3339 format"},
		common.Flag{Name: "end-at", Desc: "submission end time in RFC3339 format"},
		common.Flag{Name: "timezone", Desc: "IANA timezone, for example Asia/Shanghai"},
		common.Flag{Name: "user-submit-limit-enabled", Type: "bool", Desc: "enable or disable per-user submission limit"},
		common.Flag{Name: "user-submit-limit", Type: "int", Desc: "maximum submissions per user"},
		common.Flag{Name: "user-submit-cycle", Desc: "per-user limit cycle", Enum: []string{"total", "day", "week", "month"}},
		common.Flag{Name: "total-submit-limit-enabled", Type: "bool", Desc: "enable or disable total submission limit"},
		common.Flag{Name: "total-submit-maximum", Type: "int", Desc: "maximum total submissions"},
		common.Flag{Name: "allow-modify-submission", Type: "bool", Desc: "allow users to modify submitted records"},
		common.Flag{Name: "ai-voice-input-enabled", Type: "bool", Desc: "enable or disable AI voice input"},
	),
	Tips: []string{
		"Update exactly one top-level group per invocation: submission_period, user_submit_limit, total_submit_limit, allow_modify_submission, or ai_voice_input.",
		"Boolean flags support explicit false, for example --ai-voice-input-enabled=false.",
	},
	Validate: func(_ context.Context, runtime *common.RuntimeContext) error {
		_, err := buildFormSubmissionSettingsBody(runtime)
		return err
	},
	DryRun:  formConfigPatchDryRun("submission-settings", buildFormSubmissionSettingsBody),
	Execute: formConfigPatchExecute("submission-settings", buildFormSubmissionSettingsBody),
}

var BaseFormNotificationsGet = formConfigGetShortcut(
	"+form-notifications-get",
	"Get form notification settings",
	"notifications",
)

var BaseFormNotificationsUpdate = common.Shortcut{
	Service:     "base",
	Command:     "+form-notifications-update",
	Description: "Update form notification settings",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   authTypes(),
	Flags: appendFormConfigFlags(
		common.Flag{Name: "type", Desc: "notification type", Required: true, Enum: []string{"on-submission", "scheduled"}},
		common.Flag{Name: "enabled", Type: "bool", Desc: "enable or disable this notification", Required: true},
		common.Flag{Name: "locale", Desc: "notification locale"},
		common.Flag{Name: "receivers-json", Desc: "receiver list JSON array; open_chat_id is not supported", Input: []string{common.File, common.Stdin}},
		common.Flag{Name: "notify-time", Desc: "scheduled notify time in RFC3339 format"},
		common.Flag{Name: "repeat-type", Desc: "scheduled repeat type"},
		common.Flag{Name: "timezone", Desc: "IANA timezone, for example Asia/Shanghai"},
	),
	Tips: []string{
		"Use --type on-submission or --type scheduled; update only one notification group per invocation.",
		"Disabling a notification only accepts --type and --enabled=false, plus optional --locale.",
	},
	Validate: func(_ context.Context, runtime *common.RuntimeContext) error {
		_, err := buildFormNotificationsBody(runtime)
		return err
	},
	DryRun:  formConfigPatchDryRun("notifications", buildFormNotificationsBody),
	Execute: formConfigPatchExecute("notifications", buildFormNotificationsBody),
}

var BaseFormSubmitActionsGet = formConfigGetShortcut(
	"+form-submit-actions-get",
	"Get form post-submit actions",
	"submit-actions",
)

var BaseFormSubmitActionsUpdate = common.Shortcut{
	Service:     "base",
	Command:     "+form-submit-actions-update",
	Description: "Update one form post-submit action",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   authTypes(),
	Flags: appendFormConfigFlags(
		common.Flag{Name: "type", Desc: "submit action type", Required: true, Enum: []string{"result-page", "redirect"}},
		common.Flag{Name: "enabled", Type: "bool", Desc: "enable or disable the action", Required: true},
		common.Flag{Name: "revision", Type: "int", Desc: "current action revision"},
		common.Flag{Name: "title", Desc: "result page title"},
		common.Flag{Name: "description-json", Desc: "result page description JSON array", Input: []string{common.File, common.Stdin}},
		common.Flag{Name: "redirect-url", Desc: "redirect URL after form submit"},
	),
	Tips: []string{
		"Use result-page for the submit result page or redirect for a submit redirect; do not combine fields for both action types.",
		"--description-json supports only text, url, and mention blocks; mention blocks must use open_id.",
	},
	Validate: func(_ context.Context, runtime *common.RuntimeContext) error {
		_, err := buildFormSubmitActionsBody(runtime)
		return err
	},
	DryRun:  formConfigPatchDryRun("submit-actions", buildFormSubmitActionsBody),
	Execute: formConfigPatchExecute("submit-actions", buildFormSubmitActionsBody),
}

var BaseFormLotteryGet = formConfigGetShortcut(
	"+form-lottery-get",
	"Get form lottery settings",
	"lottery",
)

var BaseFormLotteryAction = common.Shortcut{
	Service:     "base",
	Command:     "+form-lottery-action",
	Description: "Run a form lottery action",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   authTypes(),
	Flags: appendFormConfigFlags(
		common.Flag{Name: "action", Desc: "lottery action", Required: true, Enum: []string{"enable", "disable", "update", "relink_winning_table"}},
		common.Flag{Name: "lottery-json", Desc: "lottery config JSON object; icon_token is not supported", Input: []string{common.File, common.Stdin}},
	),
	Tips: []string{
		"enable and update require --lottery-json with full lottery settings; update also requires lottery.version.",
		"Do not pass icon_token in lottery JSON.",
	},
	Validate: func(_ context.Context, runtime *common.RuntimeContext) error {
		_, err := buildFormLotteryActionBody(runtime)
		return err
	},
	DryRun: func(_ context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		body, _ := buildFormLotteryActionBody(runtime)
		return common.NewDryRunAPI().
			POST(formConfigDryRunPath("lottery/actions")).
			Body(body).
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id")).
			Set("form_id", runtime.Str("form-id"))
	},
	Execute: func(_ context.Context, runtime *common.RuntimeContext) error {
		body, err := buildFormLotteryActionBody(runtime)
		if err != nil {
			return err
		}
		data, err := baseV3Call(runtime, "POST", formConfigPath(runtime, "lottery", "actions"), nil, body)
		if err != nil {
			return err
		}
		runtime.Out(data, nil)
		return nil
	},
}

func formConfigGetShortcut(command, description, segment string) common.Shortcut {
	return common.Shortcut{
		Service:     "base",
		Command:     command,
		Description: description,
		Risk:        "read",
		Scopes:      []string{"base:form:update"},
		AuthTypes:   authTypes(),
		Flags:       formConfigCommonFlags,
		DryRun: func(_ context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
			return common.NewDryRunAPI().
				GET(formConfigDryRunPath(segment)).
				Set("base_token", runtime.Str("base-token")).
				Set("table_id", runtime.Str("table-id")).
				Set("form_id", runtime.Str("form-id"))
		},
		Execute: func(_ context.Context, runtime *common.RuntimeContext) error {
			data, err := baseV3Call(runtime, "GET", formConfigPath(runtime, segment), nil, nil)
			if err != nil {
				return err
			}
			runtime.Out(data, nil)
			return nil
		},
	}
}

func appendFormConfigFlags(flags ...common.Flag) []common.Flag {
	out := make([]common.Flag, 0, len(formConfigCommonFlags)+len(flags))
	out = append(out, formConfigCommonFlags...)
	out = append(out, flags...)
	return out
}

func formConfigDryRunPath(segment string) string {
	return "/open-apis/base/v3/bases/:base_token/tables/:table_id/forms/:form_id/" + segment
}

func formConfigPath(runtime *common.RuntimeContext, segments ...string) string {
	parts := []string{"bases", runtime.Str("base-token"), "tables", runtime.Str("table-id"), "forms", runtime.Str("form-id")}
	parts = append(parts, segments...)
	return baseV3Path(parts...)
}

func formConfigPatchDryRun(segment string, build func(*common.RuntimeContext) (map[string]interface{}, error)) func(context.Context, *common.RuntimeContext) *common.DryRunAPI {
	return func(_ context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		body, _ := build(runtime)
		return common.NewDryRunAPI().
			PATCH(formConfigDryRunPath(segment)).
			Body(body).
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id")).
			Set("form_id", runtime.Str("form-id"))
	}
}

func formConfigPatchExecute(segment string, build func(*common.RuntimeContext) (map[string]interface{}, error)) func(context.Context, *common.RuntimeContext) error {
	return func(_ context.Context, runtime *common.RuntimeContext) error {
		body, err := build(runtime)
		if err != nil {
			return err
		}
		data, err := baseV3Call(runtime, "PATCH", formConfigPath(runtime, segment), nil, body)
		if err != nil {
			return err
		}
		runtime.Out(data, nil)
		return nil
	}
}

func buildFormSubmissionSettingsBody(runtime *common.RuntimeContext) (map[string]interface{}, error) {
	groups := changedFormSubmissionGroups(runtime)
	if len(groups) != 1 {
		return nil, baseFlagErrorf("update exactly one form submission settings group; changed groups: %s", strings.Join(groups, ", "))
	}

	body := map[string]interface{}{}
	switch groups[0] {
	case "submission_period":
		if !runtime.Changed("submission-period-enabled") {
			return nil, baseFlagErrorf("--submission-period-enabled is required when updating submission period")
		}
		enabled := runtime.Bool("submission-period-enabled")
		period := map[string]interface{}{"enabled": enabled}
		if enabled {
			if runtime.Str("start-at") == "" {
				return nil, baseFlagErrorf("--start-at is required when --submission-period-enabled=true")
			}
			if runtime.Str("end-at") == "" {
				return nil, baseFlagErrorf("--end-at is required when --submission-period-enabled=true")
			}
			if runtime.Str("timezone") == "" {
				return nil, baseFlagErrorf("--timezone is required when --submission-period-enabled=true")
			}
			if err := validateRFC3339Flag("start-at", runtime.Str("start-at")); err != nil {
				return nil, err
			}
			if err := validateRFC3339Flag("end-at", runtime.Str("end-at")); err != nil {
				return nil, err
			}
			period["start_at"] = runtime.Str("start-at")
			period["end_at"] = runtime.Str("end-at")
			period["timezone"] = runtime.Str("timezone")
		} else if runtime.Changed("timezone") {
			period["timezone"] = runtime.Str("timezone")
		}
		body["submission_period"] = period
	case "user_submit_limit":
		if !runtime.Changed("user-submit-limit-enabled") {
			return nil, baseFlagErrorf("--user-submit-limit-enabled is required when updating user submit limit")
		}
		enabled := runtime.Bool("user-submit-limit-enabled")
		limit := map[string]interface{}{"enabled": enabled}
		if enabled {
			if runtime.Int("user-submit-limit") <= 0 {
				return nil, baseFlagErrorf("--user-submit-limit must be greater than 0 when --user-submit-limit-enabled=true")
			}
			if runtime.Str("user-submit-cycle") == "" {
				return nil, baseFlagErrorf("--user-submit-cycle is required when --user-submit-limit-enabled=true")
			}
			limit["frequency_limit"] = runtime.Int("user-submit-limit")
			limit["frequency_cycle"] = runtime.Str("user-submit-cycle")
		}
		body["user_submit_limit"] = limit
	case "total_submit_limit":
		if !runtime.Changed("total-submit-limit-enabled") {
			return nil, baseFlagErrorf("--total-submit-limit-enabled is required when updating total submit limit")
		}
		enabled := runtime.Bool("total-submit-limit-enabled")
		limit := map[string]interface{}{"enabled": enabled}
		if enabled {
			if runtime.Int("total-submit-maximum") <= 0 {
				return nil, baseFlagErrorf("--total-submit-maximum must be greater than 0 when --total-submit-limit-enabled=true")
			}
			limit["maximum"] = runtime.Int("total-submit-maximum")
		}
		body["total_submit_limit"] = limit
	case "allow_modify_submission":
		body["allow_modify_submission"] = runtime.Bool("allow-modify-submission")
	case "ai_voice_input":
		body["ai_voice_input"] = map[string]interface{}{"enabled": runtime.Bool("ai-voice-input-enabled")}
	}
	return body, nil
}

func changedFormSubmissionGroups(runtime *common.RuntimeContext) []string {
	groups := []string{}
	if anyChanged(runtime, "submission-period-enabled", "start-at", "end-at", "timezone") {
		groups = append(groups, "submission_period")
	}
	if anyChanged(runtime, "user-submit-limit-enabled", "user-submit-limit", "user-submit-cycle") {
		groups = append(groups, "user_submit_limit")
	}
	if anyChanged(runtime, "total-submit-limit-enabled", "total-submit-maximum") {
		groups = append(groups, "total_submit_limit")
	}
	if runtime.Changed("allow-modify-submission") {
		groups = append(groups, "allow_modify_submission")
	}
	if runtime.Changed("ai-voice-input-enabled") {
		groups = append(groups, "ai_voice_input")
	}
	return groups
}

func buildFormNotificationsBody(runtime *common.RuntimeContext) (map[string]interface{}, error) {
	enabled := runtime.Bool("enabled")
	body := map[string]interface{}{}
	if runtime.Changed("locale") {
		body["locale"] = runtime.Str("locale")
	}

	group := map[string]interface{}{"enabled": enabled}
	if runtime.Str("type") == "scheduled" {
		if enabled {
			if runtime.Str("receivers-json") == "" {
				return nil, baseFlagErrorf("--receivers-json is required when scheduled notification is enabled")
			}
			if runtime.Str("notify-time") == "" {
				return nil, baseFlagErrorf("--notify-time is required when scheduled notification is enabled")
			}
			if runtime.Str("repeat-type") == "" {
				return nil, baseFlagErrorf("--repeat-type is required when scheduled notification is enabled")
			}
			if runtime.Str("timezone") == "" {
				return nil, baseFlagErrorf("--timezone is required when scheduled notification is enabled")
			}
			receivers, err := parseJSONArrayFlag("receivers-json", runtime.Str("receivers-json"))
			if err != nil {
				return nil, err
			}
			if containsKey(receivers, "open_chat_id") {
				return nil, baseFlagErrorf("--receivers-json must not contain open_chat_id")
			}
			if err := validateRFC3339Flag("notify-time", runtime.Str("notify-time")); err != nil {
				return nil, err
			}
			group["receivers"] = receivers
			group["notify_time"] = runtime.Str("notify-time")
			group["repeat_type"] = runtime.Str("repeat-type")
			group["timezone"] = runtime.Str("timezone")
		} else if anyChanged(runtime, "receivers-json", "notify-time", "repeat-type", "timezone") {
			return nil, baseFlagErrorf("disabling scheduled notification only accepts --type, --enabled=false, and optional --locale")
		}
		body["scheduled"] = group
		return body, nil
	}

	if !enabled && anyChanged(runtime, "receivers-json", "notify-time", "repeat-type", "timezone") {
		return nil, baseFlagErrorf("disabling on-submission notification only accepts --type, --enabled=false, and optional --locale")
	}
	if runtime.Changed("receivers-json") {
		receivers, err := parseJSONArrayFlag("receivers-json", runtime.Str("receivers-json"))
		if err != nil {
			return nil, err
		}
		if containsKey(receivers, "open_chat_id") {
			return nil, baseFlagErrorf("--receivers-json must not contain open_chat_id")
		}
		group["receivers"] = receivers
	}
	body["on_submission"] = group
	return body, nil
}

func buildFormSubmitActionsBody(runtime *common.RuntimeContext) (map[string]interface{}, error) {
	enabled := runtime.Bool("enabled")
	body := map[string]interface{}{}
	if runtime.Changed("revision") {
		body["revision"] = runtime.Int("revision")
	}
	group := map[string]interface{}{"enabled": enabled}

	switch runtime.Str("type") {
	case "result-page":
		if runtime.Changed("redirect-url") {
			return nil, baseFlagErrorf("--redirect-url cannot be used with --type result-page")
		}
		if enabled {
			if runtime.Str("title") == "" {
				return nil, baseFlagErrorf("--title is required when result page is enabled")
			}
			if runtime.Str("description-json") == "" {
				return nil, baseFlagErrorf("--description-json is required when result page is enabled")
			}
			description, err := parseJSONArrayFlag("description-json", runtime.Str("description-json"))
			if err != nil {
				return nil, err
			}
			if err := validateResultPageDescription(description); err != nil {
				return nil, err
			}
			group["title"] = runtime.Str("title")
			group["description"] = description
		}
		body["result_page"] = group
		return body, nil
	case "redirect":
		if runtime.Changed("title") || runtime.Changed("description-json") {
			return nil, baseFlagErrorf("--title and --description-json cannot be used with --type redirect")
		}
		if enabled {
			if runtime.Str("redirect-url") == "" {
				return nil, baseFlagErrorf("--redirect-url is required when redirect is enabled")
			}
			group["url"] = runtime.Str("redirect-url")
		}
		body["redirect"] = group
		return body, nil
	default:
		return nil, baseFlagErrorf("--type must be result-page or redirect")
	}
}

func buildFormLotteryActionBody(runtime *common.RuntimeContext) (map[string]interface{}, error) {
	action := runtime.Str("action")
	body := map[string]interface{}{"action": action}
	if action == "enable" || action == "update" {
		if runtime.Str("lottery-json") == "" {
			return nil, baseFlagErrorf("--lottery-json is required for %s action", action)
		}
		lottery, err := parseJSONObjectFlag("lottery-json", runtime.Str("lottery-json"))
		if err != nil {
			return nil, err
		}
		if containsKey(lottery, "icon_token") {
			return nil, baseFlagErrorf("--lottery-json must not contain icon_token")
		}
		if action == "update" {
			if _, ok := lottery["version"]; !ok {
				return nil, baseFlagErrorf("--lottery-json.version is required for update action")
			}
		}
		body["lottery"] = lottery
		return body, nil
	}
	if runtime.Changed("lottery-json") {
		return nil, baseFlagErrorf("--lottery-json is only accepted for enable or update action")
	}
	return body, nil
}

func validateRFC3339Flag(name, value string) error {
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return baseFlagErrorf("--%s must be RFC3339: %v", name, err)
	}
	return nil
}

func parseJSONObjectFlag(name, value string) (map[string]interface{}, error) {
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(value), &body); err != nil {
		return nil, baseFlagErrorf("--%s must be valid JSON object: %v", name, err)
	}
	if body == nil {
		return nil, baseFlagErrorf("--%s must be valid JSON object", name)
	}
	return body, nil
}

func parseJSONArrayFlag(name, value string) ([]interface{}, error) {
	var body []interface{}
	if err := json.Unmarshal([]byte(value), &body); err != nil {
		return nil, baseFlagErrorf("--%s must be valid JSON array: %v", name, err)
	}
	return body, nil
}

func validateResultPageDescription(items []interface{}) error {
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return baseFlagErrorf("--description-json items must be JSON objects")
		}
		typ, _ := obj["type"].(string)
		switch typ {
		case "text":
			if _, ok := obj["text"].(string); !ok {
				return baseFlagErrorf("--description-json text items require text")
			}
		case "url":
			if _, ok := obj["text"].(string); !ok {
				return baseFlagErrorf("--description-json url items require text")
			}
			if _, ok := obj["url"].(string); !ok {
				return baseFlagErrorf("--description-json url items require url")
			}
		case "mention":
			if _, ok := obj["open_id"].(string); !ok {
				return baseFlagErrorf("--description-json mention items require open_id")
			}
		default:
			return baseFlagErrorf("--description-json item type must be text, url, or mention")
		}
	}
	return nil
}

func anyChanged(runtime *common.RuntimeContext, names ...string) bool {
	for _, name := range names {
		if runtime.Changed(name) {
			return true
		}
	}
	return false
}

func containsKey(v interface{}, key string) bool {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, child := range x {
			if k == key || containsKey(child, key) {
				return true
			}
		}
	case []interface{}:
		for _, child := range x {
			if containsKey(child, key) {
				return true
			}
		}
	}
	return false
}
