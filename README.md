# 🚀 ASITA - Jalan Jalan Backend API

Backend API untuk aplikasi ASITA - Jalan Jalan dengan fitur autentikasi dan OTP multi-channel.

## ✨ Features

- ✅ User Registration dengan Email & Phone Verification
- ✅ Multi-Channel OTP (Email, WhatsApp, SMS)
- ✅ Login dengan JWT Authentication
- ✅ Forgot Password & Reset Password
- ✅ Google OAuth (Coming Soon)
- ✅ Secure Password Hashing (bcrypt)
- ✅ OTP Expiration (15 minutes)

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Email:** SMTP (Gmail/SendGrid)
- **WhatsApp/SMS:** GoSMSGateway API
- **Authentication:** JWT

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Gmail account (untuk SMTP) atau SMTP server lainnya
- GoSMSGateway account (untuk WhatsApp & SMS)

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd "BE - Jalan"
```

### 2. Setup Environment Variables

```bash
cp .env.example .env
```

Edit `.env` dan isi dengan kredensial yang valid:

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

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Database Migration

Database akan otomatis di-migrate saat aplikasi pertama kali dijalankan.

### 5. Build & Run

```bash
# Build
go build -o be-jalan .

# Run
./be-jalan
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. **Register**
```http
POST /auth/register
Content-Type: application/json

{
  "fullName": "John Doe",
  "email": "john@example.com",
  "phoneNumber": "081234567890",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Registrasi berhasil. Silakan verifikasi email dan nomor telepon Anda.",
  "data": {
    "id": 1,
    "email": "john@example.com"
  },
  "debug_email_otp": "123456",
  "debug_phone_otp": "654321"
}
```

#### 2. **Verify Email**
```http
POST /auth/verify-email
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "123456"
}
```

#### 3. **Verify Phone**
```http
POST /auth/verify-phone
Content-Type: application/json

{
  "phoneNumber": "081234567890",
  "otp": "654321"
}
```

#### 4. **Login**
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

#### 5. **Forgot Password**
```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### 6. **Verify OTP (Forgot Password)**
```http
POST /auth/verify-otp
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "123456"
}
```

#### 7. **Reset Password**
```http
POST /auth/reset-password
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "123456",
  "newPassword": "newpassword123",
  "confirmPassword": "newpassword123"
}
```

## 🧪 Testing

### Automated Testing

Gunakan script testing yang sudah disediakan:

```bash
# Pastikan server sudah running
./be-jalan

# Di terminal lain, jalankan test script
./test_otp.sh
```

### Manual Testing

Lihat file `IMPLEMENTASI_OTP.md` untuk contoh curl commands lengkap.

## 📁 Project Structure

```
BE - Jalan/
├── config/              # Database configuration
├── handlers/            # HTTP handlers/controllers
│   └── auth_handler.go
├── models/              # Database models
│   └── auth_models.go
├── services/            # Business logic services
│   └── notification/    # OTP notification services
│       ├── email.go
│       ├── whatsapp.go
│       ├── sms.go
│       └── config.go
├── utils/               # Utility functions
│   └── jwt.go
├── main.go              # Application entry point
├── .env                 # Environment variables (gitignored)
├── .env.example         # Environment template
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── test_otp.sh          # Automated test script
├── IMPLEMENTASI_OTP.md  # OTP implementation guide
└── README.md            # This file
```

## 🔐 Security Features

### Password Security
- Passwords hashed using **bcrypt** (cost 14)
- Never stored or returned in plain text

### OTP Security
- OTPs hashed using **bcrypt** before storage
- 15-minute expiration time
- Cleared after successful verification
- 6-digit random generation

### JWT Security
- Configurable secret key
- Configurable expiration time
- Secure token generation

## 📝 Configuration Guide

### Gmail SMTP Setup

1. Enable 2-Factor Authentication di Google Account
2. Generate App Password:
   - Go to https://myaccount.google.com/apppasswords
   - Select "Mail" and your device
   - Copy the 16-character password
3. Use App Password sebagai `SMTP_PASSWORD`

### GoSMSGateway Setup

1. Daftar di https://api.gosmsgateway.com
2. Top up balance untuk WhatsApp & SMS
3. Copy API Key dan Secret dari dashboard
4. Paste ke `.env` file

## 🐛 Troubleshooting

### Email tidak terkirim
- Cek SMTP credentials di `.env`
- Pastikan App Password Gmail benar
- Cek firewall untuk port 587

### WhatsApp/SMS tidak terkirim
- Cek API Key dan Secret
- Pastikan balance cukup
- Cek format nomor telepon (tanpa +62)

### Database connection error
- Pastikan PostgreSQL running
- Cek credentials di `.env`
- Cek database sudah dibuat

## 📄 License

MIT License

## 👥 Contributors

- Development Team

## 📞 Support

Untuk pertanyaan atau issue, silakan hubungi tim development.

---

**Last Updated:** 2026-02-18
