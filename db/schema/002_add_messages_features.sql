CREATE TABLE IF NOT EXISTS messages (
  id TEXT PRIMARY KEY,
  sender_name TEXT NOT NULL,
  sender_phone_number TEXT NOT NULL,
  message_text TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
