-- Fix SEMUA purchases yang statusnya belum 'paid' tapi supplier_payables-nya sudah lunas
UPDATE purchases p
SET 
    paid_amount = p.total_amount,
    remaining_amount = 0,
    payment_status = 'paid'
FROM supplier_payables sp
WHERE sp.purchase_id = p.id
  AND sp.status = 'paid'
  AND p.payment_status != 'paid';

-- Fix purchase yang payment_method = 'cash' tapi statusnya masih salah
UPDATE purchases
SET 
    paid_amount = total_amount,
    remaining_amount = 0,
    payment_status = 'paid'
WHERE payment_method = 'cash' AND payment_status != 'paid';
