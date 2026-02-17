# 📮 Panduan Import ke Postman

Panduan lengkap untuk menggunakan API ASITA - Jalan Jalan di Postman.

---

## 🚀 Quick Start

### Option 1: Import Collection & Environment (Recommended)

1. **Buka Postman**

2. **Import Collection:**
   - Klik tombol **"Import"** di kiri atas
   - Pilih file: `ASITA_Jalan_Jalan_API.postman_collection.json`
   - Klik **"Import"**

3. **Import Environment:**
   - Klik tombol **"Import"** lagi
   - Pilih file: `ASITA_Jalan_Jalan.postman_environment.json`
   - Klik **"Import"**

4. **Aktifkan Environment:**
   - Di kanan atas, pilih dropdown environment
   - Pilih **"ASITA - Jalan Jalan Environment"**

5. **Selesai!** Anda siap untuk testing.

---

### Option 2: Manual Setup

Jika Anda ingin setup manual, ikuti langkah berikut:

#### 1. Buat Collection Baru

- Klik **"New"** → **"Collection"**
- Nama: `ASITA - Jalan Jalan API`
- Description: `API Collection untuk ASITA Jalan Jalan Backend`

#### 2. Buat Environment Baru

- Klik **"Environments"** di sidebar kiri
- Klik **"+"** untuk create new environment
- Nama: `ASITA - Jalan Jalan Environment`

#### 3. Set Environment Variables

Tambahkan variables berikut:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| `base_url` | `http://localhost:8080/api/v1` | `http://localhost:8080/api/v1` |
| `fullName` | `Test User` | `Test User` |
| `email` | `test@example.com` | `test@example.com` |
| `phoneNumber` | `081234567890` | `081234567890` |
| `password` | `password123` | `password123` |
| `new_password` | `newpassword456` | `newpassword456` |
| `email_otp` | (kosong) | (kosong) |
| `phone_otp` | (kosong) | (kosong) |
| `reset_otp` | (kosong) | (kosong) |
| `token` | (kosong) | (kosong) |

#### 4. Buat Requests

Lihat file `API_CURL_EXAMPLES.md` untuk detail setiap endpoint.

---

## 📝 Cara Menggunakan

### 1. Register User Baru

1. Buka request **"1. Register"**
2. Pastikan environment sudah aktif
3. (Optional) Edit variables `email`, `phoneNumber`, `fullName` di environment
4. Klik **"Send"**
5. **OTP akan otomatis tersimpan** ke environment variables:
   - `email_otp` → untuk verify email
   - `phone_otp` → untuk verify phone

**Response Example:**
```json
{
  "message": "Registrasi berhasil. Silakan verifikasi email dan nomor telepon Anda.",
  "data": {
    "id": 1,
    "email": "test@example.com"
  },
  "debug_email_otp": "123456",
  "debug_phone_otp": "654321"
}
```

---

### 2. Verify Email

1. Buka request **"2. Verify Email"**
2. OTP sudah otomatis terisi dari variable `{{email_otp}}`
3. Klik **"Send"**

**Response Example:**
```json
{
  "message": "Email verified successfully"
}
```

---

### 3. Verify Phone

1. Buka request **"3. Verify Phone"**
2. OTP sudah otomatis terisi dari variable `{{phone_otp}}`
3. Klik **"Send"**

**Response Example:**
```json
{
  "message": "Phone number verified successfully"
}
```

---

### 4. Login

1. Buka request **"4. Login"**
2. Klik **"Send"**
3. **Token akan otomatis tersimpan** ke variable `token`

**Response Example:**
```json
{
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "fullName": "Test User",
      "email": "test@example.com",
      "role": "user"
    }
  }
}
```

---

### 5. Forgot Password

1. Buka request **"5. Forgot Password"**
2. Klik **"Send"**
3. **Reset OTP akan otomatis tersimpan** ke variable `reset_otp`

**Response Example:**
```json
{
  "message": "Kode OTP untuk reset password telah dikirim.",
  "debug_otp": "789012"
}
```

---

### 6. Reset Password

1. Buka request **"7. Reset Password"**
2. OTP sudah otomatis terisi dari variable `{{reset_otp}}`
3. Password baru terisi dari variable `{{new_password}}`
4. Klik **"Send"**

**Response Example:**
```json
{
  "message": "Password updated successfully"
}
```

---

## 🔄 Complete Testing Flow

### Flow 1: Register → Verify → Login

1. **Register** → OTP tersimpan otomatis
2. **Verify Email** → Gunakan OTP yang tersimpan
3. **Verify Phone** → Gunakan OTP yang tersimpan
4. **Login** → Token tersimpan otomatis

### Flow 2: Forgot Password → Reset → Login

1. **Forgot Password** → Reset OTP tersimpan otomatis
2. **Verify OTP** (optional) → Gunakan reset OTP
3. **Reset Password** → Gunakan reset OTP
4. **Login** → Gunakan password baru

---

## 🎯 Tips & Tricks

### 1. Generate Email Unik untuk Testing

Edit variable `email` di environment dengan format:
```
test+{{$timestamp}}@example.com
```

Atau gunakan Postman dynamic variable:
```
{{$randomEmail}}
```

### 2. Generate Phone Unik untuk Testing

Edit variable `phoneNumber`:
```
0812{{$timestamp}}
```

### 3. Lihat Saved Variables

- Klik icon **"eye"** di kanan atas (sebelah environment dropdown)
- Lihat semua variables yang tersimpan
- Cek `email_otp`, `phone_otp`, `reset_otp`, `token`

### 4. Auto-Save OTP & Token

Collection sudah dilengkapi dengan **Test Scripts** yang otomatis menyimpan:
- Email OTP dari response Register
- Phone OTP dari response Register
- Reset OTP dari response Forgot Password
- JWT Token dari response Login

**Contoh Test Script (sudah include di collection):**
```javascript
var jsonData = pm.response.json();

if (jsonData.debug_email_otp) {
    pm.environment.set("email_otp", jsonData.debug_email_otp);
}

if (jsonData.data && jsonData.data.token) {
    pm.environment.set("token", jsonData.data.token);
}
```

### 5. Clear All OTPs

Jika ingin reset semua OTP:
1. Buka Environment
2. Set value `email_otp`, `phone_otp`, `reset_otp` menjadi kosong
3. Save

---

## 🔧 Troubleshooting

### 1. "Could not send request"

**Problem:** Server tidak running

**Solution:**
```bash
cd "BE - Jalan"
./be-jalan
```

### 2. "User not found"

**Problem:** Email/phone belum terdaftar

**Solution:**
- Jalankan request **"1. Register"** terlebih dahulu
- Pastikan email di environment sama dengan yang diregister

### 3. "Invalid OTP"

**Problem:** OTP salah atau sudah expired

**Solution:**
- Cek variable `email_otp`, `phone_otp`, atau `reset_otp` di environment
- Request OTP baru jika sudah expired (15 menit)

### 4. "OTP has expired"

**Problem:** OTP sudah lebih dari 15 menit

**Solution:**
- Untuk register: Daftar ulang dengan email/phone baru
- Untuk forgot password: Request OTP baru

### 5. Variables tidak tersimpan otomatis

**Problem:** Test script tidak berjalan

**Solution:**
- Pastikan menggunakan collection yang sudah di-import
- Cek tab "Tests" di setiap request, pastikan ada script
- Pastikan environment sudah aktif

---

## 📊 Request Order (Recommended)

Untuk testing lengkap, jalankan request dengan urutan:

1. ✅ **Register** → Dapat email_otp & phone_otp
2. ✅ **Verify Email** → Email terverifikasi
3. ✅ **Verify Phone** → Phone terverifikasi
4. ✅ **Login** → Dapat token
5. ✅ **Forgot Password** → Dapat reset_otp
6. ✅ **Verify OTP** (optional)
7. ✅ **Reset Password** → Password berubah
8. ✅ **Login** (dengan password baru) → Dapat token baru

---

## 🎨 Customize Environment

Anda bisa membuat multiple environments untuk testing berbeda:

### Development Environment
```
base_url: http://localhost:8080/api/v1
```

### Staging Environment
```
base_url: https://staging-api.asita-jalan.com/api/v1
```

### Production Environment
```
base_url: https://api.asita-jalan.com/api/v1
```

---

## 📱 Testing OTP Channels

### Email OTP
- Cek inbox email yang didaftarkan
- Atau lihat `debug_email_otp` di response (development only)

### WhatsApp OTP
- Cek WhatsApp di nomor yang didaftarkan
- Atau lihat `debug_phone_otp` di response (development only)

### SMS OTP
- Cek SMS di nomor yang didaftarkan
- Atau lihat `debug_phone_otp` di response (development only)

**Note:** Di production, `debug_*` fields akan dihapus!

---

## 🔐 Security Notes

1. **Jangan commit** file environment dengan credentials asli
2. **Hapus debug OTP fields** sebelum production
3. **Gunakan HTTPS** di production
4. **Rotate JWT secret** secara berkala

---

## 📞 Support

Untuk pertanyaan atau issue:
- Lihat file `API_CURL_EXAMPLES.md` untuk curl examples
- Lihat file `IMPLEMENTASI_OTP.md` untuk detail implementasi
- Lihat file `README.md` untuk overview project

---

**Happy Testing! 🚀**

**Last Updated:** 2026-02-18
