-- Add tenant-safe references without changing historical migrations or IDs.
ALTER TABLE users ADD CONSTRAINT users_tenant_id_id_unique UNIQUE(tenant_id,id);
ALTER TABLE roles ADD CONSTRAINT roles_tenant_id_id_unique UNIQUE(tenant_id,id);
ALTER TABLE sessions ADD CONSTRAINT sessions_tenant_user_id_unique UNIQUE(tenant_id,user_id,id);
ALTER TABLE tenant_users ADD CONSTRAINT tenant_users_user_scope FOREIGN KEY(tenant_id,user_id) REFERENCES users(tenant_id,id);
ALTER TABLE tenant_users ADD CONSTRAINT tenant_users_role_scope FOREIGN KEY(tenant_id,role_id) REFERENCES roles(tenant_id,id);
ALTER TABLE sessions ADD CONSTRAINT sessions_user_scope FOREIGN KEY(tenant_id,user_id) REFERENCES users(tenant_id,id);
ALTER TABLE refresh_tokens ADD CONSTRAINT refresh_session_scope FOREIGN KEY(tenant_id,user_id,session_id) REFERENCES sessions(tenant_id,user_id,id);
ALTER TABLE audit_logs ADD CONSTRAINT audit_actor_scope FOREIGN KEY(tenant_id,actor_id) REFERENCES users(tenant_id,id);

ALTER TABLE fleet_vehicle_categories ADD CONSTRAINT fleet_category_scope_unique UNIQUE(tenant_id,id);
ALTER TABLE fleet_vehicle_types ADD CONSTRAINT fleet_type_scope_unique UNIQUE(tenant_id,id);
ALTER TABLE fleet_vehicle_brands ADD CONSTRAINT fleet_brand_scope_unique UNIQUE(tenant_id,id);
ALTER TABLE fleet_vehicle_models ADD CONSTRAINT fleet_model_scope_unique UNIQUE(tenant_id,id);
ALTER TABLE fleet_assets ADD CONSTRAINT fleet_asset_scope_unique UNIQUE(tenant_id,id);
ALTER TABLE fleet_vehicle_types ADD CONSTRAINT fleet_type_category_scope FOREIGN KEY(tenant_id,category_id) REFERENCES fleet_vehicle_categories(tenant_id,id);
ALTER TABLE fleet_vehicle_models ADD CONSTRAINT fleet_model_brand_scope FOREIGN KEY(tenant_id,brand_id) REFERENCES fleet_vehicle_brands(tenant_id,id);
ALTER TABLE fleet_vehicle_models ADD CONSTRAINT fleet_model_category_scope FOREIGN KEY(tenant_id,category_id) REFERENCES fleet_vehicle_categories(tenant_id,id);
ALTER TABLE fleet_vehicle_models ADD CONSTRAINT fleet_model_type_scope FOREIGN KEY(tenant_id,type_id) REFERENCES fleet_vehicle_types(tenant_id,id);
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_category_scope FOREIGN KEY(tenant_id,category_id) REFERENCES fleet_vehicle_categories(tenant_id,id);
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_type_scope FOREIGN KEY(tenant_id,type_id) REFERENCES fleet_vehicle_types(tenant_id,id);
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_brand_scope FOREIGN KEY(tenant_id,brand_id) REFERENCES fleet_vehicle_brands(tenant_id,id);
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_model_scope FOREIGN KEY(tenant_id,model_id) REFERENCES fleet_vehicle_models(tenant_id,id);
ALTER TABLE fleet_vehicles ADD CONSTRAINT fleet_vehicle_asset_scope FOREIGN KEY(tenant_id,asset_id) REFERENCES fleet_assets(tenant_id,id);
