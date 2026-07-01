-- 1. Revertir cambios en SAMPLES
ALTER TABLE samples DROP CONSTRAINT IF EXISTS check_samples_status;
ALTER TABLE samples ALTER COLUMN status SET DEFAULT 'processing'; -- O el default que tuvieras antes

-- 2. Revertir cambios en NOTIFICATIONS
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS check_notifications_status;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS check_notifications_type;
ALTER TABLE notifications ALTER COLUMN status SET DEFAULT 'unread';