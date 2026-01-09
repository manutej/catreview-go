```mermaid
flowchart BT
    subgraph L0["Foundations"]
        import_fmt["import:fmt<br/>A:94 E:0"]
        import_strings["import:strings<br/>A:85 E:0"]
        import_context["import:context<br/>A:80 E:0"]
        import_os["import:os<br/>A:53 E:0"]
        import_time["import:time<br/>A:53 E:0"]
        import_charm_land_lipgloss_v2["v2<br/>A:49 E:0"]
        config_Get["config.Get<br/>A:44 E:0"]
        import_charm_land_bubbletea_v2["v2<br/>A:43 E:0"]
        import_github_com_charmbracelet_crush_internal_tui_styles["styles<br/>A:39 E:0"]
        import_path_filepath["filepath<br/>A:38 E:0"]
    end
    subgraph L1["Core"]
        internal_tui_components_chat_messages_renderer_go["renderer.go<br/>A:0 E:93"]
        internal_tui_exp_list_list_go["list.go<br/>A:0 E:82"]
        internal_config_config_go["config.go<br/>A:0 E:79"]
        internal_tui_page_chat_chat_go["chat.go<br/>A:0 E:70"]
        internal_agent_coordinator_go["coordinator.go<br/>A:0 E:70"]
        internal_tui_components_chat_messages_tool_go["tool.go<br/>A:0 E:69"]
        internal_tui_components_chat_splash_splash_go["splash.go<br/>A:0 E:67"]
        internal_tui_components_chat_editor_editor_go["editor.go<br/>A:0 E:63"]
        internal_message_content_go["content.go<br/>A:0 E:62"]
        internal_tui_components_chat_chat_go["chat.go<br/>A:0 E:61"]
    end
    subgraph L2["Support"]
        internal_tui_exp_list_filterable_go["filterable.go<br/>A:0 E:40"]
        internal_tui_components_dialogs_models_models_go["models.go<br/>A:0 E:39"]
        internal_uicmd_uicmd_go["uicmd.go<br/>A:0 E:39"]
        internal_cmd_root_go["root.go<br/>A:0 E:36"]
        internal_tui_styles_theme_go["theme.go<br/>A:0 E:35"]
        internal_agent_tools_grep_go["grep.go<br/>A:0 E:33"]
        internal_tui_components_dialogs_reasoning_reasoning_go["reasoning.go<br/>A:0 E:33"]
        internal_tui_components_dialogs_filepicker_filepicker_go["filepicker.go<br/>A:0 E:32"]
        internal_agent_prompt_prompt_go["prompt.go<br/>A:0 E:32"]
        internal_tui_components_completions_completions_go["completions.go<br/>A:0 E:30"]
    end
    subgraph L3["Applications"]
        main_go["main.go<br/>A:0 E:7"]
    end

    internal_agent_tools_diagnostics_go --> import_log_slog
    internal_cmd_update_providers_go --> import_github_com_charmbracelet_crush_internal_config
    internal_db_connect_go --> import_github_com_ncruces_go_sqlite3
    internal_db_connect_go --> import_path_filepath
    internal_agent_tools_search_go --> tools_SearchResult
    internal_tui_components_dialogs_permissions_permissions_go --> permissions_NewPermissionDialogCmp
    internal_agent_prompts_go --> import_embed
    internal_tui_components_chat_splash_splash_go --> import_charm_land_bubbles_v2_spinner
    internal_tui_components_chat_messages_messages_go --> messages__messageCmp_Update
    internal_tui_components_core_status_status_go --> status__statusCmp_ToggleFullHelp
    internal_update_update_go --> update_Info_Available
    internal_agent_tools_view_go --> import_io
    internal_tui_components_dialogs_keys_go --> import_charm_land_bubbles_v2_key
    internal_skills_skills_go --> skills_splitFrontmatter
    internal_tui_components_dialogs_dialogs_go --> dialogs_dialogCmp
    internal_tui_exp_list_list_go --> list_renderedItem
    internal_uicmd_uicmd_go --> uicmd_createCommandHandler
    internal_config_config_go --> config_MCPConfig
    internal_tui_components_chat_messages_tool_go --> messages__toolCallCmp_copyTool
    internal_tui_components_image_image_go --> image_Model_Init
    internal_agent_tools_mcp_init_go --> import_github_com_charmbracelet_crush_internal_csync
    internal_session_session_go --> session__service_List
    internal_tui_exp_diffview_diffview_go --> diffview__DiffView_normalizeLineEndings
    internal_tui_exp_list_list_go --> list_List
    internal_app_app_go --> import_github_com_charmbracelet_x_exp_charmtone
    internal_tui_components_dialogs_models_models_go --> models__modelDialogCmp_modelTypeRadio
    internal_session_session_go --> import_github_com_charmbracelet_crush_internal_event
    main_go --> import_os
    internal_config_config_go --> config_resolveAllowedTools
    internal_tui_components_dialogs_claude_oauth_go --> claude__OAuth2_SetDefaults
    internal_tui_exp_diffview_diffview_go --> diffview__DiffView_Height
    internal_tui_exp_list_filterable_group_go --> list___SetGroups
    internal_agent_tools_references_go --> import_path_filepath
    internal_cmd_logs_go --> import_github_com_nxadm_tail
    internal_tui_page_chat_chat_go --> import_github_com_charmbracelet_crush_internal_tui_components_dialogs_filepicker
    internal_config_provider_go --> config_UpdateHyper
    internal_message_content_go --> message_FinishReason
    internal_agent_tools_write_go --> import_github_com_charmbracelet_crush_internal_diff
    internal_tui_components_chat_chat_go --> chat__messageListCmp_GetSelectedText
    internal_tui_components_chat_splash_splash_go --> splash__splashCmp_mcpBlock
    internal_oauth_copilot_oauth_go --> import_fmt
    internal_tui_components_chat_sidebar_sidebar_go --> import_charm_land_lipgloss_v2
    internal_tui_components_completions_completions_go --> import_github_com_charmbracelet_crush_internal_tui_exp_list
    internal_tui_components_dialogs_models_models_go --> import_charm_land_bubbles_v2_help
    internal_agent_coordinator_go --> agent__coordinator_buildBedrockProvider
    internal_tui_components_dialogs_claude_method_go --> claude_AuthMethodChooser
    internal_tui_components_anim_anim_go --> anim__Anim_SetLabel
    internal_tui_components_chat_messages_tool_go --> messages__toolCallCmp_SetSize
    internal_agent_tools_glob_go --> import_path_filepath
    internal_agent_tools_mcp_init_go --> import_cmp

    style L0 fill:#e8f5e9
    style L1 fill:#fff3e0
    style L2 fill:#e3f2fd
    style L3 fill:#fce4ec
```
