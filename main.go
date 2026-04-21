package main

import (
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed embed/*
var staticFiles embed.FS

var protoRegistry *ProtoRegistry

func main() {
	protoRegistry = NewProtoRegistry()

	if err := CheckProtoc(); err != nil {
		fmt.Printf("警告: %v\n", err)
		fmt.Println("提示: 请安装 protoc 以支持 proto 解析功能")
		fmt.Println("安装方式: https://github.com/protocolbuffers/protobuf/releases")
	}

	r := http.NewServeMux()

	r.HandleFunc("GET /api/proto/types", handleProtoTypes)
	r.HandleFunc("POST /api/proto/upload", handleProtoUpload)
	r.HandleFunc("POST /api/proto/upload-directory", handleDirectoryUpload)
	r.HandleFunc("DELETE /api/proto/types/{name}", handleProtoDelete)
	r.HandleFunc("POST /api/proto/decode", handleProtoDecode)

	_, err := fs.Sub(staticFiles, "embed")
	if err != nil {
		r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./embed/proto-debugger.html")
		})
		r.Handle("GET /static/", http.FileServer(http.Dir("./")))
	} else {
		subFS, _ := fs.Sub(staticFiles, "embed")
		r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "embed/proto-debugger.html")
		})
		r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(subFS))))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("服务启动在 :%s\n", port)
	http.ListenAndServe(":"+port, r)
}

type ProtoTypesResponse struct {
	Types []string `json:"types"`
}

func handleProtoTypes(w http.ResponseWriter, r *http.Request) {
	types := protoRegistry.GetLoadedTypes()
	resp := ProtoTypesResponse{Types: types}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type UploadResponse struct {
	Types []string `json:"types,omitempty"`
	Error string   `json:"error,omitempty"`
}

func handleProtoUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(UploadResponse{Error: "读取请求体失败"})
		return
	}

	fileName := r.Header.Get("X-Proto-Name")
	if fileName == "" {
		fileName = "upload.proto"
	}
	if filepath.Ext(fileName) != ".proto" {
		fileName = fileName + ".proto"
	}

	types, err := protoRegistry.LoadProto(body, fileName)
	if err != nil {
		json.NewEncoder(w).Encode(UploadResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(UploadResponse{Types: types})
}

func handleDirectoryUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	types, err := protoRegistry.LoadDirectory(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(UploadResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(UploadResponse{Types: types})
}

func handleProtoDelete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if err := protoRegistry.UnloadByName(name); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"success": "deleted"})
}

type DecodeRequest struct {
	Data     string `json:"data"`
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
}

type DecodeResponse struct {
	JSON  string `json:"json,omitempty"`
	Error string `json:"error,omitempty"`
}

func handleProtoDecode(w http.ResponseWriter, r *http.Request) {
	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(DecodeResponse{Error: "无效的请求格式"})
		return
	}

	if req.Type == "" {
		json.NewEncoder(w).Encode(DecodeResponse{Error: "请选择 proto 类型"})
		return
	}

	if err := CheckProtoc(); err != nil {
		json.NewEncoder(w).Encode(DecodeResponse{Error: "protoc 未安装: " + err.Error()})
		return
	}

	encoding := req.Encoding
	if encoding == "" {
		encoding = "hex"
	}

	var binaryData []byte
	var err error

	switch encoding {
	case "hex":
		binaryData, err = hex.DecodeString(req.Data)
	case "base64":
		binaryData, err = base64.StdEncoding.DecodeString(req.Data)
	default:
		json.NewEncoder(w).Encode(DecodeResponse{Error: "不支持的编码格式，仅支持 hex 或 base64"})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(DecodeResponse{Error: fmt.Sprintf("无效的 %s 数据格式: %v", encoding, err)})
		return
	}

	jsonResult, err := protoRegistry.Decode(binaryData, req.Type)
	if err != nil {
		json.NewEncoder(w).Encode(DecodeResponse{Error: fmt.Sprintf("反序列化失败: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(DecodeResponse{JSON: jsonResult})
}
