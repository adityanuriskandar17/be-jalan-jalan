# 📠 Implementasi API SMS (GoSMSGateway Masking)

Sistem menggunakan endpoint API Masking dari GoSMSGateway.

## 🔗 Endpoint & Auth

**Base URL:** `https://secure.gosmsgateway.com/masking/api/sendsms.php`
**Authentication:** Basic Auth dengan Username & Password

## 🔧 Konfigurasi

Pastikan file `.env` sudah diupdate:

```bash
# SMS Configuration (Masking API)
SMS_BASE_URL=https://secure.gosmsgateway.com/masking/api/sendsms.php
SMS_USERNAME=asitate
SMS_PASSWORD=rxS5aTkfybV8BZXW # password Anda yang sebenarnya
```

## 💻 Contoh Penggunaan dalam Kode

```go
package main

import "be-jalan/services/notification"

func main() {
    // 1. Dapatkan config
    cfg := notification.GetSMSConfig()

    // 2. Kirim SMS
    // lineType 0 = reguler
    err := notification.SendSMS(cfg, "081234567890", "Halo, ini pesan test", 0)
    if err != nil {
        panic(err)
    }
}
```

## 🛠 Troubleshooting

Jika SMS gagal terkirim:
1. Pastikan **pulsa masking** tersedia di akun GoSMSGateway.
2. Pastikan nomor tujuan **valid** (format internasional tanpa + misal `62812...` atau `0812...` tergantung provider).
3. Cek apakah IP server di whitelist (jika ada batasan IP).
