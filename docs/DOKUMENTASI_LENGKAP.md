# 📚 Dokumentasi Lengkap API - Quick Reference

File ini berisi daftar lengkap semua dokumentasi yang tersedia untuk API ASITA - Jalan Jalan.

---

## 📁 File Dokumentasi

### 1. **README.md** 
📖 **Overview Project & Quick Start Guide**
- Pengenalan project
- Tech stack yang digunakan
- Cara setup dan instalasi
- Struktur project
- Security features

**Baca file ini untuk:** Memahami project secara keseluruhan dan cara setup awal.

---

### 2. **API_CURL_EXAMPLES.md** ⭐
🔧 **Koleksi Lengkap Curl Commands**
- Semua endpoint dengan curl examples
- Request & response examples
- Complete flow examples
- Testing tips dengan bash

**Baca file ini untuk:** Copy-paste curl commands untuk testing manual atau di terminal.

---

### 3. **POSTMAN_GUIDE.md** ⭐⭐
📮 **Panduan Lengkap Postman**
- Cara import collection & environment
- Cara menggunakan setiap endpoint
- Auto-save OTP & token
- Tips & tricks Postman
- Troubleshooting

**Baca file ini untuk:** Setup dan testing menggunakan Postman (RECOMMENDED).

---

### 4. **IMPLEMENTASI_OTP.md**
📱 **Detail Implementasi OTP Multi-Channel**
- Cara kerja OTP system
- Flow diagram
- Security features
- Configuration guide
- Troubleshooting OTP delivery

**Baca file ini untuk:** Memahami detail implementasi OTP (Email, WhatsApp, SMS).

---

### 5. **.env.example**
🔐 **Template Environment Variables**
- Semua environment variables yang diperlukan
- Penjelasan setiap variable
- Cara mendapatkan credentials

**Baca file ini untuk:** Setup environment variables dengan benar.

---

## 📦 File Postman

### 1. **ASITA_Jalan_Jalan_API.postman_collection.json** ⭐⭐⭐
📮 **Postman Collection (IMPORT FILE INI!)**
- Semua endpoint API
- Auto-save scripts untuk OTP & token
- Request examples

**Import file ini ke Postman untuk testing.**

---

### 2. **ASITA_Jalan_Jalan.postman_environment.json** ⭐⭐⭐
🌍 **Postman Environment (IMPORT FILE INI!)**
- Pre-configured variables
- Test data
- Auto-populated OTP & token

**Import file ini ke Postman untuk environment setup otomatis.**

---

## 🧪 File Testing

### **test_otp.sh**
🤖 **Automated Testing Script**
- Test semua endpoint secara otomatis
- Interactive flow
- Colored output

**Jalankan script ini untuk automated testing:**
```bash
./test_otp.sh
```

---

## 🚀 Quick Start Recommendations

### Untuk Testing dengan Postman (RECOMMENDED):

1. ✅ Baca: **POSTMAN_GUIDE.md**
2. ✅ Import: **ASITA_Jalan_Jalan_API.postman_collection.json**
3. ✅ Import: **ASITA_Jalan_Jalan.postman_environment.json**
4. ✅ Aktifkan environment di Postman
5. ✅ Mulai testing!

### Untuk Testing dengan Curl:

1. ✅ Baca: **API_CURL_EXAMPLES.md**
2. ✅ Copy-paste curl commands
3. ✅ Jalankan di terminal

### Untuk Automated Testing:

1. ✅ Jalankan server: `./be-jalan`
2. ✅ Jalankan test: `./test_otp.sh`

---

## 📊 API Endpoints Summary

| No | Endpoint | Method | Description |
|----|----------|--------|-------------|
| 1 | `/auth/register` | POST | Register user baru + kirim OTP |
| 2 | `/auth/verify-email` | POST | Verifikasi email dengan OTP |
| 3 | `/auth/verify-phone` | POST | Verifikasi phone dengan OTP |
| 4 | `/auth/login` | POST | Login dan dapat JWT token |
| 5 | `/auth/forgot-password` | POST | Request OTP untuk reset password |
| 6 | `/auth/verify-otp` | POST | Verifikasi OTP (optional) |
| 7 | `/auth/reset-password` | POST | Reset password dengan OTP |
| 8 | `/auth/google` | POST | Login dengan Google OAuth |

---

## 🎯 Testing Flow

### Complete Registration Flow:
```
Register → Verify Email → Verify Phone → Login
```

### Forgot Password Flow:
```
Forgot Password → Reset Password → Login (with new password)
```

---

## 🔧 Configuration Files

### **.env** (HARUS DIISI!)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=bosani
DB_PASSWORD=1234567890
DB_NAME=astpay_jalan

# JWT
JWT_SECRET=supersecretkey_change_me_in_production
JWT_EXPIRATION_HOURS=24

# SMTP Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@asita-jalan.com
SMTP_FROM_NAME=ASITA - Jalan Jalan

# WhatsApp & SMS (GoSMSGateway)
GOSMS_API_KEY=your-api-key
GOSMS_API_SECRET=your-api-secret
SMS_API_KEY=your-api-key
SMS_API_SECRET=your-api-secret
SMS_SENDER_ID=ASITA
```

---

## 📱 OTP Channels

| Channel | Provider | Status |
|---------|----------|--------|
| Email | SMTP (Gmail/SendGrid) | ✅ Active |
| WhatsApp | GoSMSGateway API | ✅ Active |
| SMS | GoSMSGateway API | ✅ Active |

---

## 🔐 Security Features

- ✅ Password hashing (bcrypt)
- ✅ OTP hashing (bcrypt)
- ✅ OTP expiration (15 minutes)
- ✅ JWT authentication
- ✅ Multi-channel redundancy

---

## 📞 Need Help?

| Topic | File to Read |
|-------|--------------|
| Setup project | `README.md` |
| Testing dengan Postman | `POSTMAN_GUIDE.md` |
| Testing dengan curl | `API_CURL_EXAMPLES.md` |
| Detail OTP implementation | `IMPLEMENTASI_OTP.md` |
| Environment setup | `.env.example` |

---

## 🎉 Ready to Test!

**Pilih salah satu:**

### Option 1: Postman (Recommended)
```
1. Import ASITA_Jalan_Jalan_API.postman_collection.json
2. Import ASITA_Jalan_Jalan.postman_environment.json
3. Baca POSTMAN_GUIDE.md
4. Start testing!
```

### Option 2: Curl
```
1. Baca API_CURL_EXAMPLES.md
2. Copy-paste curl commands
3. Start testing!
```

### Option 3: Automated Script
```bash
./test_otp.sh
```

---

**Happy Testing! 🚀**

**Last Updated:** 2026-02-18
