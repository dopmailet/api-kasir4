-- Verifikasi status purchase #16 setelah perbaikan
SELECT 
    id,
    supplier_name,
    total_amount,
    paid_amount,
    remaining_amount,
    payment_method,
    payment_status
FROM purchases
WHERE id = 16;
