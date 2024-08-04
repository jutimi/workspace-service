CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE workspaces (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    name_slug VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    address TEXT,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at BIGINT
);

CREATE TYPE user_workspace_role AS ENUM ('admin', 'user', 'owner');
CREATE TABLE user_workspaces (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(20),
    email VARCHAR(100),
    role user_workspace_role NOT NULL,
    is_active BOOLEAN DEFAULT true,
    user_id uuid NOT NULL,
    workspace_id uuid NOT NULL,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at BIGINT
);

CREATE TABLE user_workspace_details (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    full_name_slug VARCHAR(100) NOT NULL,
    user_workspace_id uuid NOT NULL,
    workspace_id uuid NOT NULL,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at BIGINT
);

CREATE TABLE organizations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    name_slug VARCHAR(100) NOT NULL,
    level INTEGER NOT NULL,
    parent_organization_ids TEXT,
    parent_organization_id uuid,
    manager_id uuid,
    workspace_id uuid NOT NULL,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at BIGINT
);

CREATE TYPE user_workspace_organization_role AS ENUM ('leader', 'sub_leader', 'member');
CREATE TABLE user_workspace_organizations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id uuid NOT NULL,
    role user_workspace_organization_role NOT NULL,
    leader_ids TEXT,
    user_workspace_id uuid NOT NULL,
    workspace_id uuid NOT NULL,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at BIGINT
);