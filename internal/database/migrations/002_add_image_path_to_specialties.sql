-- Add image_path column to specialties table
ALTER TABLE specialties ADD COLUMN IF NOT EXISTS image_path VARCHAR(255);
