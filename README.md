# gRPC Polyglot Lab

Contoh praktikum komunikasi antarlayanan menggunakan **gRPC** dengan
bahasa pemrograman yang berbeda:

-   **Go** sebagai gRPC server.
-   **Python** sebagai client.
-   **Node.js** sebagai client.
-   **Protocol Buffers (`user.proto`)** sebagai definisi pesan dan
    layanan.

Server menyediakan layanan `GetUserProfile`. Client Python dan Node.js
mengirim permintaan dengan `user_id` dan menerima nama serta email dari
server.

## 1. Struktur proyek

Setelah semua langkah selesai, struktur proyek utama kurang lebih
seperti berikut:

``` text
grpc-polyglot-lab/
├── .gitignore
├── user.proto
├── goserver/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── pb/
│       ├── user.pb.go
│       └── user_grpc.pb.go
├── python-client/
│   ├── .venv/                 # dibuat secara lokal, tidak di-commit
│   ├── client.py
│   └── pb/
│       ├── user_pb2.py
│       └── user_pb2_grpc.py
└── node-client/
    ├── node_modules/           # dibuat melalui npm install
    ├── client.js
    ├── package.json
    └── package-lock.json
```

## 2. Prasyarat

Pasang perangkat berikut terlebih dahulu:

-   Git
-   Go
-   Python 3
-   Node.js dan npm
-   Protocol Buffers compiler (`protoc`)

Periksa instalasi dari Git Bash:

``` bash
git --version
go version
py --version
node --version
npm --version
protoc --version
```

Pastikan perintah-perintah tersebut dikenali. Di Windows, gunakan
`py -3` atau versi Python 3 yang terpasang untuk membuat virtual
environment jika `python` menunjuk ke Python yang berbeda.

## 3. Siapkan workspace dan Git

Buat folder proyek dan inisialisasi repository:

``` bash
mkdir grpc-polyglot-lab
cd grpc-polyglot-lab
git init
```

Buat file `.gitignore` di folder utama dengan isi:

``` gitignore
node_modules/
__pycache__/
.venv/
*.exe
```

Tambahkan dan simpan file tersebut:

``` bash
git add .gitignore
git commit -m "Initial commit: Add gitignore rules"
```

> Jika Git meminta identitas, atur nama dan email Git terlebih dahulu
> menggunakan `git config --global user.name "Nama Anda"` dan
> `git config --global user.email "email@example.com"`.

## 4. Buat definisi layanan `user.proto`

Buat file `user.proto` di folder utama proyek:

``` proto
syntax = "proto3";

package user;

option go_package = "./pb";

message UserRequest {
    int32 user_id = 1;
}

message UserResponse {
    string name = 1;
    string email = 2;
}

service UserService {
    rpc GetUserProfile(UserRequest) returns (UserResponse);
}
```

File ini mendefinisikan permintaan `UserRequest`, respons
`UserResponse`, dan RPC `GetUserProfile`.

## 5. Implementasikan server Go

### 5.1 Buat modul Go dan pasang dependensi

Dari folder utama proyek, jalankan:

``` bash
mkdir goserver
cd goserver
go mod init goserver
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

### 5.2 Pasang generator Go

``` bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Pastikan folder bin Go dapat ditemukan oleh shell. Pada Windows,
generator biasanya berada di `$(go env GOPATH)/bin`. Jika `protoc` tidak
menemukan generator, tambahkan folder tersebut ke `PATH`, lalu buka
ulang terminal.

### 5.3 Generate kode Go

Pastikan posisi terminal masih di folder `goserver`. Jalankan:

``` bash
mkdir -p pb
protoc --proto_path=.. --go_out=. --go-grpc_out=. ../user.proto
```

Periksa hasilnya:

``` bash
ls pb
```

Seharusnya ada `user.pb.go` dan `user_grpc.pb.go`.

### 5.4 Buat `main.go`

Buat file `goserver/main.go` dengan isi:

``` go
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "goserver/pb"
)

type server struct {
    pb.UnimplementedUserServiceServer
}

func (s *server) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
    log.Printf("Menerima permintaan untuk user ID: %d", req.GetUserId())

    return &pb.UserResponse{
        Name:  "Maftuh",
        Email: "maftuh@example.com",
    }, nil
}

func main() {
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Gagal membuka port: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, &server{})

    log.Println("gRPC server berjalan pada port 50051")

    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Gagal menjalankan server: %v", err)
    }
}
```

Nama dan email di atas adalah data contoh untuk praktikum; silakan ubah
sesuai kebutuhan.

### 5.5 Jalankan server

Dari folder `goserver`:

``` bash
go run .
```

Output yang diharapkan:

``` text
gRPC server berjalan pada port 50051
```

Biarkan terminal server tetap terbuka saat menguji client.

## 6. Buat client Python

Buka terminal Git Bash baru. Dari folder utama proyek, jalankan:

``` bash
cd /d/grpc-polyglot-lab/python-client
```

Jika belum ada folder `python-client`, buat dari folder utama terlebih
dahulu:

``` bash
cd /d/grpc-polyglot-lab
mkdir python-client
cd python-client
```

### 6.1 Buat dan aktifkan virtual environment

Di Windows dengan Git Bash:

``` bash
py -3 -m venv .venv
source .venv/Scripts/activate
```

Jika `py -3` tidak tersedia, gunakan launcher Python 3 yang terpasang di
komputer. Setelah aktif, prompt biasanya menampilkan `(.venv)`.

### 6.2 Pasang paket Python

``` bash
python -m pip install grpcio grpcio-tools
```

### 6.3 Generate kode Python

Dari folder `python-client`:

``` bash
mkdir -p pb
python -m grpc_tools.protoc -I.. --python_out=pb --grpc_python_out=pb ../user.proto
```

Periksa hasil:

``` bash
ls pb
```

Seharusnya muncul `user_pb2.py` dan `user_pb2_grpc.py`.

### 6.4 Buat `client.py`

Buat file `python-client/client.py`:

``` python
import grpc
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "pb"))

import user_pb2
import user_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)
        request = user_pb2.UserRequest(user_id=1)
        response = stub.GetUserProfile(request)

        print("=== HASIL DATA USER ===")
        print(f"Nama  : {response.name}")
        print(f"Email : {response.email}")


if __name__ == "__main__":
    run()
```

### 6.5 Jalankan client Python

Pastikan server Go masih berjalan. Di terminal `python-client` dengan
virtual environment aktif:

``` bash
python client.py
```

Output yang diharapkan:

``` text
=== HASIL DATA USER ===
Nama  : Maftuh
Email : maftuh@example.com
```

## 7. Buat client Node.js

Buka terminal baru. Dari folder utama proyek:

``` bash
cd /d/grpc-polyglot-lab
mkdir node-client
cd node-client
npm init -y
npm install @grpc/grpc-js @grpc/proto-loader
```

Jika folder `node-client` sudah ada, cukup masuk ke folder itu dan
jalankan `npm install`.

### 7.1 Buat `client.js`

Buat file `node-client/client.js`:

``` javascript
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path = require("path");

const PROTO_PATH = path.join(__dirname, "..", "user.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const userPackage = protoDescriptor.user;

const client = new userPackage.UserService(
    "localhost:50051",
    grpc.credentials.createInsecure()
);

client.GetUserProfile({ userId: 1 }, (error, response) => {
    if (error) {
        console.error("Gagal mengambil data user:", error.message);
        client.close();
        process.exitCode = 1;
        return;
    }

    console.log("=== HASIL DATA USER ===");
    console.log("Nama  :", response.name);
    console.log("Email :", response.email);

    client.close();
});
```

### 7.2 Jalankan client Node.js

Pastikan server Go masih berjalan, lalu dari folder `node-client`:

``` bash
node client.js
```

Output yang diharapkan:

``` text
=== HASIL DATA USER ===
Nama  : Maftuh
Email : maftuh@example.com
```

## 8. Alur pengujian

Jalankan server Go pada satu terminal:

``` bash
cd /d/grpc-polyglot-lab/goserver
go run .
```

Kemudian, dari terminal lain, jalankan client Python:

``` bash
cd /d/grpc-polyglot-lab/python-client
source .venv/Scripts/activate
python client.py
```

Jalankan client Node.js dari terminal lain:

``` bash
cd /d/grpc-polyglot-lab/node-client
node client.js
```

Kedua client harus menampilkan data user yang sama. Server Go akan
mencatat permintaan yang masuk di terminal server.

## 9. Simpan perubahan ke Git

Dari folder utama proyek:

``` bash
cd /d/grpc-polyglot-lab
git status
git add .
git commit -m "Implement gRPC polyglot communication"
git status
```

Jika semua perubahan sudah di-commit, status akhir akan menunjukkan:

``` text
nothing to commit, working tree clean
```

Jangan menjalankan commit kedua dengan pesan yang sama jika tidak ada
perubahan baru. Jika Git mengatakan `nothing to commit`, itu berarti
tidak ada perubahan yang perlu disimpan.

## 10. Troubleshooting

### `protoc: command not found`

Pastikan Protocol Buffers compiler terpasang dan perintah
`protoc --version` bisa dijalankan dari terminal.

### `protoc-gen-go: program not found`

Pastikan generator sudah dipasang dan folder bin Go ditambahkan ke
`PATH`. Tutup dan buka ulang terminal setelah mengubah `PATH`.

### `connection refused` atau client gagal terhubung

Pastikan server Go sedang berjalan dan mendengarkan port `50051`.
Jalankan client di terminal lain; jangan hentikan terminal server.

### Modul Python tidak ditemukan

Aktifkan virtual environment:

``` bash
source .venv/Scripts/activate
```

Kemudian pastikan dependensi terpasang:

``` bash
python -m pip install grpcio grpcio-tools
```

### Paket Node.js tidak ditemukan

Dari folder `node-client`, jalankan:

``` bash
npm install
```

### Git menampilkan peringatan LF/CRLF

Peringatan ini berkaitan dengan konversi akhir baris pada Windows dan
tidak otomatis berarti proses gagal. Periksa status serta hasil perintah
untuk memastikan.

## 11. Ringkasan

Praktikum ini menunjukkan komunikasi gRPC lintas bahasa: server Go
menyediakan RPC `GetUserProfile`, sementara client Python dan Node.js
memanggil layanan yang sama menggunakan definisi bersama `user.proto`.

**Catatan keamanan:** contoh ini menggunakan koneksi tanpa TLS
(`insecure`) dan data user statis. Gunakan hanya untuk praktikum lokal,
bukan sebagai konfigurasi produksi.
