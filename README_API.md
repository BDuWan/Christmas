# API Server cho Christmas Wishes

## Cài đặt

1. Cài đặt Go: https://golang.org/dl/

2. Cài đặt dependencies:
```bash
go mod download
```

3. Chạy server:
```bash
go run server.go
```

Hoặc build và chạy:
```bash
go build -o server server.go
./server
```

## Sử dụng

1. Mở trình duyệt và truy cập: `http://localhost:8080/1.html`

2. Server sẽ tự động tạo file `music/data.json` nếu chưa có

## API Endpoints

### GET /api/wishes
Lấy tất cả lời chúc

**Response:**
```json
[
  {
    "id": 1,
    "sender": "Anh",
    "content": "Chúc em một Giáng Sinh ấm áp và hạnh phúc! 🎄✨"
  }
]
```

### POST /api/wishes
Thêm lời chúc mới

**Request Body:**
```json
{
  "sender": "Tên người gửi",
  "content": "Nội dung lời chúc"
}
```

**Response:**
```json
{
  "id": 5,
  "sender": "Tên người gửi",
  "content": "Nội dung lời chúc"
}
```

### PUT /api/wishes
Cập nhật toàn bộ danh sách lời chúc

**Request Body:**
```json
[
  {
    "id": 1,
    "sender": "Anh",
    "content": "Lời chúc mới"
  }
]
```

### DELETE /api/wishes/{id}
Xóa lời chúc theo ID

**Response:** 204 No Content

## Cấu trúc file

```
.
├── server.go          # Server API
├── go.mod             # Go module
├── music/
│   └── data.json     # File lưu lời chúc (tự động tạo)
└── 1.html            # File HTML chính
```

