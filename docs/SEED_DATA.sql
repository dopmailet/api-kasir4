-- ============================================================
-- SEED DATA LENGKAP - KasirPOS API
-- Jalankan di Supabase SQL Editor
-- Data dummy realistis untuk testing & demo
-- ============================================================
-- URUTAN EKSEKUSI (ikuti urutan karena ada foreign key):
--   1. categories
--   2. products
--   3. discounts
--   4. employees
--   5. transactions + transaction_details
--   6. purchases + purchase_items
--   7. expenses
--   8. payroll
-- ============================================================


-- ============================================================
-- 1. CATEGORIES (10 kategori)
-- ============================================================
INSERT INTO categories (nama, description, discount_type, discount_value) VALUES
  ('Minuman',        'Berbagai jenis minuman segar dan kemasan',             NULL,           0),
  ('Makanan',        'Makanan ringan dan berat siap saji',                   NULL,           0),
  ('Snack',          'Camilan dan makanan ringan kemasan',                   'percentage',   5),
  ('Rokok',          'Produk tembakau dan aksesoris rokok',                  NULL,           0),
  ('Kebersihan',     'Produk kebersihan dan perawatan tubuh',                'fixed',     1000),
  ('Elektronik',     'Aksesoris dan perangkat elektronik kecil',             NULL,           0),
  ('Frozen Food',    'Makanan beku siap saji',                               NULL,           0),
  ('Bumbu & Dapur',  'Bumbu masak, minyak, dan kebutuhan dapur',             NULL,           0),
  ('Obat & Vitamin', 'Obat bebas, vitamin, dan suplemen kesehatan',          'percentage',  10),
  ('Alat Tulis',     'Peralatan tulis kantor dan sekolah',                   NULL,           0)
ON CONFLICT DO NOTHING;


-- ============================================================
-- 2. PRODUCTS (50 produk)
-- ============================================================
INSERT INTO products (nama, harga, harga_beli, stok, barcode, category_id, is_featured, created_by) VALUES
  -- MINUMAN (category_id = 1)
  ('Aqua 600ml',              4000,   2500,  200, '8000000010001', 1, true,  1),
  ('Aqua 1500ml',             7000,   4500,  150, '8000000010002', 1, false, 1),
  ('Teh Botol Sosro 450ml',   5000,   3200,  130, '8000000010003', 1, true,  1),
  ('Pocari Sweat 500ml',      8000,   5500,  100, '8000000010004', 1, true,  1),
  ('Coca Cola 330ml',         7000,   4800,  120, '8000000010005', 1, false, 1),
  ('Sprite 330ml',            7000,   4800,   90, '8000000010006', 1, false, 1),
  ('Fanta Strawberry 330ml',  7000,   4800,   80, '8000000010007', 1, false, 1),
  ('Kopi Good Day 200ml',     4500,   2800,  110, '8000000010008', 1, false, 1),
  ('Milo Actigen 200ml',      6000,   4000,   95, '8000000010009', 1, false, 1),
  ('Mizone Apple 500ml',      6500,   4200,   75, '8000000010010', 1, false, 1),

  -- MAKANAN (category_id = 2)
  ('Indomie Goreng Spesial',  3500,   2200,  250, '8000000020001', 2, true,  1),
  ('Indomie Soto Ayam',       3500,   2200,  200, '8000000020002', 2, false, 1),
  ('Mie Sedaap Goreng',       3000,   1900,  180, '8000000020003', 2, false, 1),
  ('Roti Tawar Sari Roti',   15000,   9500,   50, '8000000020004', 2, false, 1),
  ('Beng-Beng Share It',      3000,   1800,  150, '8000000020005', 2, false, 1),
  ('Biskuit Roma Kelapa',     8500,   5500,   80, '8000000020006', 2, false, 1),
  ('Pop Mie Cup 75gr',        5500,   3500,  100, '8000000020007', 2, false, 1),
  ('Roti Coklat Sari Roti',  12000,   7800,   60, '8000000020008', 2, false, 1),
  ('Makaroni Pedas Mamang',   5000,   3200,  120, '8000000020009', 2, false, 1),
  ('Permen Relaxa',           1000,    500,  300, '8000000020010', 2, false, 1),

  -- SNACK (category_id = 3)
  ('Chitato Original 68gr',  11000,   7000,   70, '8000000030001', 3, true,  1),
  ('Chitato Sapi Panggang',  11000,   7000,   65, '8000000030002', 3, false, 1),
  ('Lays Classic 65gr',      12000,   7800,   60, '8000000030003', 3, false, 1),
  ('Oreo Vanila 133gr',       8000,   5000,   90, '8000000030004', 3, false, 1),
  ('Oreo Stroberi 133gr',     8000,   5000,   85, '8000000030005', 3, false, 1),
  ('Taro Net 115gr',         10000,   6500,   55, '8000000030006', 3, false, 1),
  ('Cheetos Jagung 60gr',     9000,   5800,   70, '8000000030007', 3, false, 1),
  ('Qtela Singkong 65gr',    10000,   6500,   50, '8000000030008', 3, false, 1),
  
  -- ROKOK (category_id = 4)
  ('Gudang Garam Surya 16',  28000,  25000,   60, '8000000040001', 4, false, 1),
  ('Djarum Super 16',        27000,  24000,   55, '8000000040002', 4, false, 1),
  ('Marlboro Red 20',        42000,  38000,   40, '8000000040003', 4, false, 1),
  ('Sampoerna Mild 16',      30000,  27000,   50, '8000000040004', 4, false, 1),

  -- KEBERSIHAN (category_id = 5)
  ('Sabun Lifebuoy 110gr',    8000,   5000,  100, '8000000050001', 5, false, 1),
  ('Shampoo Sunsilk 170ml',  16000,  10500,   60, '8000000050002', 5, false, 1),
  ('Pasta Gigi Pepsodent',   12000,   7800,   70, '8000000050003', 5, false, 1),
  ('Dettol Hand Sanitizer',  18000,  12000,   45, '8000000050004', 5, false, 1),
  ('Rinso Bubuk 800gr',      22000,  15000,   40, '8000000050005', 5, false, 1),
  ('Molto Pelembut 900ml',   20000,  13500,   35, '8000000050006', 5, false, 1),

  -- ELEKTRONIK (category_id = 6)
  ('Baterai ABC AA (2pcs)',    8000,   5500,   50, '8000000060001', 6, false, 1),
  ('Earphone Unipin EP C11', 25000,  16000,   25, '8000000060002', 6, false, 1),
  ('Kabel USB Type-C 1m',    18000,  11000,   30, '8000000060003', 6, false, 1),
  ('Charger Dus 1A',         25000,  16000,   20, '8000000060004', 6, false, 1),

  -- FROZEN FOOD (category_id = 7)
  ('Nugget Fiesta 500gr',    35000,  25000,   35, '8000000070001', 7, true,  1),
  ('Sosis Ayam So Nice 375gr',28000, 19000,   40, '8000000070002', 7, false, 1),
  ('Bakso Sapi Kanzler',     32000,  22000,   30, '8000000070003', 7, false, 1),

  -- BUMBU & DAPUR (category_id = 8)
  ('Minyak Goreng Bimoli 2L',32000,  25000,   45, '8000000080001', 8, false, 1),
  ('Kecap Bango 220ml',      15000,  10000,   60, '8000000080002', 8, false, 1),
  ('Saos Sambal ABC 335ml',  16000,  11000,   55, '8000000080003', 8, false, 1),

  -- OBAT & VITAMIN (category_id = 9)
  ('Paracetamol Generik 10s', 5000,   3000,   80, '8000000090001', 9, false, 1),
  ('Antangin JRG Sachet',     4000,   2500,   90, '8000000090002', 9, false, 1)
ON CONFLICT DO NOTHING;


-- ============================================================
-- 3. DISCOUNTS (10 diskon)
-- ============================================================
INSERT INTO discounts (name, type, value, min_order_amount, product_id, category_id, start_date, end_date, is_active) VALUES
  ('Diskon Weekend 10%',         'PERCENTAGE', 10,    50000, NULL, NULL, NOW() - INTERVAL '30 days', NOW() + INTERVAL '60 days',  true),
  ('Potongan Rp 5.000',          'FIXED',      5000,  30000, NULL, NULL, NOW() - INTERVAL '30 days', NOW() + INTERVAL '60 days',  true),
  ('Promo Minuman 15%',          'PERCENTAGE', 15,    20000, NULL, 1,    NOW() - INTERVAL '7 days',  NOW() + INTERVAL '30 days',  true),
  ('Diskon Snack 5%',            'PERCENTAGE',  5,    10000, NULL, 3,    NOW() - INTERVAL '7 days',  NOW() + INTERVAL '30 days',  true),
  ('Flash Sale Aqua 600ml',      'FIXED',      1000,      0, 1,   NULL, NOW() - INTERVAL '1 days',  NOW() + INTERVAL '7 days',   true),
  ('Promo Frozen Food 8%',       'PERCENTAGE',  8,    25000, NULL, 7,    NOW() - INTERVAL '5 days',  NOW() + INTERVAL '25 days',  true),
  ('Diskon Obat 10%',            'PERCENTAGE', 10,     5000, NULL, 9,    NOW() - INTERVAL '5 days',  NOW() + INTERVAL '25 days',  true),
  ('Promo Kemerdekaan 17%',      'PERCENTAGE', 17,   100000, NULL, NULL, NOW() - INTERVAL '7 days',  NOW() + INTERVAL '14 days',  true),
  ('Potongan Rp 10.000',         'FIXED',      10000, 75000, NULL, NULL, NOW() - INTERVAL '14 days', NOW() + INTERVAL '14 days',  true),
  ('Diskon Lama (Expired)',       'PERCENTAGE', 20,   100000, NULL, NULL, NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 days',  false)
ON CONFLICT DO NOTHING;


-- ============================================================
-- 4. EMPLOYEES (8 karyawan)
-- ============================================================
INSERT INTO employees (nama, posisi, gaji_pokok, no_hp, alamat, tanggal_masuk, aktif, user_id) VALUES
  ('Budi Santoso',    'Kasir',          2500000, '081234567890', 'Jl. Melati No. 5, Jakarta Selatan',    '2024-01-15', true,  2),
  ('Siti Rahayu',     'Kasir',          2500000, '082345678901', 'Jl. Mawar No. 12, Jakarta Timur',      '2024-02-01', true,  3),
  ('Ahmad Fauzi',     'Stok Keeper',    2200000, '083456789012', 'Jl. Anggrek No. 8, Jakarta Barat',     '2023-11-01', true,  NULL),
  ('Dewi Lestari',    'Supervisor',     3500000, '084567890123', 'Jl. Kenanga No. 3, Jakarta Pusat',     '2023-06-01', true,  NULL),
  ('Rizky Pratama',   'Kasir',          2500000, '085678901234', 'Jl. Dahlia No. 17, Depok',             '2025-01-10', true,  NULL),
  ('Hana Fitriani',   'Admin Gudang',   2300000, '086789012345', 'Jl. Cempaka No. 22, Tangerang',        '2024-08-20', false, NULL),
  ('Eko Purnomo',     'Security',       2000000, '087890123456', 'Jl. Wijaya No. 9, Bekasi',             '2024-03-01', true,  NULL),
  ('Lina Marlina',    'Cleaning Staff', 1800000, '088901234567', 'Jl. Bougenville No. 31, Jakarta Utara','2024-05-15', true,  NULL)
ON CONFLICT DO NOTHING;


-- ============================================================
-- 5. TRANSACTIONS (50 transaksi, 3 bulan terakhir)
-- ============================================================
INSERT INTO transactions (total_amount, discount_id, discount_amount, payment_amount, change_amount, created_by, created_at) VALUES
  -- Januari 2026
  ( 36000, NULL,  0,      50000,  14000, 2, '2026-01-03 09:12:00+07'),
  ( 74500, 1,     7450,   80000,   5500, 3, '2026-01-03 11:30:00+07'),
  ( 53000, NULL,  0,      60000,   7000, 2, '2026-01-05 10:05:00+07'),
  (128000, 1,    12800,  150000,  22000, 2, '2026-01-05 14:45:00+07'),
  ( 47000, 2,     5000,   50000,   3000, 3, '2026-01-07 09:50:00+07'),
  ( 88000, NULL,  0,     100000,  12000, 2, '2026-01-08 16:20:00+07'),
  ( 33500, NULL,  0,      40000,   6500, 3, '2026-01-10 08:30:00+07'),
  ( 62000, 1,     6200,   70000,   8000, 2, '2026-01-12 13:15:00+07'),
  ( 44000, NULL,  0,      50000,   6000, 3, '2026-01-14 10:40:00+07'),
  (175000, 9,    10000,  200000,  25000, 2, '2026-01-15 15:00:00+07'),
  ( 29000, NULL,  0,      30000,   1000, 3, '2026-01-17 09:05:00+07'),
  ( 91500, 1,     9150,  100000,   9000, 2, '2026-01-18 11:35:00+07'),
  ( 38000, 2,     5000,   40000,   2000, 3, '2026-01-20 14:20:00+07'),
  ( 56000, NULL,  0,      60000,   4000, 2, '2026-01-21 16:50:00+07'),
  (110000, 8,    18700,  120000,  10000, 3, '2026-01-22 10:10:00+07'),
  ( 22000, NULL,  0,      30000,   8000, 2, '2026-01-24 09:30:00+07'),
  ( 67500, 1,     6750,   70000,   2500, 3, '2026-01-25 13:00:00+07'),
  ( 49000, NULL,  0,      50000,   1000, 2, '2026-01-27 11:15:00+07'),
  ( 83000, 2,     5000,   90000,   7000, 3, '2026-01-28 15:40:00+07'),
  (145000, 1,    14500,  150000,   5000, 2, '2026-01-30 17:00:00+07'),

  -- Februari 2026
  ( 41000, NULL,  0,      50000,   9000, 3, '2026-02-02 09:20:00+07'),
  ( 97000, 1,     9700,  100000,   3000, 2, '2026-02-03 12:05:00+07'),
  ( 58500, 2,     5000,   60000,   1500, 3, '2026-02-05 10:35:00+07'),
  (133000, 9,    10000,  150000,  17000, 2, '2026-02-06 14:50:00+07'),
  ( 26000, NULL,  0,      30000,   4000, 3, '2026-02-07 09:00:00+07'),
  ( 79000, 1,     7900,   80000,    100, 2, '2026-02-09 11:25:00+07'),
  ( 45000, NULL,  0,      50000,   5000, 3, '2026-02-10 13:40:00+07'),
  (115000, 8,    19550,  120000,    550, 2, '2026-02-12 16:15:00+07'),
  ( 37000, 2,     5000,   40000,   3000, 3, '2026-02-13 10:00:00+07'),
  ( 68000, NULL,  0,      70000,   2000, 2, '2026-02-14 12:30:00+07'),
  (200000, 1,    20000,  200000,      0, 3, '2026-02-15 15:00:00+07'),
  ( 51000, NULL,  0,      60000,   9000, 2, '2026-02-17 09:45:00+07'),
  ( 89500, 1,     8950,  100000,  10000, 3, '2026-02-18 14:00:00+07'),
  ( 32000, 2,     5000,   30000,  -2000, 2, '2026-02-20 11:10:00+07'),
  ( 76000, NULL,  0,      80000,   4000, 3, '2026-02-21 16:30:00+07'),
  (160000, 8,    27200,  150000,  -9800, 2, '2026-02-22 10:20:00+07'),
  ( 43500, NULL,  0,      50000,   6500, 3, '2026-02-24 13:55:00+07'),
  ( 94000, 1,     9400,  100000,   5000, 2, '2026-02-25 15:20:00+07'),
  ( 28000, NULL,  0,      30000,   2000, 3, '2026-02-26 09:30:00+07'),
  (155000, 9,    10000,  160000,   5000, 2, '2026-02-28 17:00:00+07'),

  -- Maret 2026
  ( 55000, NULL,  0,      60000,   5000, 2, '2026-03-01 10:00:00+07'),
  ( 87000, 1,     8700,   90000,   3000, 3, '2026-03-02 12:15:00+07'),
  ( 43000, 2,     5000,   40000,  -3000, 2, '2026-03-03 09:30:00+07'),
  (130000, 9,    10000,  130000,      0, 3, '2026-03-04 14:00:00+07'),
  ( 39000, NULL,  0,      40000,   1000, 2, '2026-03-05 11:45:00+07'),
  ( 72000, 1,     7200,   80000,   8000, 3, '2026-03-06 16:00:00+07'),
  ( 61500, NULL,  0,      70000,   8500, 2, '2026-03-07 10:30:00+07'),
  (185000, 8,    31450,  190000,   5000, 3, '2026-03-07 15:00:00+07'),
  ( 47000, 2,     5000,   50000,   3000, 2, '2026-03-08 09:15:00+07'),
  ( 93000, NULL,  0,     100000,   7000, 3, '2026-03-08 13:45:00+07')
ON CONFLICT DO NOTHING;


-- ============================================================
-- TRANSACTION DETAILS (item per transaksi)
-- ============================================================
-- Kolom sesuai schema DB: transaction_id, product_id, quantity,
-- price, subtotal, harga_beli, discount_type, discount_value, discount_amount
-- (TIDAK ada product_name — nama produk di-JOIN saat query read)

INSERT INTO transaction_details (transaction_id, product_id, quantity, price, subtotal, harga_beli, discount_type, discount_value, discount_amount)
SELECT
  t.id                                                            AS transaction_id,
  p.id                                                            AS product_id,
  d.qty                                                           AS quantity,
  p.harga                                                         AS price,
  p.harga * d.qty                                                 AS subtotal,
  COALESCE(p.harga_beli, ROUND((p.harga * 0.65)::numeric, 0))    AS harga_beli,
  NULL                                                            AS discount_type,
  0                                                               AS discount_value,
  0                                                               AS discount_amount
FROM (
  SELECT id, ROW_NUMBER() OVER (ORDER BY created_at ASC) AS rn
  FROM transactions
  ORDER BY created_at ASC
) t
JOIN LATERAL (
  VALUES
    (1,  1, 3),(1, 11,2),(1, 20,1),
    (2,  4,2),(2, 21,2),(2, 11,3),(2, 24,1),
    (3,  3,2),(3,  6,1),(3, 11,2),(3, 15,1),
    (4, 29,2),(4,  4,3),(4,  1,5),(4, 11,2),(4, 21,2),
    (5,  9,2),(5, 10,1),(5, 24,2),(5, 15,1),
    (6, 30,2),(6,  5,3),(6,  1,4),(6, 11,2),
    (7, 15,1),(7, 13,2),(7, 20,2),
    (8,  4,3),(8, 31,1),(8, 21,2),(8, 11,2),
    (9, 11,4),(9, 15,2),(9,  1,2),
    (10,44,2),(10,43,1),(10,34,2),(10,36,1),(10,45,1),(10,11,5),
    (11, 1,3),(11,20,2),
    (12,32,2),(12, 4,4),(12,21,3),(12,11,2),
    (13,15,2),(13,24,2),(13,13,1),
    (14,22,2),(14, 4,2),(14,11,4),(14, 1,2),
    (15,29,2),(15,32,1),(15, 4,3),(15,21,2),(15,11,3),
    (16,15,2),(16,11,2),
    (17,10,2),(17,22,1),(17, 4,3),(17,24,2),
    (18,11,5),(18, 1,4),(18,13,1),
    (19,30,2),(19,21,2),(19, 4,3),(19,15,1),
    (20,29,3),(20,32,1),(20, 4,5),(20,21,2),(20,11,3),(20,36,1),
    (21,17,1),(21,11,3),(21, 1,2),(21,24,1),
    (22, 4,3),(22,31,1),(22,21,2),(22,11,4),(22, 1,2),
    (23, 9,2),(23,15,1),(23,11,3),(23,24,1),
    (24,44,1),(24,43,1),(24,34,1),(24,35,1),(24,11,5),(24, 4,2),
    (25,15,2),(25,20,2),
    (26, 4,4),(26,32,1),(26,21,2),(26,11,2),
    (27,11,4),(27,24,2),(27,15,1),
    (28,30,1),(28,31,1),(28, 4,3),(28,21,2),(28,11,3),
    (29,15,2),(29,11,2),(29,13,1),
    (30, 3,3),(30,11,2),(30,21,1),(30,20,2),
    (31,44,2),(31,43,1),(31,36,1),(31,34,2),(31,35,1),(31,11,5),(31, 4,5),
    (32,11,4),(32,15,1),(32, 1,2),(32,24,2),
    (33, 4,4),(33,32,1),(33,21,2),(33,11,3),
    (34,15,2),(34,11,2),(34,20,1),
    (35, 9,2),(35, 4,3),(35,11,3),(35,21,1),
    (36,29,2),(36,30,1),(36, 4,3),(36,21,2),(36,11,3),(36,36,1),
    (37,24,2),(37,15,1),(37,11,2),(37, 1,2),
    (38, 4,4),(38,32,1),(38,21,3),(38,11,2),
    (39,11,2),(39,20,2),
    (40,44,1),(40,43,1),(40,35,1),(40,34,1),(40,36,1),(40,11,5),(40, 4,3),
    (41, 4,3),(41,11,3),(41,21,1),(41,24,1),
    (42,32,1),(42, 4,4),(42,11,2),(42, 1,2),(42,21,2),
    (43,15,2),(43,11,3),(43,24,1),
    (44,44,1),(44,43,1),(44,34,2),(44,35,1),(44,11,5),(44, 4,2),
    (45, 1,3),(45,11,2),(45,15,1),
    (46, 4,3),(46,21,2),(46,11,2),(46,24,2),
    (47,30,1),(47, 4,3),(47,11,3),(47,21,1),(47,15,1),
    (48,29,2),(48,30,1),(48, 4,4),(48,21,2),(48,11,3),(48,36,1),(48,34,1),
    (49, 9,2),(49, 4,2),(49,15,1),(49,11,2),
    (50,32,1),(50, 4,4),(50,21,2),(50,11,4),(50, 1,2)
) AS d(rn, pid, qty) ON t.rn = d.rn
JOIN products p ON p.id = d.pid;


-- ============================================================
-- 6. PURCHASES (Pembelian dari Supplier, 10 transaksi)
-- ============================================================
INSERT INTO purchases (supplier_name, total_amount, notes, created_by, created_at) VALUES
  ('PT Indofood Sukses Makmur',  2250000, 'Restok mie instan dan snack bulanan',           1, '2026-01-05 10:00:00+07'),
  ('CV Berkah Jaya Abadi',       1850000, 'Pembelian rokok dan produk kebersihan',          1, '2026-01-12 09:00:00+07'),
  ('Distributor Elektronik Maju', 980000, 'Pengadaan aksesoris elektronik',                1, '2026-01-20 14:00:00+07'),
  ('PT Aqua Golden Mississippi', 1200000, 'Restok minuman: Aqua, Pocari, Teh Botol',       1, '2026-02-03 10:00:00+07'),
  ('PT Unilever Indonesia',      1650000, 'Produk kebersihan: Lifebuoy, Sunsilk, Rinso',   1, '2026-02-10 09:00:00+07'),
  ('CV Frozen Jaya Makmur',       980000, 'Frozen food: nugget, sosis, bakso',             1, '2026-02-17 11:00:00+07'),
  ('Toko Besar Sumber Rejeki',   1750000, 'Restok besar minuman dan makanan',              1, '2026-02-24 10:00:00+07'),
  ('PT Coca-Cola Indonesia',      840000, 'Restok minuman bersoda dan berkarbonasi',       1, '2026-03-02 09:00:00+07'),
  ('Apotek Sehat Farma',          450000, 'Obat-obatan dan vitamin bebas',                 1, '2026-03-05 13:00:00+07'),
  ('Grosir Bumbu Nusantara',      680000, 'Minyak goreng, kecap, dan saus',                1, '2026-03-07 10:30:00+07')
ON CONFLICT DO NOTHING;

INSERT INTO purchase_items (purchase_id, product_id, product_name, quantity, buy_price, subtotal, created_at)
SELECT
  pu.id,
  pr.id,
  pr.nama,
  d.qty,
  pr.harga_beli,
  pr.harga_beli * d.qty,
  pu.created_at
FROM (
  SELECT id, created_at, ROW_NUMBER() OVER (ORDER BY created_at ASC) AS rn
  FROM purchases
  ORDER BY created_at ASC
) pu
JOIN LATERAL (
  VALUES
    (1, 11, 100),(1, 12, 80),(1, 13, 60),(1, 21, 50),(1, 22, 50),(1, 26, 40),
    (2, 29, 30), (2, 30, 30),(2, 31, 20),(2, 33, 50),(2, 34, 30),(2, 37, 20),
    (3, 39, 25), (3, 40, 15),(3, 41, 15),(3, 42, 10),
    (4, 1,  100),(4, 2, 80), (4, 3, 60), (4, 4, 50), (4, 8, 60), (4, 9, 40),
    (5, 33, 50), (5, 34, 30),(5, 35, 25),(5, 36, 20),(5, 37, 20),(5, 38, 15),
    (6, 43, 20), (6, 44, 20),(6, 45, 15),
    (7, 1,  80), (7, 3, 60), (7, 5, 50), (7, 11,80), (7, 13,50),(7, 15,50),
    (8, 5,  50), (8, 6, 50), (8, 7, 40),
    (9, 49, 40), (9, 50, 50),
    (10,46, 25),(10,47, 30),(10,48, 25)
) AS d(rn, pid, qty) ON pu.rn = d.rn
JOIN products pr ON pr.id = d.pid;


-- ============================================================
-- 7. EXPENSES (Pengeluaran 3 bulan, 25 item)
-- ============================================================
INSERT INTO expenses (category, description, amount, expense_date, is_recurring, recurring_period, notes, created_by) VALUES
  -- Januari 2026
  ('Utilitas',     'Tagihan listrik Januari 2026',               440000, '2026-01-05', true,  'monthly', 'PLN - No rek 123456789',           1),
  ('Utilitas',     'Tagihan air Januari 2026',                   118000, '2026-01-05', true,  'monthly', 'PDAM Jaya',                        1),
  ('Utilitas',     'Internet & WiFi Januari 2026',               250000, '2026-01-05', true,  'monthly', 'Telkom Indihome 20Mbps',            1),
  ('Gaji',         'Gaji karyawan kontrak Januari',             1800000, '2026-01-31', false, NULL,      'Karyawan harian 6 orang x 300rb',  1),
  ('Operasional',  'Plastik & kantong belanja',                   75000, '2026-01-10', false, NULL,      '3 pak @ 25rb',                     1),
  ('Operasional',  'Tinta printer struk',                         45000, '2026-01-15', false, NULL,      NULL,                               1),
  ('Iklan',        'Cetak brosur promosi',                       120000, '2026-01-20', false, NULL,      '100 lembar A5',                    1),
  -- Februari 2026
  ('Utilitas',     'Tagihan listrik Februari 2026',              450000, '2026-02-05', true,  'monthly', 'PLN - No rek 123456789',           1),
  ('Utilitas',     'Tagihan air Februari 2026',                  120000, '2026-02-05', true,  'monthly', 'PDAM Jaya',                        1),
  ('Utilitas',     'Internet & WiFi Februari 2026',              250000, '2026-02-05', true,  'monthly', 'Telkom Indihome 20Mbps',            1),
  ('Operasional',  'Service AC toko',                            350000, '2026-02-08', false, NULL,      'Servis dan cuci AC 1 unit',        1),
  ('Operasional',  'Plastik & kantong belanja',                   75000, '2026-02-12', false, NULL,      '3 pak @ 25rb',                     1),
  ('Kebersihan',   'Pembersih lantai & alat kebersihan toko',     90000, '2026-02-15', false, NULL,      NULL,                               1),
  ('Iklan',        'Boost postingan Instagram Februari',         100000, '2026-02-20', false, NULL,      'Campaign 7 hari Valentine week',   1),
  ('Operasional',  'Perbaikan etalase kaca',                     250000, '2026-02-24', false, NULL,      'Retak akibat benturan',            1),
  ('Gaji',         'Gaji karyawan kontrak Februari',            1800000, '2026-02-28', false, NULL,      'Karyawan harian 6 orang x 300rb',  1),
  -- Maret 2026
  ('Utilitas',     'Tagihan listrik Maret 2026',                 460000, '2026-03-05', true,  'monthly', 'PLN - No rek 123456789',           1),
  ('Utilitas',     'Tagihan air Maret 2026',                     122000, '2026-03-05', true,  'monthly', 'PDAM Jaya',                        1),
  ('Utilitas',     'Internet & WiFi Maret 2026',                 250000, '2026-03-05', true,  'monthly', 'Telkom Indihome 20Mbps',            1),
  ('Operasional',  'Plastik & kantong belanja',                   75000, '2026-03-05', false, NULL,      '3 pak @ 25rb',                     1),
  ('Operasional',  'Tinta printer struk',                         45000, '2026-03-06', false, NULL,      NULL,                               1),
  ('Operasional',  'Isi ulang galon air minum karyawan',          24000, '2026-03-07', false, NULL,      '4 galon x 6.000',                  1),
  ('Iklan',        'Boost postingan Instagram Maret',            100000, '2026-03-08', false, NULL,      'Campaign 7 hari promo bulanan',    1),
  ('Pemeliharaan', 'Cat ulang papan nama toko',                  300000, '2026-03-08', false, NULL,      'Cat dan jasa pengecatan',          1),
  ('Lain-lain',    'Konsumsi rapat bulanan',                      85000, '2026-03-08', false, NULL,      'Snack & minuman 10 orang',         1)
ON CONFLICT DO NOTHING;


-- ============================================================
-- 8. PAYROLL (Penggajian 3 bulan x 7 karyawan aktif)
-- ============================================================
INSERT INTO payroll (employee_id, periode, gaji_pokok, bonus, potongan, total, catatan, paid_at, created_by) VALUES
  -- Januari 2026
  (1, 'Januari 2026', 2500000,  200000,      0, 2700000, 'Bonus kinerja baik',             '2026-02-01 10:00:00+07', 1),
  (2, 'Januari 2026', 2500000,       0, 100000, 2400000, 'Potongan izin 1 hari',           '2026-02-01 10:00:00+07', 1),
  (3, 'Januari 2026', 2200000,       0,      0, 2200000, NULL,                              '2026-02-01 10:00:00+07', 1),
  (4, 'Januari 2026', 3500000,  750000,      0, 4250000, 'Bonus supervisor Q4 2025',       '2026-02-01 10:00:00+07', 1),
  (5, 'Januari 2026', 2500000,       0,      0, 2500000, NULL,                              '2026-02-01 10:00:00+07', 1),
  (7, 'Januari 2026', 2000000,       0,  50000, 1950000, 'Potongan keterlambatan 2x',      '2026-02-01 10:00:00+07', 1),
  (8, 'Januari 2026', 1800000,       0,      0, 1800000, NULL,                              '2026-02-01 10:00:00+07', 1),
  -- Februari 2026
  (1, 'Februari 2026', 2500000,      0,      0, 2500000, NULL,                              '2026-03-01 10:00:00+07', 1),
  (2, 'Februari 2026', 2500000, 200000,      0, 2700000, 'Bonus pelayanan terbaik',        '2026-03-01 10:00:00+07', 1),
  (3, 'Februari 2026', 2200000,      0,      0, 2200000, NULL,                              '2026-03-01 10:00:00+07', 1),
  (4, 'Februari 2026', 3500000, 500000,      0, 4000000, 'Bonus supervisor bulanan',       '2026-03-01 10:00:00+07', 1),
  (5, 'Februari 2026', 2500000,      0, 200000, 2300000, 'Potongan keterlambatan 4x',      '2026-03-01 10:00:00+07', 1),
  (7, 'Februari 2026', 2000000,      0,      0, 2000000, NULL,                              '2026-03-01 10:00:00+07', 1),
  (8, 'Februari 2026', 1800000, 100000,      0, 1900000, 'Bonus kerja lembur 2 hari',      '2026-03-01 10:00:00+07', 1),
  -- Maret 2026 (dibayar awal April, tapi kita catat pending)
  (1, 'Maret 2026',   2500000, 300000,      0, 2800000, 'Bonus target penjualan tercapai', '2026-03-31 10:00:00+07', 1),
  (2, 'Maret 2026',   2500000,      0,      0, 2500000, NULL,                               '2026-03-31 10:00:00+07', 1),
  (3, 'Maret 2026',   2200000,      0, 100000, 2100000, 'Potongan sakit tanpa ket.',        '2026-03-31 10:00:00+07', 1),
  (4, 'Maret 2026',   3500000, 500000,      0, 4000000, 'Bonus supervisor bulanan',         '2026-03-31 10:00:00+07', 1),
  (5, 'Maret 2026',   2500000,      0,      0, 2500000, NULL,                               '2026-03-31 10:00:00+07', 1),
  (7, 'Maret 2026',   2000000, 150000,      0, 2150000, 'Bonus lembur 3 malam',             '2026-03-31 10:00:00+07', 1),
  (8, 'Maret 2026',   1800000,      0,      0, 1800000, NULL,                               '2026-03-31 10:00:00+07', 1)
ON CONFLICT DO NOTHING;


-- ============================================================
-- VERIFIKASI DATA
-- ============================================================
SELECT 'users'             AS tabel, COUNT(*) AS jumlah FROM users
UNION ALL SELECT 'categories',       COUNT(*) FROM categories
UNION ALL SELECT 'products',         COUNT(*) FROM products
UNION ALL SELECT 'discounts',        COUNT(*) FROM discounts
UNION ALL SELECT 'employees',        COUNT(*) FROM employees
UNION ALL SELECT 'transactions',     COUNT(*) FROM transactions
UNION ALL SELECT 'transaction_details', COUNT(*) FROM transaction_details
UNION ALL SELECT 'purchases',        COUNT(*) FROM purchases
UNION ALL SELECT 'purchase_items',   COUNT(*) FROM purchase_items
UNION ALL SELECT 'expenses',         COUNT(*) FROM expenses
UNION ALL SELECT 'payroll',          COUNT(*) FROM payroll
ORDER BY tabel;
