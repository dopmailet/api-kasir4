-- Find the TOKO MOHON BERSANDAR purchase that has partial payment
-- Then directly update it and its linked supplier_payable to 'paid'

-- 1. First check the state
SELECT p.id, p.supplier_name, p.total_amount, p.paid_amount, p.remaining_amount, p.payment_status,
       sp.id AS sp_id, sp.amount AS sp_amount, sp.paid_amount AS sp_paid, sp.status AS sp_status
FROM purchases p
LEFT JOIN supplier_payables sp ON sp.purchase_id = p.id
WHERE p.payment_status = 'partial'
ORDER BY p.created_at DESC;
