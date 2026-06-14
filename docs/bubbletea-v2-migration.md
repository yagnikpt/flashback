# Bubble Tea v2 Migration (Flashback)

## Web-sourced migration checklist (applied)

Sources used:
- https://github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md
- https://github.com/charmbracelet/bubbles/blob/main/UPGRADE_GUIDE_V2.md
- https://github.com/charmbracelet/bubbletea/releases/tag/v2.0.0

1. **Upgrade module paths and versions together**
   - `charm.land/bubbletea/v2`
   - `charm.land/bubbles/v2`
   - `charm.land/lipgloss/v2`
   - Run `go mod tidy`.

2. **Switch imports from `github.com/charmbracelet/...` to `charm.land/.../v2`**
   - Bubble Tea, Bubbles subpackages (`list`, `key`, `spinner`, `textarea`), and Lip Gloss.

3. **Update keyboard message types**
   - Replace `tea.KeyMsg` handling with `tea.KeyPressMsg` in all `Update` switches.

4. **Migrate views from `string` to `tea.View`**
   - Change all `View() string` methods to `View() tea.View`.
   - Wrap rendered strings with `tea.NewView(...)`.

5. **Move alt screen/window title behavior to declarative view fields**
   - Remove v1 options/commands like `tea.WithAltScreen()` and `tea.SetWindowTitle(...)`.
   - Set fields on returned `tea.View` instead:
     - `v.AltScreen = true`
     - `v.WindowTitle = "flashback"`

6. **Adjust Bubbles textarea styling API usage**
   - In v2, use `styles := model.Styles()` + `model.SetStyles(styles)`.
   - Do not write directly through old style fields.

7. **Reformat and validate**
   - Run `gofmt`.
   - Run `go test ./...` (passes).

## Notes
- Space key matching should be `"space"` if used as a string-match key. (No migration needed in this codebase.)
- `go.mod`/`go.sum` were updated by `go mod tidy` to reflect the new dependency graph.
