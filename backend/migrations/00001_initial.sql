-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('system', 'owner', 'admin', 'doctor', 'nurse', 'technician', 'reception');
    CREATE TYPE appt_status AS ENUM ('scheduled', 'waiting', 'in_progress', 'completed', 'cancelled');
    CREATE TYPE pay_method AS ENUM ('cash', 'card', 'click', 'payme', 'debt');
    CREATE TYPE lab_status AS ENUM ('ordered', 'received', 'in_process', 'ready', 'delivered');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    owner_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    subscription_end_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    UNIQUE(tenant_id, slug)
);

CREATE SEQUENCE IF NOT EXISTS global_staff_seq START 1;
CREATE SEQUENCE IF NOT EXISTS global_patient_seq START 1;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    full_name VARCHAR(100) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role user_role NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE staff_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    display_id VARCHAR(20) UNIQUE, 
    primary_branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    specialty VARCHAR(100),
    room_number VARCHAR(10),
    percentage_share DECIMAL(5,2) DEFAULT 0
);

CREATE TABLE doctor_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staff_id UUID REFERENCES staff_profiles(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL, 
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);

CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    display_id VARCHAR(20) UNIQUE, 
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    birth_date DATE,
    gender VARCHAR(10),
    balance DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id UUID REFERENCES staff_profiles(id) ON DELETE SET NULL,
    scheduled_time TIMESTAMP NOT NULL,
    status appt_status DEFAULT 'scheduled',
    queue_number INT,
    complaint TEXT,
    diagnosis TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    code VARCHAR(20)
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    patient_id UUID REFERENCES patients(id),
    appointment_id UUID REFERENCES appointments(id),
    amount DECIMAL(15, 2) NOT NULL,
    method pay_method NOT NULL,
    transaction_ref VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    min_alert_qty INT DEFAULT 10
);

CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    quantity DECIMAL(10, 3) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, product_id)
);

CREATE TABLE service_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID REFERENCES services(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    quantity_required DECIMAL(10, 3) NOT NULL
);

CREATE TABLE lab_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    doctor_id UUID REFERENCES staff_profiles(id),
    technician_id UUID REFERENCES staff_profiles(id),
    patient_id UUID REFERENCES patients(id),
    item_name VARCHAR(150) NOT NULL,
    price DECIMAL(15, 2),
    status lab_status DEFAULT 'ordered',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- TRIGGERS
CREATE OR REPLACE FUNCTION generate_staff_display_id()
RETURNS TRIGGER AS $$
DECLARE
    role_prefix CHAR(1);
    seq_val BIGINT;
    user_role_val user_role;
BEGIN
    SELECT role INTO user_role_val FROM users WHERE id = NEW.user_id;
    
    CASE user_role_val
        WHEN 'doctor' THEN role_prefix := 'D';
        WHEN 'nurse' THEN role_prefix := 'N';
        WHEN 'technician' THEN role_prefix := 'T';
        WHEN 'reception' THEN role_prefix := 'R';
        ELSE role_prefix := 'S';
    END CASE;

    seq_val := nextval('global_staff_seq');
    NEW.display_id := role_prefix || LPAD(seq_val::TEXT, 6, '0');
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_staff_display_id
BEFORE INSERT ON staff_profiles
FOR EACH ROW
EXECUTE FUNCTION generate_staff_display_id();

CREATE OR REPLACE FUNCTION generate_patient_display_id()
RETURNS TRIGGER AS $$
DECLARE
    seq_val BIGINT;
BEGIN
    seq_val := nextval('global_patient_seq');
    NEW.display_id := 'P' || LPAD(seq_val::TEXT, 6, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_patient_display_id
BEFORE INSERT ON patients
FOR EACH ROW
EXECUTE FUNCTION generate_patient_display_id();

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_patient_display_id ON patients;
DROP TRIGGER IF EXISTS set_staff_display_id ON staff_profiles;

DROP FUNCTION IF EXISTS generate_patient_display_id;
DROP FUNCTION IF EXISTS generate_staff_display_id;

DROP TABLE IF EXISTS lab_orders;
DROP TABLE IF EXISTS service_recipes;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS doctor_schedules;
DROP TABLE IF EXISTS staff_profiles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS branches;
DROP TABLE IF EXISTS tenants;

DROP SEQUENCE IF EXISTS global_patient_seq;
DROP SEQUENCE IF EXISTS global_staff_seq;

DROP TYPE IF EXISTS lab_status;
DROP TYPE IF EXISTS pay_method;
DROP TYPE IF EXISTS appt_status;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose StatementEnd