-- Fix purchase #16: recalculate paid_amount, remaining_amount, and payment_status
-- berdasarkan total pembayaran nyata dari payable_payments + paid_amount awal

WITH calc AS (
    SELECT
        p.id AS purchase_id,
        p.total_amount,
        -- paid_amount awal (DP saat beli) + cicilan dari payable_payments
        p.paid_amount + COALESCE((
            SELECT SUM(pp.amount)
            FROM payable_payments pp
            JOIN supplier_payables sp ON pp.payable_id = sp.id
            WHERE sp.purchase_id = p.id
        ), 0) AS total_paid_now
    FROM purchases p
    WHERE p.id = 16
)
UPDATE purchases
SET
    paid_amount      = LEAST(calc.total_paid_now, calc.total_amount),
    remaining_amount = GREATEST(calc.total_amount - calc.total_paid_now, 0),
    payment_status   = CASE
                           WHEN calc.total_paid_now >= calc.total_amount THEN 'paid'
                           WHEN calc.total_paid_now > 0                  THEN 'partial'
                           ELSE 'unpaid'
                       END
FROM calc
WHERE purchases.id = calc.purchase_id;
