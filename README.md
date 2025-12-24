# Christmas Tree Project

Dự án cây thông Giáng Sinh tương tác với cử chỉ tay.

## Cấu trúc thư mục

```
/
├── index.html          # File chính
├── _redirects          # Cloudflare Pages redirects
├── img/                # Thư mục ảnh
│   ├── a1.jpg
│   ├── a2.jpg
│   ├── a3.jpg
│   └── a4.jpg
└── music/              # Thư mục nhạc
    ├── song1.mp3
    ├── song2.mp3
    └── song3.mp3
```

## Deploy lên Cloudflare Pages

1. Đăng nhập vào Cloudflare Dashboard
2. Vào Pages → Create a project
3. Kết nối với Git repository hoặc upload thư mục
4. Build command: (để trống)
5. Build output directory: `/`
6. Deploy!

## Tính năng

- Tương tác bằng cử chỉ tay (MediaPipe Hands)
- 5 ngón: Mở cây thông
- 2 ngón: Hiển thị chữ 3D
- 3 ngón: Hiển thị trái tim
- Phát nhạc tự động
- Responsive cho mobile

