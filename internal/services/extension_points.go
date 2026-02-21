package services

// ExtensionPoint defines where plugins can inject UI components
type ExtensionPoint string

const (
	// Page-level: Full MainView replacement
	ExtensionPointGlobalPage ExtensionPoint = "Global.Page"

	// Panel-level: Sidebar sections
	ExtensionPointPoolSidebarSection      ExtensionPoint = "Pool.Sidebar.Section"
	ExtensionPointWorkspaceSidebarSection ExtensionPoint = "Workspace.Sidebar.Section"
	ExtensionPointArtifactSidebarSection  ExtensionPoint = "Artifact.Sidebar.Section"
	ExtensionPointRightsSidebarSection    ExtensionPoint = "Rights.Sidebar.Section"

	// Panel-level: Inspector tabs
	ExtensionPointPoolInspectorTab      ExtensionPoint = "Pool.Inspector.Tab"
	ExtensionPointWorkspaceInspectorTab ExtensionPoint = "Workspace.Inspector.Tab"
	ExtensionPointArtifactInspectorTab  ExtensionPoint = "Artifact.Inspector.Tab"
	ExtensionPointRightsInspectorTab    ExtensionPoint = "Rights.Inspector.Tab"

	// Action-level: Toolbar buttons
	ExtensionPointPoolToolbarAction      ExtensionPoint = "Pool.Toolbar.Action"
	ExtensionPointWorkspaceToolbarAction ExtensionPoint = "Workspace.Toolbar.Action"
	ExtensionPointArtifactToolbarAction  ExtensionPoint = "Artifact.Toolbar.Action"
	ExtensionPointRightsToolbarAction    ExtensionPoint = "Rights.Toolbar.Action"

	// Action-level: Context menu items
	ExtensionPointPoolContextMenuItem      ExtensionPoint = "Pool.ContextMenu.Item"
	ExtensionPointWorkspaceContextMenuItem ExtensionPoint = "Workspace.ContextMenu.Item"
	ExtensionPointArtifactContextMenuItem  ExtensionPoint = "Artifact.ContextMenu.Item"

	// StatusBar-level: Bottom bar widgets
	ExtensionPointStatusBarLeft  ExtensionPoint = "Global.StatusBar.Left"
	ExtensionPointStatusBarRight ExtensionPoint = "Global.StatusBar.Right"
)

// ValidExtensionPoints maps all valid extension points
var ValidExtensionPoints = map[ExtensionPoint]string{
	ExtensionPointGlobalPage: "全局页面（独占MainView）",

	ExtensionPointPoolSidebarSection:      "素材库左侧边栏面板",
	ExtensionPointWorkspaceSidebarSection: "工作台左侧边栏面板",
	ExtensionPointArtifactSidebarSection:  "成品库左侧边栏面板",
	ExtensionPointRightsSidebarSection:    "版权中心左侧边栏面板",

	ExtensionPointPoolInspectorTab:      "素材库右侧检查器标签页",
	ExtensionPointWorkspaceInspectorTab: "工作台右侧检查器标签页",
	ExtensionPointArtifactInspectorTab:  "成品库右侧检查器标签页",
	ExtensionPointRightsInspectorTab:    "版权中心右侧检查器标签页",

	ExtensionPointPoolToolbarAction:      "素材库工具栏按钮",
	ExtensionPointWorkspaceToolbarAction: "工作台工具栏按钮",
	ExtensionPointArtifactToolbarAction:  "成品库工具栏按钮",
	ExtensionPointRightsToolbarAction:    "版权中心工具栏按钮",

	ExtensionPointPoolContextMenuItem:      "素材右键菜单项",
	ExtensionPointWorkspaceContextMenuItem: "工作台右键菜单项",
	ExtensionPointArtifactContextMenuItem:  "成品右键菜单项",

	ExtensionPointStatusBarLeft:  "状态栏左侧小部件",
	ExtensionPointStatusBarRight: "状态栏右侧小部件",
}

// IsValidExtensionPoint checks if a slot is valid
func IsValidExtensionPoint(slot string) bool {
	_, ok := ValidExtensionPoints[ExtensionPoint(slot)]
	return ok
}

// GetExtensionPointDescription returns human-readable description
func GetExtensionPointDescription(slot string) string {
	return ValidExtensionPoints[ExtensionPoint(slot)]
}
