# pasar_malam

Aplikasi mobile marketplace jajanan pasar malam. User bisa browsing produk makanan & minuman, masukin ke keranjang, checkout, dan lihat status pesanan. Dibangun pakai Flutter dengan state management Provider.

# Struktur project
```text 
pasar_malam/
├── lib/
│   ├── core/                              # Core / shared modules
│   │   ├── constants/                     # API URL, warna, string
│   │   │   ├── api_constants.dart
│   │   │   ├── app_colors.dart
│   │   │   └── app_strings.dart
│   │   ├── providers/                     # Theme provider (dark mode)
│   │   │   └── theme_provider.dart
│   │   ├── routes/                        # Routing & navigasi
│   │   │   └── app_router.dart
│   │   ├── services/                      # Service layer
│   │   │   ├── dio_client.dart            # HTTP client setup
│   │   │   ├── secure_storage.dart        # Token storage
│   │   │   ├── notification_service.dart  # FCM handler
│   │   │   ├── biometric_lock_provider.dart
│   │   │   └── global_institute_pay_service.dart
│   │   ├── theme/                         # Theme data (light & dark)
│   │   │   └── app_theme.dart
│   │   └── widgets/                       # Shared widgets
│   │       ├── biometric_lock_screen.dart
│   │       └── swiss.dart
│   ├── features/                          # Feature modules
│   │   ├── auth/                          # Autentikasi
│   │   │   ├── data/models/               # Auth response model
│   │   │   ├── data/repositories/         # Auth repository impl
│   │   │   ├── domain/repositories/       # Auth repository interface
│   │   │   └── presentation/
│   │   │       ├── pages/                 # Login, Register, Verify Email
│   │   │       ├── providers/             # AuthProvider (state)
│   │   │       └── widgets/               # Button, text field, dll
│   │   ├── dashboard/                     # Produk & beranda
│   │   │   ├── data/models/               # Product model
│   │   │   ├── data/repositories/         # Product repository impl
│   │   │   ├── domain/repositories/       # Product repository interface
│   │   │   └── presentation/
│   │   │       ├── pages/                 # Dashboard page
│   │   │       └── providers/             # ProductProvider
│   │   ├── cart/                          # Keranjang
│   │   │   ├── data/models/               # Cart model
│   │   │   ├── data/repositories/         # Cart repository impl
│   │   │   ├── domain/repositories/       # Cart repository interface
│   │   │   └── presentation/
│   │   │       ├── pages/                 # Cart page
│   │   │       └── providers/             # CartProvider
│   │   └── order/                         # Pesanan
│   │       ├── data/models/               # Order model
│   │       ├── data/repositories/         # Order repository impl
│   │       ├── domain/repositories/       # Order repository interface
│   │       └── presentation/
│   │           ├── pages/                 # Checkout, My Orders, dll
│   │           └── providers/             # OrderProvider
│   ├── firebase_options.dart              # Konfigurasi Firebase
│   └── main.dart                          # Entry point
├── packages/
│   └── flutter_biometric_kit/             # Library biometric lokal
├── assets/
│   └── icons/
├── pubspec.yaml
└── README.md
```


# Screenshoot
<p align="center">
  <img src="assets/redme/login.png" width="220"/>
  <img src="assets/redme/cart.png" width="220"/>
  <img src="assets/redme/dashboard.png" width="220"/>
</p>

# Teknologi

- Flutter
- Provider
- Firebase Authentication
- REST API
- MySQL