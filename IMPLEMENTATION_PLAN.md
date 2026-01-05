# 🛠️ CARIGO IMPLEMENTATION PLAN (SENIOR LEVEL)

Bu plan, `task.md` içerisindeki manifestoya %100 sadık kalınarak hazırlanmıştır. Junior ihmalleri kabul edilmez.

## 📅 FAZ 1: TEMEL VE DOMAIN (THE CORE)
**Hedef:** İş kurallarının (Business Logic) dış dünyadan izole bir şekilde, test edilebilir halde inşası.

- [ ] **1.1. Proje İskeleti Kurulumu**
    - Go module init (`go mod init carigo`)
    - Dizin yapısının `task.md`'ye birebir uygun oluşturulması.
    - `cmd`, `internal`, `web`, `migrations` klasörlerinin açılması.

- [ ] **1.2. Domain Layer: Money Value Object (Task D1)**
    - `internal/domain/money.go`
    - `int64` tabanlı kuruş/cent hesabı.
    - `Add`, `Subtract` metodları (Error safe).
    - Immutable yapı.

- [ ] **1.3. Domain Layer: Entities (Task D2)**
    - `internal/domain/invoice.go`: Fatura durum makinesi (Open -> Paid).
    - `internal/domain/payment.go`: Ödeme entity'si.
    - `internal/domain/customer.go`: Müşteri entity'si.
    - `internal/domain/allocation.go`: Ödeme ve fatura eşleşmesi.

- [ ] **1.4. Domain Layer: Unit Tests**
    - `money_test.go`: Kuruş hesabı şaşamaz.
    - `invoice_test.go`: Status geçişleri kontrolü.

## 📅 FAZ 2: APPLICATION LAYER (USE CASES)
**Hedef:** Domain objelerini yöneten, senaryoları (Use Cases) işleten katman.

- [ ] **2.1. Port Tanımları (Interfaces)**
    - `internal/application/ports/repositories.go`: Repository interface'leri (ICustomerRepo, IInvoiceRepo).
    - `internal/application/ports/clock.go`: Zaman bağımlılığını soyutlama.

- [ ] **2.2. Use Cases (Task A1, A2, A3)**
    - `CreateInvoice`: Fatura oluşturma.
    - `RegisterPayment`: Tahsilat girişi ve otomatik dağıtım (Allocation) stratejisi.
    - `GenerateStatement`: Hesap ekstresi.

## 📅 FAZ 3: INFRASTRUCTURE LAYER (PERSISTENCE)
**Hedef:** Veritabanı ve dış dünya entegrasyonu. Domain burayı bilmez.

- [ ] **3.1. SQLite Setup**
    - GORM veya raw SQL (Performance odaklı seçim).
    - `internal/infrastructure/persistence/sqlite/`.
    - Repository interface maplemeleri.

- [ ] **3.2. Migrations**
    - `migrations/` klasörü altında SQL dosyaları.
    - Uygulama başlangıcında opsiyonel auto-migrate.

## 📅 FAZ 4: HTTP / INTERFACE LAYER
**Hedef:** Dış dünyadan gelen istekleri Use Case'lere çevirmek.

- [ ] **4.1. Gin Setup & Middleware**
    - Router yapılandırması (`internal/interfaces/http/router.go`).
    - Error handling middleware.

- [ ] **4.2. Handlers & DTOs**
    - Request/Response struct'ları (`dto/`).
    - Handler fonksiyonları (Logic barındırmaz!).

## 📅 FAZ 5: WIRING (MAIN)
**Hedef:** Tüm bağımlılıkların (Dependency Injection) bağlanması.

- [ ] **5.1. Bootstrap**
    - `internal/bootstrap/container.go`: Manual DI.
    - `cmd/api/main.go`: Uygulamayı ayağa kaldırma.

---
**NOT:** Her adımda "Clean Architecture" kuralları ihlal edilirse PR reddedilir (Simülasyon).
