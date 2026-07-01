-- 1. Modificar la tabla de SAMPLES (para soportar los nuevos estados)
ALTER TABLE samples 
  ALTER COLUMN status TYPE VARCHAR(20),
  ALTER COLUMN status SET DEFAULT 'processing',
  ADD CONSTRAINT check_samples_status CHECK (status IN ('starting', 'processing', 'succeeded', 'failed', 'canceled'));

-- 2. Modificar la tabla de NOTIFICATIONS (para soportar el NotificationStatus y TypeNotification)
ALTER TABLE notifications 
  ALTER COLUMN status TYPE VARCHAR(20),
  ALTER COLUMN status SET DEFAULT 'unread',
  ALTER COLUMN type TYPE VARCHAR(50),
  ADD CONSTRAINT check_notifications_status CHECK (status IN ('unread', 'read')),
  ADD CONSTRAINT check_notifications_type CHECK (type IN ('replicate_error', 'replicate_success', 'payment', 'info'));