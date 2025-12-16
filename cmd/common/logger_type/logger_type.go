package loggertype

const (
	TenantActionCreate = "TENANT_CREATE"
	TenantActionUpdate = "TENANT_UPDATE"
	TenantActionDelete = "TENANT_DELETE"
	UserActionCreate   = "USER_CREATE"
	UserActionUpdate   = "USER_UPDATE"
	UserActionDelete   = "USER_DELETE"
)

func InfoTenantCreate() string {
	return TenantActionCreate
}
