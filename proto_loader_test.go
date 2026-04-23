package main

import (
	"encoding/hex"
	"testing"
)

func TestExtractGrpcWebPayload(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		wantPayload []byte
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "testdata/1.hex - DATA帧后跟TRAILER帧",
			hexData:     "00000000030a01308000000044687270632d7374617475733a20300d0a782d726571756573742d69643a2031396462356462622d663730302d343030302d383639342d6534306634373636643530300d0a",
			wantPayload: []byte{0x0a, 0x01, 0x30},
			wantErr:     false,
		},
		{
			name:        "仅有DATA帧",
			hexData:     "00000000030a0130",
			wantPayload: []byte{0x0a, 0x01, 0x30},
			wantErr:     false,
		},
		{
			name:    "压缩数据应返回错误",
			hexData: "01000000030a0130",
			wantErr: true,
			errMsg:  "不支持压缩的 gRPC-Web 数据",
		},
		{
			name:    "仅有TRAILER帧（无DATA帧）",
			hexData: "8000000044687270",
			wantErr: true,
			errMsg:  "未找到有效的 gRPC-Web DATA 帧",
		},
		{
			name:        "多个连续DATA帧",
			hexData:     "00000000020a01" + "00000000021020",
			wantPayload: []byte{0x0a, 0x01, 0x10, 0x20},
			wantErr:     false,
		},
		{
			name:        "空DATA帧（length=0）",
			hexData:     "0000000000" + "00000000030a0130",
			wantPayload: []byte{0x0a, 0x01, 0x30},
			wantErr:     false,
		},
		{
			name:    "非gRPC-Web数据（太短）",
			hexData:     "0a01",
			wantErr: true,
			errMsg:  "未找到有效的 gRPC-Web DATA 帧",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hexData)
			if err != nil {
				t.Fatalf("hex.DecodeString failed: %v", err)
			}

			gotPayload, err := extractGrpcWebPayload(data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractGrpcWebPayload() expected error %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("extractGrpcWebPayload() error = %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("extractGrpcWebPayload() unexpected error: %v", err)
				return
			}

			if string(gotPayload) != string(tt.wantPayload) {
				t.Errorf("extractGrpcWebPayload() payload = %x, want %x", gotPayload, tt.wantPayload)
			}
		})
	}
}
