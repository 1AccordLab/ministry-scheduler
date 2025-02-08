CREATE TABLE USERS (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  phone_number TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  -- OAuth2 related fields
  line_id TEXT UNIQUE NOT NULL,
  line_name TEXT NOT NULL,
  avatar_url TEXT NOT NULL,
  status_message TEXT NOT NULL
)
