# 🚀 CariGo – Senior-Level Agresif Mimari ile KOBİ Tahsilat & Mutabakat Sistemi (MVP+)

> **Bu doküman JUNIOR işi değildir.**
> CariGo, *"çalışsın yeter"* projesi değil; **domain-driven**, **test edilebilir**, **yarın SaaS’a evrilebilir**, **senior showcase** bir backend mimarisiyle inşa edilir.

---

# 🧠 MİMARİ MANİFESTO (OKUMADAN KOD YAZMA)

### ❗ KIRMIZI ÇİZGİLER

* ❌ God handler yok
* ❌ DB modeli = domain modeli değil
* ❌ Handler içinde iş kuralı YOK
* ❌ Float / money hack YOK
* ❌ "Sonra refactor" YOK

### ✅ ZORUNLU PRENSİPLER

* Clean Architecture + DDD (pragmatik)
* Dependency Inversion
* Explicit boundaries
* Test-first domain
* Infrastructure detayları **içeri sızamaz**

---

# 🧱 KATMANLI MİMARİ (NET SINIRLAR)

```txt
┌─────────────────────────────┐
│        HTTP / UI Layer      │  → Gin handlers, DTO
├─────────────────────────────┤
│      Application Layer      │  → UseCases / Services
├─────────────────────────────┤
│        Domain Layer         │  → Entities, Rules
├─────────────────────────────┤
│   Infrastructure Layer     │  → SQLite, Repo impl
└─────────────────────────────┘
```

> **Üst katman alt katmanı TANIMAZ**
> Sadece interface görür.

---

# 📁 KESİN PROJE YAPISI (DEĞİŞMEZ)

```txt
carigo/
├── cmd/
│   └── api/
│       └── main.go          # Sadece wiring
│
├── internal/
│   ├── domain/
│   │   ├── customer.go
│   │   ├── invoice.go
│   │   ├── payment.go
│   │   ├── allocation.go
│   │   ├── money.go         # VALUE OBJECT
│   │   └── errors.go
│   │
│   ├── application/
│   │   ├── ports/           # Interfaces
│   │   │   ├── repositories.go
│   │   │   └── clock.go
│   │   ├── usecases/
│   │   │   ├── create_invoice.go
│   │   │   ├── register_payment.go
│   │   │   └── generate_statement.go
│   │   └── dto/
│   │
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── sqlite/
│   │   │   │   ├── customer_repo.go
│   │   │   │   ├── invoice_repo.go
│   │   │   │   └── allocation_repo.go
│   │   └── clock/
│   │
│   ├── interfaces/
│   │   └── http/
│   │       ├── handlers/
│   │       ├── middleware/
│   │       └── router.go
│   │
│   └── bootstrap/
│       └── container.go     # DI manual
│
├── migrations/
├── web/                     # Iconic template
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

# 🧠 DOMAIN LAYER (EN KUTSAL YER)

## 🎯 Amaç

**İş kuralları burada yaşar.**
DB, HTTP, JSON umurunda değil.

### Task D1 – Money Value Object

* amount_cents
* Safe add / subtract
* Negatif koruması

### Task D2 – Invoice Entity

* Open / Partial / Paid state machine
* Allocation sonrası otomatik status değişimi

### Task D3 – Allocation Rules

* Bir payment → N invoice
* Bir invoice ← N payment
* Overpayment desteklenir

> ❗ Allocation logic **domain testleri olmadan geçmez**

---

# 🧠 APPLICATION LAYER (AKIL)

## 🎯 Amaç

Use-case bazlı, senaryoya göre çalışan katman.

### Task A1 – RegisterPayment UseCase

* Input: customer_id, amount, date
* Açık faturaları çek
* Allocation planı üret

### Task A2 – ConfirmAllocation UseCase

* Allocation onay
* Transactional commit

### Task A3 – GenerateStatement UseCase

* Mutabakat hesapla
* DTO üret

---

# 🗄️ INFRASTRUCTURE LAYER (DETAYLAR)

## 🎯 Amaç

Domain’e **asla** iş kuralı sokmadan persistence.

### Task I1 – SQLite Repo Implementasyonları

* Transaction support
* Lock & concurrency safe

### Task I2 – Migration Runner

* Idempotent
* Startup’ta çalışır

---

# 🌐 INTERFACE / HTTP LAYER

## 🎯 Amaç

Sadece translate eder.
Request → UseCase → Response

### Task H1 – DTO Mapping

* Validation burada
* Domain leak yok

### Task H2 – HTML Page Wiring

* Iconic template
* Server-render

---

# 🐳 BOOTSTRAP & DI

### Task B1 – Manual Dependency Injection

* Constructor injection
* Interface binding

---

# 🧪 TEST STRATEJİSİ (ZORUNLU)

* Domain: %100 unit
* UseCase: happy path + edge
* Repo: integration

---

# 🏁 MVP BİTİŞ KRİTERİ (NET)

* Allocation %100 doğru
* Domain testleri yeşil
* UI → gerçek senaryo çalışıyor

---

> **Bu mimariyle yapılan CariGo:**
>
> * CV’ye girer
> * Senior interview’da anlatılır
> * 2 yıl sonra çöpe atılmaz

🔥 **Bu noktadan sonra bu bir oyuncak değil, mühendisliktir.**
