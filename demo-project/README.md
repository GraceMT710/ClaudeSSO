# Todo App — React + Golang (MVP demo)

Contoh project full-stack minimal: React (Vite) di frontend, Go (stdlib `net/http`,
tanpa framework/dependency eksternal) di backend. Data disimpan in-memory —
cukup untuk belajar alur, tinggal ganti `store` di `main.go` kalau mau
persistence (Postgres/Firestore).

## Struktur

```
demo-project/
├── backend/       # Go REST API (port 8080)
│   ├── go.mod
│   └── main.go
└── frontend/      # React + Vite (port 5173)
    └── src/
        ├── App.jsx
        └── App.css
```

## Menjalankan

**1. Backend**
```bash
cd backend
go run main.go
# -> backend running on :8080
```

**2. Frontend** (terminal terpisah)
```bash
cd frontend
npm install
npm run dev
# -> buka http://localhost:5173
```

## API

| Method | Path              | Body                  | Keterangan       |
|--------|-------------------|-----------------------|-------------------|
| GET    | /api/todos        | -                     | List semua todo   |
| POST   | /api/todos        | `{"title": "..."}`    | Tambah todo       |
| PATCH  | /api/todos/{id}   | -                     | Toggle done       |
| DELETE | /api/todos/{id}   | -                     | Hapus todo        |
