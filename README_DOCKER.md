# Hướng dẫn chạy với Docker

## Cách 1: Dùng docker-compose (Khuyên dùng)

```bash
# Build và chạy
docker-compose up -d

# Xem logs
docker-compose logs -f

# Dừng
docker-compose down
```

## Cách 2: Dùng docker trực tiếp

```bash
# Build image
docker build -t christmas-full .

# Chạy container
docker run -d \
  --name christmas-app \
  -p 8080:8080 \
  -v $(pwd)/music:/app/music \
  christmas-full

# Xem logs
docker logs -f christmas-app

# Dừng và xóa container
docker stop christmas-app
docker rm christmas-app
```

## Truy cập

Sau khi chạy, mở trình duyệt:
- **Frontend:** http://localhost:8080/1.html
- **API:** http://localhost:8080/api/wishes

## Lưu ý

- Thư mục `music` được mount ra ngoài để file `data.json` được lưu vĩnh viễn
- Nếu muốn thay đổi port, sửa `8080:8080` thành `PORT_KHAC:8080` trong docker-compose.yml

