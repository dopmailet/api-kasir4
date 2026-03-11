-- Fix: Sync supplier_payables.paid_amount agar sesuai dengan purchases.paid_amount
-- untuk semua payable yang terhubung ke purchase

-- 1. Update supplier_payables paid_amount agar match dengan purchases paid_amount
UPDATE supplier_payables sp
SET 
    paid_amount = p.paid_amount,
    status = CASE 
        WHEN p.paid_amount >= sp.amount THEN 'paid'
        WHEN p.paid_amount > 0 THEN 'partial'
        ELSE 'unpaid'
    END,
    updated_at = NOW()
FROM purchases p
WHERE sp.purchase_id = p.id
  AND sp.paid_amount != p.paid_amount;

-- 2. Recalculate total_payable di semua suppliers
UPDATE suppliers s
SET total_payable = COALESCE((
    SELECT SUM(sp.amount - sp.paid_amount)
    FROM supplier_payables sp
    WHERE sp.supplier_id = s.id AND sp.status != 'paid'
), 0),
updated_at = NOW();
