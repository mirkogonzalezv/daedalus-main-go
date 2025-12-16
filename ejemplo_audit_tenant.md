# Ejemplo de Audit Log para Creación de Tenants

## Configuración Implementada

Se ha integrado el `AuditService` en el `TenantUseCase` para registrar automáticamente los logs de auditoría cuando se crea un tenant.

### Cambios Realizados:

1. **TenantUseCase**: Ahora incluye `AuditService` como dependencia
2. **Container**: Configurado para inyectar el `AuditService` 
3. **Audit Log**: Se registra automáticamente al crear un tenant

## Ejemplo de Uso

### 1. Crear un Tenant (Request)

```bash
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mi Empresa",
    "slug": "mi-empresa", 
    "plan": "pro"
  }'
```

### 2. Response Esperado

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Mi Empresa",
  "slug": "mi-empresa",
  "subscription_plan": "pro",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 3. Audit Log Generado Automáticamente

El sistema registrará automáticamente en la tabla `audit_logs`:

```sql
SELECT * FROM daedalus.audit_logs WHERE action = 'tenant.created';
```

**Resultado esperado:**
```
id: uuid-generado
tenant_id: 550e8400-e29b-41d4-a716-446655440000
user_id: NULL (porque es creación de tenant, no hay usuario aún)
action: tenant.created
meta: {
  "tenant_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_name": "Mi Empresa", 
  "tenant_slug": "mi-empresa",
  "plan": "pro"
}
ip: NULL (opcional)
created_at: 2024-01-15T10:30:00Z
```

## Flujo de Ejecución

1. **Controller** recibe el request
2. **TenantUseCase.CreateTenant()** valida y crea el tenant
3. **Repository** guarda el tenant en la base de datos
4. **AuditService.LogAction()** registra el log de auditoría (asíncrono)
5. **Controller** retorna la respuesta

## Ventajas de esta Implementación

- ✅ **Automático**: No requiere código adicional en el controller
- ✅ **Asíncrono**: El audit log no bloquea la respuesta
- ✅ **Consistente**: Siempre se registra cuando se crea un tenant
- ✅ **Metadata Rica**: Incluye información detallada del tenant creado
- ✅ **Escalable**: Fácil de extender a otras acciones (update, delete, etc.)

## Próximos Pasos

Para extender el audit log a otras operaciones:

1. **Actualización de Tenant**: Agregar audit log en `ActualizarTenant()`
2. **Creación de Usuario**: Implementar en `UserUseCase`
3. **Login/Logout**: Registrar eventos de autenticación
4. **Operaciones Críticas**: Auditar cambios de permisos, configuraciones, etc.