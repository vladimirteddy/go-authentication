-- +goose Up
-- +goose StatementBegin

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    UNIQUE(resource, action)
);

-- Create user_roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES 
    ('admin', 'Administrator with full access to all resources'),
    ('user', 'Regular user with limited access')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (resource, action, description) VALUES 
    ('users', 'read', 'Read user information'),
    ('users', 'create', 'Create new users'),
    ('users', 'update', 'Update user information'),
    ('users', 'delete', 'Delete users'),
    ('roles', 'read', 'Read role information'),
    ('roles', 'create', 'Create new roles'),
    ('roles', 'update', 'Update role information'),
    ('roles', 'delete', 'Delete roles'),
    ('permissions', 'read', 'Read permission information'),
    ('permissions', 'create', 'Create new permissions'),
    ('permissions', 'update', 'Update permission information'),
    ('permissions', 'delete', 'Delete permissions')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign all permissions to admin role
DO $$
DECLARE
    admin_role_id INTEGER;
BEGIN
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    
    -- Insert admin role permissions
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT admin_role_id, id FROM permissions
    ON CONFLICT DO NOTHING;
    
    -- Assign read permissions to the user role
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT 
        (SELECT id FROM roles WHERE name = 'user'),
        id
    FROM permissions 
    WHERE action = 'read'
    ON CONFLICT DO NOTHING;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

-- +goose StatementEnd 