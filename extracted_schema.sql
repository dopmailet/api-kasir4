-- EXTRACTED SCHEMA FROM LIMA DATABASE --

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    category character varying(50) NOT NULL,
    description character varying(255) NOT NULL,
    amount numeric NOT NULL,
    expense_date date NOT NULL,
    is_recurring boolean DEFAULT false,
    recurring_period character varying(20),
    notes text,
    created_by integer,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    total_amount numeric NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    kasir_id integer,
    discount_id integer,
    discount_amount numeric DEFAULT 0,
    payment_amount numeric DEFAULT 0,
    change_amount numeric DEFAULT 0,
    created_by integer
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    supplier_name character varying(150),
    total_amount numeric DEFAULT 0 NOT NULL,
    notes text,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE purchase_items (
    id SERIAL PRIMARY KEY,
    purchase_id integer NOT NULL,
    product_id integer,
    product_name character varying(150) NOT NULL,
    quantity integer NOT NULL,
    buy_price numeric NOT NULL,
    sell_price numeric,
    category_id integer,
    subtotal numeric NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    nama character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    discount_type character varying(20) DEFAULT NULL::character varying,
    discount_value numeric DEFAULT 0
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    nama character varying(255) NOT NULL,
    harga integer NOT NULL,
    stok integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    category_id integer,
    harga_beli numeric,
    created_by integer,
    default_discount_type character varying(20),
    default_discount_value numeric,
    barcode character varying(100),
    is_featured boolean DEFAULT false
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    nama_lengkap character varying(100) NOT NULL,
    role character varying(20) DEFAULT 'kasir'::character varying NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE discounts (
    id SERIAL PRIMARY KEY,
    name character varying(100) NOT NULL,
    type character varying(20) NOT NULL,
    value numeric NOT NULL,
    min_order_amount numeric DEFAULT 0,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    product_id integer,
    category_id integer
);

CREATE TABLE payroll (
    id SERIAL PRIMARY KEY,
    employee_id integer NOT NULL,
    periode character varying(20),
    gaji_pokok numeric NOT NULL,
    bonus numeric DEFAULT 0,
    potongan numeric DEFAULT 0,
    total numeric NOT NULL,
    catatan text,
    paid_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id integer NOT NULL,
    product_id integer NOT NULL,
    quantity integer NOT NULL,
    price numeric NOT NULL,
    subtotal numeric NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    harga_beli numeric,
    discount_type character varying(10) DEFAULT NULL::character varying,
    discount_value numeric DEFAULT 0,
    discount_amount numeric DEFAULT 0
);

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    nama character varying(100) NOT NULL,
    posisi character varying(50) NOT NULL,
    gaji_pokok numeric DEFAULT 0 NOT NULL,
    no_hp character varying(20),
    alamat text,
    tanggal_masuk date DEFAULT CURRENT_DATE,
    aktif boolean DEFAULT true,
    user_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

