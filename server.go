package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	port        = ":8080"
	dataFile    = "./music/data.json"
	defaultData = `[
  {
    "id": 1,
    "sender": "Anh",
    "content": "Chúc em một Giáng Sinh ấm áp và hạnh phúc! 🎄✨"
  },
  {
    "id": 2,
    "sender": "Mẹ",
    "content": "Con yêu, chúc con luôn vui vẻ và thành công trong cuộc sống! ❤️"
  },
  {
    "id": 3,
    "sender": "Bạn thân",
    "content": "Merry Christmas! Chúc bạn năm mới nhiều may mắn và hạnh phúc! 🎁"
  },
  {
    "id": 4,
    "sender": "Gia đình",
    "content": "Chúc cả nhà một mùa Giáng Sinh an lành và đầy ắp tiếng cười! 🎅"
  }
]`
)

type Wish struct {
	ID      int    `json:"id"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Đọc file JSON
func readWishes() ([]Wish, error) {
	// Tạo file mặc định nếu chưa tồn tại
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		// Tạo thư mục nếu chưa có
		dir := filepath.Dir(dataFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		// Ghi file mặc định
		if err := ioutil.WriteFile(dataFile, []byte(defaultData), 0644); err != nil {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, fmt.Errorf("không đọc được file: %v", err)
	}

	// Kiểm tra file rỗng
	if len(data) == 0 {
		// File rỗng, tạo lại với dữ liệu mặc định
		if err := ioutil.WriteFile(dataFile, []byte(defaultData), 0644); err != nil {
			return nil, fmt.Errorf("không tạo được file mặc định: %v", err)
		}
		data = []byte(defaultData)
	}

	var wishes []Wish
	if err := json.Unmarshal(data, &wishes); err != nil {
		// Nếu parse lỗi, tạo lại file với dữ liệu mặc định
		log.Printf("Lỗi parse JSON: %v, tạo lại file mặc định", err)
		if err := ioutil.WriteFile(dataFile, []byte(defaultData), 0644); err != nil {
			return nil, fmt.Errorf("không tạo được file mặc định: %v", err)
		}
		// Parse lại dữ liệu mặc định
		if err := json.Unmarshal([]byte(defaultData), &wishes); err != nil {
			return nil, fmt.Errorf("lỗi parse dữ liệu mặc định: %v", err)
		}
	}

	return wishes, nil
}

// Ghi file JSON
func writeWishes(wishes []Wish) error {
	data, err := json.MarshalIndent(wishes, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dataFile, data, 0644)
}

// GET /api/wishes - Lấy tất cả lời chúc
func getWishes(w http.ResponseWriter, r *http.Request) {
	wishes, err := readWishes()
	if err != nil {
		log.Printf("Lỗi khi đọc wishes: %v", err)
		// Trả về mảng rỗng thay vì lỗi để frontend không crash
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	// Đảm bảo luôn trả về array, không phải null
	if wishes == nil {
		wishes = []Wish{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(wishes); err != nil {
		log.Printf("Lỗi khi encode JSON: %v", err)
	}
}

// POST /api/wishes - Thêm lời chúc mới
func addWish(w http.ResponseWriter, r *http.Request) {
	var newWish Wish
	if err := json.NewDecoder(r.Body).Decode(&newWish); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate
	if strings.TrimSpace(newWish.Sender) == "" || strings.TrimSpace(newWish.Content) == "" {
		http.Error(w, "Sender và Content không được để trống", http.StatusBadRequest)
		return
	}

	wishes, err := readWishes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tìm ID lớn nhất và tăng lên 1
	maxID := 0
	for _, wish := range wishes {
		if wish.ID > maxID {
			maxID = wish.ID
		}
	}
	newWish.ID = maxID + 1

	wishes = append(wishes, newWish)

	if err := writeWishes(wishes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWish)
}

// PUT /api/wishes - Cập nhật toàn bộ danh sách
func updateWishes(w http.ResponseWriter, r *http.Request) {
	var wishes []Wish
	if err := json.NewDecoder(r.Body).Decode(&wishes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := writeWishes(wishes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishes)
}

// DELETE /api/wishes/{id} - Xóa lời chúc
func deleteWish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	wishes, err := readWishes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tìm và xóa
	found := false
	for i, wish := range wishes {
		if wish.ID == id {
			wishes = append(wishes[:i], wishes[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Không tìm thấy lời chúc", http.StatusNotFound)
		return
	}

	if err := writeWishes(wishes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/wishes", getWishes).Methods("GET")
	api.HandleFunc("/wishes", addWish).Methods("POST")
	api.HandleFunc("/wishes", updateWishes).Methods("PUT")
	api.HandleFunc("/wishes/{id:[0-9]+}", deleteWish).Methods("DELETE")

	// Serve static files (HTML, CSS, JS)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	fmt.Printf("🚀 Server đang chạy tại http://localhost%s\n", port)
	fmt.Printf("📝 API: http://localhost%s/api/wishes\n", port)
	fmt.Printf("📄 Mở file: http://localhost%s/1.html\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

