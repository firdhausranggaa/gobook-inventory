# API RESTful Inventaris Buku (Golang)

API RESTful tingkat produksi yang tangguh untuk manajemen perpustakaan dan inventaris. Dibangun menggunakan Go (Golang), proyek ini mendemonstrasikan arsitektur *backend* skala *enterprise*, kontrol konkurensi tingkat lanjut, dan mekanisme otorisasi yang aman.

## 🚀 Fitur Skala Enterprise

*   **Arsitektur RESTful:** *Backend* sepenuhnya *headless* yang merespons dengan JSON terstruktur dan kode status HTTP yang ketat.
*   **Kontrol Konkurensi:** Mengimplementasikan **Pessimistic Locking** (`FOR UPDATE`) dan **Database Transactions** untuk mencegah bentrok data (*race condition*) selama proses peminjaman dan pengembalian buku.
*   **Optimasi Sumber Daya:** Konfigurasi **Connection Pooling** PostgreSQL (batas koneksi aktif/menganggur) untuk menangani lalu lintas tinggi secara efisien.
*   **Penghentian Aman (Graceful Shutdown):** Server secara aman menyelesaikan transaksi *database* yang sedang berjalan sebelum benar-benar mati saat menerima sinyal interupsi.
*   **Rekam Jejak & Soft Delete:** Terintegrasi dengan fitur *Soft Delete* GORM (`DeletedAt`) untuk keamanan audit data tanpa penghapusan permanen.
*   **Autentikasi Aman:** Pendaftaran dan *login* pengguna memanfaatkan **Bcrypt** untuk enkripsi kata sandi dan **JWT** untuk autentikasi *stateless*[cite: 26].
*   **Kontrol Akses Berbasis Peran (RBAC):** Pemisahan tingkat izin yang jelas antara `Admin` (akses penuh CRUD) dan `Member` (akses baca dan pinjam)[cite: 28].

## 🛠️ Teknologi yang Digunakan

*   **Bahasa Pemrograman:** Go (Golang)
*   **Framework:** Gin-Gonic (dilengkapi Middleware CORS)
*   **ORM:** GORM[cite: 24, 27]
*   **Database:** PostgreSQL[cite: 27]
*   **Keamanan:** Golang-JWT (v4), X/Crypto (Bcrypt)[cite: 26]

## ⚙️ Persyaratan Sistem

*   Go terinstal di komputer Anda.
*   Server PostgreSQL yang berjalan secara lokal atau *remote*.

## 📦 Instalasi & Konfigurasi Lokal

1. **Kloning repositori:**
```bash
git clone [https://github.com/firdhausranggaa/gobook-inventory.git](https://github.com/firdhausranggaa/gobook-inventory.git)
cd gobook-inventory

```

2. **Variabel Lingkungan:**
Buat *file* `.env` di direktori utama dan definisikan variabel berikut:
```env
POSTGRES_URL="host=localhost user=postgres password=password_anda dbname=database_anda port=5432 sslmode=disable"
SUPER_USER="admin"
SUPER_PASS="123"
SUPER_SECRET="kunci-rahasia-jwt-anda"
APP_PORT="8080"

```


3. **Unduh Dependensi:**
```bash
go mod tidy

```


4. **Jalankan Server:**
```bash
go run main.go

```


*Catatan: Pada saat pertama kali dijalankan, API akan secara otomatis melakukan migrasi tabel, membuat Foreign Key, dan mengisi data awal (seeder) untuk admin serta katalog buku.*

## 📡 Referensi Endpoint API

Semua rute yang dilindungi (*protected routes*) mewajibkan *header* `Authorization` dengan format: `Bearer <token_anda>`.

### Autentikasi (Publik)

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `POST` | `/api/register` | Mendaftarkan akun *member* baru |
| `POST` | `/api/login` | Autentikasi pengguna dan menerima token JWT |

### Manajemen Buku (Terlindungi)

| Method | Endpoint | Peran Akses | Deskripsi |
| --- | --- | --- | --- |
| `GET` | `/api/books` | Admin / Member | Menampilkan semua buku. Mendukung parameter `?search=`, `?page=`, `?limit=`<br> |
| `GET` | `/api/books/:id` | Admin / Member | Menampilkan detail spesifik dari satu buku

 |
| `POST` | `/api/books` | **Khusus Admin** | Menambahkan buku baru ke inventaris

 |
| `PUT` | `/api/books/:id` | **Khusus Admin** | Memperbarui detail buku yang sudah ada

 |
| `DELETE` | `/api/books/:id` | **Khusus Admin** | Menghapus buku (Menggunakan sistem Audit/Soft Delete)

 |

### Sistem Peminjaman (Terlindungi & Transaksional)

| Method | Endpoint | Peran Akses | Deskripsi |
| --- | --- | --- | --- |
| `GET` | `/api/borrowings/me` | Admin / Member | Menarik riwayat peminjaman aktif/lampau milik pengguna (Otomatis memuat *Foreign Key*) |
| `POST` | `/api/borrow` | Admin / Member | Meminjam buku (Menerapkan penguncian *Pessimistic Locking*) |
| `POST` | `/api/return/:id` | Admin / Member | Mengembalikan buku yang dipinjam berdasarkan ID Transaksi |
