-- Add payment columns to purchases table
ALTER TABLE purchases 
ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50) DEFAULT 'cash',
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) DEFAULT 'paid',
ADD COLUMN IF NOT EXISTS paid_amount NUMERIC(15,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS remaining_amount NUMERIC(15,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS due_date DATE,
ADD COLUMN IF NOT EXISTS payment_notes TEXT;

-- Update existing records to reflect them as cash purchases fully paid
UPDATE purchases
SET 
    paid_amount = total_amount,
    payment_method = 'cash',
    payment_status = 'paid',
    remaining_amount = 0
WHERE paid_amount = 0 AND total_amount > 0 AND payment_method IS NULL;
