-- Tabla de proyectos

CREATE TABLE IF NOT EXISTS projects (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  student_email TEXT NOT NULL,
  student_name TEXT NOT NULL,
  tech_stack TEXT NOT NULL DEFAULT '[]',
  categories TEXT NOT NULL DEFAULT '[]',
  stars INTEGER NOT NULL DEFAULT 0,
  repo_url TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
