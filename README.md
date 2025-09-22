# Transfer Management System (Go Microservices)

Sistem ini terdiri dari beberapa microservices berbasis Golang, Echo framework, database PostgreSQL, dan komunikasi antar service menggunakan gRPC. Sistem ini dirancang untuk mengelola data klub, pemain, dan transfer secara terpisah namun saling terintegrasi.

## Arsitektur Microservices

### 1. Club Service
- Manajemen data klub sepak bola
- Komunikasi gRPC: menerima dan mengirim data klub
- Database: PostgreSQL
- Struktur: entities, repository, service, grpcserver, proto

### 2. Pemain Service
- Manajemen data pemain
- Komunikasi gRPC: menerima dan mengirim data pemain, serta client untuk komunikasi ke Club Service
- Database: PostgreSQL
- Struktur: entities, repository, service, grpcserver, proto, clients

### 3. Transfer Service
- Manajemen data transfer pemain antar klub
- Komunikasi gRPC: menerima dan mengirim data transfer, serta client untuk komunikasi ke Pemain Service
- Database: PostgreSQL
- Struktur: entities, repository, service, grpcserver, proto, clients

Setiap service berjalan secara independen, memiliki entry point (`main.go`), konfigurasi environment, dan port sendiri.

## Teknologi
- Golang
- Echo Framework
- PostgreSQL
- gRPC (Protocol Buffers)

## Struktur Folder

```
microservices/
├── go-echo-Club/
├── go-echo-pemain/
├── go-echo-Transfer/
└── README.md
```

Contoh struktur di masing-masing service:

```
go-echo-Club/
├── main.go
├── config/
├── entities/
├── grpcserver/
├── protobuf/
├── repository/
├── service/
└── ...
```

## Contoh Alur Transfer Pemain

1. Client mengirim request transfer ke Transfer Service (gRPC).
2. Transfer Service validasi pemain ke Pemain Service (gRPC).
3. Transfer Service validasi klub tujuan ke Club Service (gRPC).
4. Jika valid, Transfer Service menyimpan data transfer dan mengirim response sukses ke client.
5. Jika gagal, Transfer Service mengirim response error ke client.

---

## Alur Proses

1. **Club Service** menerima permintaan pembuatan/daftar klub baru.
2. **Pemain Service** menerima permintaan pendaftaran pemain baru, dan dapat melakukan validasi klub melalui Club Service (gRPC).
3. **Transfer Service** menerima permintaan transfer pemain antar klub, melakukan validasi pemain dan klub melalui Pemain Service dan Club Service (gRPC).
4. Semua data disimpan di masing-masing database PostgreSQL sesuai domain.

## Contoh Request & Response API (gRPC)

### 1. Club Service
**Request (gRPC):**
```protobuf
message CreateClubRequest {
	string name = 1;
	string country = 2;
}
```
**Response:**
```protobuf
message ClubResponse {
	int64 id = 1;
	string name = 2;
	string country = 3;
}
```

### 2. Pemain Service
**Request (gRPC):**
```protobuf
message CreatePemainRequest {
	string name = 1;
	int32 age = 2;
	int64 club_id = 3;
}
```
**Response:**
```protobuf
message PemainResponse {
	int64 id = 1;
	string name = 2;
	int32 age = 3;
	int64 club_id = 4;
}
```

### 3. Transfer Service
**Request (gRPC):**
```protobuf
message CreateTransferRequest {
	int64 pemain_id = 1;
	int64 from_club_id = 2;
	int64 to_club_id = 3;
	string transfer_date = 4;
}
```
**Response:**
```protobuf
message TransferResponse {
	int64 id = 1;
	int64 pemain_id = 2;
	int64 from_club_id = 3;
	int64 to_club_id = 4;
	string transfer_date = 5;
	string status = 6;
}
```

## Cara Menjalankan
1. Clone repository dan masuk ke folder masing-masing service
2. Buat file `.env` sesuai kebutuhan (lihat contoh di masing-masing service)
3. Jalankan database PostgreSQL
4. Generate file gRPC dari `.proto` jika belum:
	```bash
	protoc --go_out=. --go-grpc_out=. path/to/your.proto
	```
5. Jalankan service dengan:
	```bash
	go run main.go
	```

## Komunikasi Antar Service
- Antar service berkomunikasi menggunakan gRPC dengan message dan service yang didefinisikan di file `.proto`.
- Setiap service dapat menjadi client dan server gRPC sesuai kebutuhan bisnis.

## Contoh Endpoint gRPC
Misal, untuk transfer pemain:
- Transfer Service menerima request transfer pemain dari Pemain Service melalui gRPC.
- Club Service dapat diakses oleh Pemain Service untuk validasi klub tujuan.

## Catatan
- Setiap service dapat dikembangkan, diuji, dan di-deploy secara independen.
- Pastikan port dan konfigurasi environment tidak bentrok antar service.
- Untuk pengujian, gunakan tool seperti grpcurl atau Postman (gRPC).

---
Sistem ini scalable, maintainable, dan siap dikembangkan lebih lanjut sesuai kebutuhan bisnis transfer pemain sepak bola.
