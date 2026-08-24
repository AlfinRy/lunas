-- Lunas schema. Idempotent — applied at startup.
CREATE TABLE IF NOT EXISTS clients (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  name               TEXT    NOT NULL,
  email              TEXT    NOT NULL,
  payment_terms_days INTEGER NOT NULL DEFAULT 14,
  relationship_note  TEXT    NOT NULL DEFAULT '',
  created_at         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS invoices (
  id                INTEGER PRIMARY KEY AUTOINCREMENT,
  client_id         INTEGER NOT NULL REFERENCES clients(id),
  number            TEXT    NOT NULL,
  amount_cents      INTEGER NOT NULL,
  amount_paid_cents INTEGER NOT NULL DEFAULT 0,
  currency          TEXT    NOT NULL DEFAULT 'USD',
  issued_on         TEXT    NOT NULL,
  due_on            TEXT    NOT NULL,
  status            TEXT    NOT NULL DEFAULT 'scheduled',
  notes             TEXT    NOT NULL DEFAULT '',
  agent_state       TEXT    NOT NULL DEFAULT 'idle',
  current_stage     TEXT,
  next_action_on    TEXT,
  created_at        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  UNIQUE (client_id, number)
);

CREATE TABLE IF NOT EXISTS email_drafts (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  invoice_id INTEGER NOT NULL REFERENCES invoices(id),
  stage      TEXT    NOT NULL,
  subject    TEXT    NOT NULL,
  body       TEXT    NOT NULL,
  status     TEXT    NOT NULL DEFAULT 'pending',
  created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  sent_at    TEXT
);

CREATE TABLE IF NOT EXISTS outbox_emails (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  invoice_id     INTEGER NOT NULL REFERENCES invoices(id),
  draft_id       INTEGER REFERENCES email_drafts(id),
  invoice_number TEXT    NOT NULL,
  to_name        TEXT    NOT NULL,
  to_email       TEXT    NOT NULL,
  subject        TEXT    NOT NULL,
  body           TEXT    NOT NULL,
  sent_at        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS payments (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  invoice_id   INTEGER REFERENCES invoices(id),
  amount_cents INTEGER NOT NULL,
  paid_on      TEXT    NOT NULL,
  source       TEXT    NOT NULL DEFAULT 'manual',
  confidence   TEXT,
  raw_text     TEXT,
  created_at   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS activities (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  invoice_id INTEGER REFERENCES invoices(id),
  type       TEXT    NOT NULL,
  message    TEXT    NOT NULL,
  created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS settings (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  sender_name        TEXT NOT NULL DEFAULT 'Rani Prameswari',
  sender_email       TEXT NOT NULL DEFAULT 'billing@ranistudio.co',
  default_terms_days INTEGER NOT NULL DEFAULT 14,
  global_mode        TEXT NOT NULL DEFAULT 'approve_each',
  sim_now            TEXT
);

CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_client ON invoices(client_id);
CREATE INDEX IF NOT EXISTS idx_drafts_invoice ON email_drafts(invoice_id);
CREATE INDEX IF NOT EXISTS idx_activities_invoice ON activities(invoice_id);
