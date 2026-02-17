# 📡 API Curl Examples - ASITA Jalan Jalan

Koleksi lengkap curl commands untuk testing semua endpoint API.

**Base URL:** `http://localhost:8080/api/v1`

---

## 📋 Table of Contents

1. [Register](#1-register)
2. [Verify Email](#2-verify-email)
3. [Verify Phone](#3-verify-phone)
4. [Login](#4-login)
5. [Forgot Password](#5-forgot-password)
6. [Verify OTP (Forgot Password)](#6-verify-otp-forgot-password)
7. [Reset Password](#7-reset-password)
8. [Google Auth](#8-google-auth)

---

## 1. Register

**Endpoint:** `POST /api/v1/auth/register`

**Description:** Register user baru dan kirim OTP ke Email, WhatsApp, dan SMS

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "John Doe",
    "email": "john.doe@example.com",
    "phoneNumber": "081234567890",
    "password": "password123"
  }'
```

### Request Body:

```json
{
  "fullName": "John Doe",
  "email": "john.doe@example.com",
  "phoneNumber": "081234567890",
  "password": "password123"
}
```

### Success Response (201 Created):

```json
{
  "message": "Registrasi berhasil. Silakan verifikasi email dan nomor telepon Anda.",
  "data": {
    "id": 1,
    "email": "john.doe@example.com"
  },
  "debug_email_otp": "123456",
  "debug_phone_otp": "654321"
}
```

### Error Responses:

**400 Bad Request:**
```json
{
  "error": "Invalid input"
}
```

**409 Conflict:**
```json
{
  "error": "Email already registered"
}
```

---

## 2. Verify Email

**Endpoint:** `POST /api/v1/auth/verify-email`

**Description:** Verifikasi email menggunakan OTP yang dikirim saat registrasi

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/verify-email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "otp": "123456"
  }'
```

### Request Body:

```json
{
  "email": "john.doe@example.com",
  "otp": "123456"
}
```

### Success Response (200 OK):

```json
{
  "message": "Email verified successfully"
}
```

### Error Responses:

**400 Bad Request:**
```json
{
  "error": "Invalid OTP"
}
```

```json
{
  "error": "OTP has expired"
}
```

**404 Not Found:**
```json
{
  "error": "User not found"
}
```

---

## 3. Verify Phone

**Endpoint:** `POST /api/v1/auth/verify-phone`

**Description:** Verifikasi nomor telepon menggunakan OTP yang dikirim via WhatsApp/SMS

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/verify-phone" \
  -H "Content-Type: application/json" \
  -d '{
    "phoneNumber": "081234567890",
    "otp": "654321"
  }'
```

### Request Body:

```json
{
  "phoneNumber": "081234567890",
  "otp": "654321"
}
```

### Success Response (200 OK):

```json
{
  "message": "Phone number verified successfully"
}
```

### Error Responses:

**400 Bad Request:**
```json
{
  "error": "Invalid OTP"
}
```

```json
{
  "error": "OTP has expired"
}
```

**404 Not Found:**
```json
{
  "error": "User not found"
}
```

---

## 4. Login

**Endpoint:** `POST /api/v1/auth/login`

**Description:** Login user dan dapatkan JWT token

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Request Body:

```json
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Success Response (200 OK):

```json
{
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "fullName": "John Doe",
      "email": "john.doe@example.com",
      "role": "user"
    }
  }
}
```

### Error Response (401 Unauthorized):

```json
{
  "error": "Invalid email or password"
}
```

---

## 5. Forgot Password

**Endpoint:** `POST /api/v1/auth/forgot-password`

**Description:** Request OTP untuk reset password. OTP akan dikirim via Email atau WhatsApp/SMS tergantung identifier.

### Curl Command (Email):

```bash
curl -X POST "http://localhost:8080/api/v1/auth/forgot-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com"
  }'
```

### Curl Command (Phone Number):

```bash
curl -X POST "http://localhost:8080/api/v1/auth/forgot-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "081234567890"
  }'
```

### Request Body:

```json
{
  "email": "john.doe@example.com"
}
```

**Note:** Field `email` bisa diisi dengan email atau nomor telepon.

### Success Response (200 OK):

```json
{
  "message": "Kode OTP untuk reset password telah dikirim.",
  "debug_otp": "789012"
}
```

### Error Response (404 Not Found):

```json
{
  "error": "User not found"
}
```

---

## 6. Verify OTP (Forgot Password)

**Endpoint:** `POST /api/v1/auth/verify-otp`

**Description:** Verifikasi OTP sebelum reset password (optional step)

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/verify-otp" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "otp": "789012"
  }'
```

### Request Body:

```json
{
  "email": "john.doe@example.com",
  "otp": "789012"
}
```

### Success Response (200 OK):

```json
{
  "message": "OTP valid"
}
```

### Error Responses:

**400 Bad Request:**
```json
{
  "error": "Invalid OTP"
}
```

```json
{
  "error": "OTP has expired"
}
```

---

## 7. Reset Password

**Endpoint:** `POST /api/v1/auth/reset-password`

**Description:** Reset password menggunakan OTP yang valid

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "otp": "789012",
    "newPassword": "newpassword456",
    "confirmPassword": "newpassword456"
  }'
```

### Request Body:

```json
{
  "email": "john.doe@example.com",
  "otp": "789012",
  "newPassword": "newpassword456",
  "confirmPassword": "newpassword456"
}
```

### Success Response (200 OK):

```json
{
  "message": "Password updated successfully"
}
```

### Error Responses:

**400 Bad Request:**
```json
{
  "error": "Invalid OTP"
}
```

```json
{
  "error": "OTP has expired"
}
```

**404 Not Found:**
```json
{
  "error": "User not found"
}
```

---

## 8. Google Auth

**Endpoint:** `POST /api/v1/auth/google`

**Description:** Login menggunakan Google OAuth (Coming Soon)

### Curl Command:

```bash
curl -X POST "http://localhost:8080/api/v1/auth/google" \
  -H "Content-Type: application/json" \
  -d '{
    "idToken": "google-id-token-here"
  }'
```

### Request Body:

```json
{
  "idToken": "google-id-token-here"
}
```

### Success Response (200 OK):

```json
{
  "message": "Login Google berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 123,
      "fullName": "Google User",
      "email": "google@example.com",
      "role": "user"
    }
  }
}
```

---

## 🔄 Complete Flow Examples

### Flow 1: Register → Verify Email → Verify Phone → Login

```bash
# Step 1: Register
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Jane Smith",
    "email": "jane.smith@example.com",
    "phoneNumber": "081298765432",
    "password": "securepass123"
  }'

# Response akan berisi debug_email_otp dan debug_phone_otp
# Contoh: "debug_email_otp": "123456", "debug_phone_otp": "654321"

# Step 2: Verify Email (gunakan OTP dari response)
curl -X POST "http://localhost:8080/api/v1/auth/verify-email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com",
    "otp": "123456"
  }'

# Step 3: Verify Phone (gunakan OTP dari response)
curl -X POST "http://localhost:8080/api/v1/auth/verify-phone" \
  -H "Content-Type: application/json" \
  -d '{
    "phoneNumber": "081298765432",
    "otp": "654321"
  }'

# Step 4: Login
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com",
    "password": "securepass123"
  }'
```

---

### Flow 2: Forgot Password → Reset Password → Login

```bash
# Step 1: Request OTP untuk reset password
curl -X POST "http://localhost:8080/api/v1/auth/forgot-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com"
  }'

# Response akan berisi debug_otp
# Contoh: "debug_otp": "789012"

# Step 2: (Optional) Verify OTP
curl -X POST "http://localhost:8080/api/v1/auth/verify-otp" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com",
    "otp": "789012"
  }'

# Step 3: Reset Password
curl -X POST "http://localhost:8080/api/v1/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com",
    "otp": "789012",
    "newPassword": "mynewpassword789",
    "confirmPassword": "mynewpassword789"
  }'

# Step 4: Login dengan password baru
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@example.com",
    "password": "mynewpassword789"
  }'
```

---

## 📝 Notes untuk Postman

### Import ke Postman:

1. **Buat Collection baru** dengan nama "ASITA - Jalan Jalan API"

2. **Set Environment Variables:**
   - `base_url`: `http://localhost:8080/api/v1`
   - `email`: `test@example.com`
   - `phone`: `081234567890`
   - `password`: `password123`
   - `token`: (akan di-set otomatis setelah login)

3. **Untuk setiap endpoint:**
   - Method: POST
   - URL: `{{base_url}}/auth/register` (contoh)
   - Headers: `Content-Type: application/json`
   - Body: raw JSON

4. **Auto-save token setelah login:**
   - Di tab "Tests" pada endpoint Login, tambahkan:
   ```javascript
   var jsonData = pm.response.json();
   pm.environment.set("token", jsonData.data.token);
   ```

5. **Gunakan token untuk protected endpoints:**
   - Headers: `Authorization: Bearer {{token}}`

---

## 🧪 Testing Tips

### 1. Generate Email Unik untuk Testing:
```bash
# Gunakan timestamp untuk email unik
EMAIL="test$(date +%s)@example.com"
```

### 2. Generate Phone Unik untuk Testing:
```bash
# Gunakan timestamp untuk phone unik
PHONE="0812$(date +%s | tail -c 9)"
```

### 3. Pretty Print JSON Response:
```bash
curl ... | jq '.'
```

### 4. Save Response to Variable (Bash):
```bash
RESPONSE=$(curl -s -X POST ...)
EMAIL_OTP=$(echo "$RESPONSE" | jq -r '.debug_email_otp')
```

---

## 🔒 Security Notes

1. **Debug OTP Fields:**
   - `debug_email_otp`, `debug_phone_otp`, `debug_otp` hanya untuk development
   - **HARUS dihapus** di production!

2. **Password Requirements:**
   - Minimum 6 karakter
   - Akan di-hash dengan bcrypt sebelum disimpan

3. **OTP Expiration:**
   - OTP berlaku selama **15 menit**
   - Setelah expired, harus request OTP baru

4. **Rate Limiting:**
   - Belum diimplementasikan
   - Recommended untuk production

---

## 📞 Support

Untuk pertanyaan atau issue, silakan hubungi tim development.

---

**Last Updated:** 2026-02-18
**Version:** 1.0.0
